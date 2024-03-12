// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/text/gstr"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

// DoInsert inserts or updates data for given table.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by oracle driver`,
		)
	}
	var (
		keys   []string
		values []string
		params []interface{}
	)
	// Retrieve the table fields and length.
	var (
		listLength  = len(list)
		valueHolder = make([]string, 0)
	)
	for k := range list[0] {
		keys = append(keys, k)
		valueHolder = append(valueHolder, "?")
	}
	var (
		batchResult    = new(gdb.SqlResult)
		charL, charR   = d.GetChars()
		keyStr         = charL + strings.Join(keys, charL+","+charR) + charR
		valueHolderStr = strings.Join(valueHolder, ",")
	)
	// Format "INSERT...INTO..." statement.
	intoStrArray := make([]string, 0)
	for i := 0; i < len(list); i++ {
		for _, k := range keys {
			if s, ok := list[i][k].(gdb.Raw); ok {
				params = append(params, gconv.String(s))
			} else {
				params = append(params, list[i][k])
			}
		}
		values = append(values, valueHolderStr)
		intoStrArray = append(
			intoStrArray,
			fmt.Sprintf(
				"INTO %s(%s) VALUES(%s)",
				table, keyStr, valueHolderStr,
			),
		)
		if len(intoStrArray) == option.BatchCount || (i == listLength-1 && len(valueHolder) > 0) {
			r, err := d.DoExec(ctx, link, fmt.Sprintf(
				"INSERT ALL %s SELECT * FROM DUAL",
				strings.Join(intoStrArray, " "),
			), params...)
			if err != nil {
				return r, err
			}
			if n, err := r.RowsAffected(); err != nil {
				return r, err
			} else {
				batchResult.Result = r
				batchResult.Affected += n
			}
			params = params[:0]
			intoStrArray = intoStrArray[:0]
		}
	}
	return batchResult, nil
}

// doSave support upsert for Oracle
func (d *Driver) doSave(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	if len(option.OnConflict) == 0 {
		return nil, gerror.NewCode(
			gcode.CodeMissingParameter, `Please specify conflict columns`,
		)
	}

	if len(list) == 0 {
		return nil, gerror.NewCode(
			gcode.CodeInvalidRequest, `Save operation list is empty by oracle driver`,
		)
	}

	var (
		one                    = list[0]
		charL, charR           = d.GetChars()
		valueCharL, valueCharR = "'", "'"

		conflictKeys   = option.OnConflict
		conflictKeySet = gset.New(true)

		// insertKeys:   Handle valid keys that need to be inserted
		// insertValues: Handle values that need to be inserted
		// updateValues: Handle values that need to be updated
		// queryValues:  Handle data that need to be upsert
		queryValues, insertKeys, insertValues, updateValues []string
	)

	// conflictKeys slice type conv to set type
	for _, conflictKey := range conflictKeys {
		conflictKeySet.Add(gstr.ToUpper(conflictKey))
	}

	for key, value := range one {
		saveValue := gconv.String(value)
		queryValues = append(
			queryValues,
			fmt.Sprintf(
				valueCharL+"%s"+valueCharR+" AS "+charL+"%s"+charR,
				saveValue, key,
			),
		)

		insertKeys = append(insertKeys, charL+key+charR)
		insertValues = append(insertValues, "T2."+charL+key+charR)

		// filter conflict keys in updateValues
		if !conflictKeySet.Contains(key) {
			updateValues = append(
				updateValues,
				fmt.Sprintf(`T1.%s = T2.%s`, charL+key+charR, charL+key+charR),
			)
		}
	}

	batchResult := new(gdb.SqlResult)
	sqlStr := parseSqlForUpsert(table, queryValues, insertKeys, insertValues, updateValues, conflictKeys)
	r, err := d.DoExec(ctx, link, sqlStr)
	if err != nil {
		return r, err
	}
	if n, err := r.RowsAffected(); err != nil {
		return r, err
	} else {
		batchResult.Result = r
		batchResult.Affected += n
	}
	return batchResult, nil
}

// parseSqlForUpsert
// MERGE INTO {{table}} T1
// USING ( SELECT {{queryValues}} FROM DUAL T2
// ON (T1.{{duplicateKey}} = T2.{{duplicateKey}} AND ...)
// WHEN NOT MATCHED THEN
// INSERT {{insertKeys}} VALUES {{insertValues}}
// WHEN MATCHED THEN
// UPDATE SET {{updateValues}}
func parseSqlForUpsert(table string,
	queryValues, insertKeys, insertValues, updateValues, duplicateKey []string,
) (sqlStr string) {
	var (
		queryValueStr   = strings.Join(queryValues, ",")
		insertKeyStr    = strings.Join(insertKeys, ",")
		insertValueStr  = strings.Join(insertValues, ",")
		updateValueStr  = strings.Join(updateValues, ",")
		duplicateKeyStr string
		pattern         = gstr.Trim(`MERGE INTO %s T1 USING (SELECT %s FROM DUAL) T2 ON (%s) WHEN NOT MATCHED THEN INSERT(%s) VALUES (%s) WHEN MATCHED THEN UPDATE SET %s`)
	)

	for index, keys := range duplicateKey {
		if index != 0 {
			duplicateKeyStr += " AND "
		}
		duplicateTmp := fmt.Sprintf("T1.%s = T2.%s", keys, keys)
		duplicateKeyStr += duplicateTmp
	}

	return fmt.Sprintf(pattern,
		table,
		queryValueStr,
		duplicateKeyStr,
		insertKeyStr,
		insertValueStr,
		updateValueStr,
	)
}
