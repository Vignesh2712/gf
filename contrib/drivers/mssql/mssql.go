// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package mssql implements gdb.Driver, which supports operations for database MSSql.
//
// Note:
// 1. It does not support Save/Replace features.
// 2. It does not support LastInsertId.
package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
)

// Driver is the driver for SQL server database.
type Driver struct {
	*gdb.Core
}

const (
	quoteChar = `"`
)

var (
	tablesSqlTmp      = `SELECT NAME FROM SYSOBJECTS WHERE XTYPE='U' AND STATUS >= 0 ORDER BY NAME`
	tableFieldsSqlTmp = `
SELECT 
	a.name Field,
	CASE b.name 
		WHEN 'datetime' THEN 'datetime'
		WHEN 'numeric' THEN b.name + '(' + convert(varchar(20), a.xprec) + ',' + convert(varchar(20), a.xscale) + ')' 
		WHEN 'char' THEN b.name + '(' + convert(varchar(20), a.length)+ ')'
		WHEN 'varchar' THEN b.name + '(' + convert(varchar(20), a.length)+ ')'
		ELSE b.name + '(' + convert(varchar(20),a.length)+ ')' END AS Type,
	CASE WHEN a.isnullable=1 THEN 'YES' ELSE 'NO' end AS [Null],
	CASE WHEN exists (
		SELECT 1 FROM sysobjects WHERE xtype='PK' AND name IN (
			SELECT name FROM sysindexes WHERE indid IN (
				SELECT indid FROM sysindexkeys WHERE id = a.id AND colid=a.colid
			)
		)
	) THEN 'PRI' ELSE '' END AS [Key],
	CASE WHEN COLUMNPROPERTY(a.id,a.name,'IsIdentity')=1 THEN 'auto_increment' ELSE '' END Extra,
	isnull(e.text,'') AS [Default],
	isnull(g.[value],'') AS [Comment]
FROM syscolumns a
LEFT JOIN systypes b ON a.xtype=b.xtype AND a.xusertype=b.xusertype
INNER JOIN sysobjects d ON a.id=d.id AND d.xtype='U' AND d.name<>'dtproperties'
LEFT JOIN syscomments e ON a.cdefault=e.id
LEFT JOIN sys.extended_properties g ON a.id=g.major_id AND a.colid=g.minor_id
LEFT JOIN sys.extended_properties f ON d.id=f.major_id AND f.minor_id =0
WHERE d.name='%s'
ORDER BY a.id,a.colorder
`
	selectSqlTmp          = `SELECT * FROM (SELECT TOP %d * FROM (SELECT TOP %d %s) as TMP1_ ) as TMP2_ `
	selectWithOrderSqlTmp = `
SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s) as ROWNUMBER_, %s ) as TMP_ 
WHERE TMP_.ROWNUMBER_ > %d AND TMP_.ROWNUMBER_ <= %d
`
)

func init() {
	var err error
	tableFieldsSqlTmp = formatSqlTmp(tableFieldsSqlTmp)
	if err = gdb.Register(`mssql`, New()); err != nil {
		panic(err)
	}
}

// formatSqlTmp formats sql template string into one line.
func formatSqlTmp(sqlTmp string) string {
	var err error
	// format sql template string.
	sqlTmp, err = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(sqlTmp))
	if err != nil {
		panic(err)
	}
	sqlTmp, err = gregex.ReplaceString(`\s{2,}`, " ", gstr.Trim(sqlTmp))
	if err != nil {
		panic(err)
	}
	return sqlTmp
}

// New create and returns a driver that implements gdb.Driver, which supports operations for Mssql.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for SQL server.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{
		Core: core,
	}, nil
}

// Open creates and returns an underlying sql.DB object for mssql.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "sqlserver"
	)
	if config.Link != "" {
		// ============================================================================
		// Deprecated from v2.2.0.
		// ============================================================================
		source = config.Link
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`database=([\w\.\-]+)+`, "database="+config.Name, source)
		}
	} else {
		source = fmt.Sprintf(
			"user id=%s;password=%s;server=%s;port=%s;database=%s;encrypt=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
		if config.Extra != "" {
			var extraMap map[string]interface{}
			if extraMap, err = gstr.Parse(config.Extra); err != nil {
				return nil, err
			}
			for k, v := range extraMap {
				source += fmt.Sprintf(`;%s=%s`, k, v)
			}
		}
	}

	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}

// GetChars returns the security char for this type of database.
func (d *Driver) GetChars() (charLeft string, charRight string) {
	return quoteChar, quoteChar
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, sql string, args []interface{},
) (newSql string, newArgs []interface{}, err error) {
	var index int
	// Convert placeholder char '?' to string "@px".
	newSql, err = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf("@p%d", index)
	})
	if err != nil {
		return "", nil, err
	}
	newSql, err = gregex.ReplaceString("\"", "", newSql)
	if err != nil {
		return "", nil, err
	}
	newSql, err = d.parseSql(newSql)
	if err != nil {
		return "", nil, err
	}
	newArgs = args
	return d.Core.DoFilter(ctx, link, newSql, newArgs)
}

// parseSql does some replacement of the sql before commits it to underlying driver,
// for support of microsoft sql server.
func (d *Driver) parseSql(toBeCommittedSql string) (string, error) {
	var (
		err       error
		operation = gstr.StrTillEx(toBeCommittedSql, " ")
		keyword   = strings.ToUpper(gstr.Trim(operation))
	)
	switch keyword {
	case "SELECT":
		toBeCommittedSql, err = d.handleSelectSqlReplacement(toBeCommittedSql)
		if err != nil {
			return "", err
		}
	}
	return toBeCommittedSql, nil
}

func (d *Driver) handleSelectSqlReplacement(toBeCommittedSql string) (newSql string, err error) {
	// SELECT * FROM USER WHERE ID=1 LIMIT 1
	match, err := gregex.MatchString(`^SELECT(.+)LIMIT 1$`, toBeCommittedSql)
	if err != nil {
		return "", err
	}
	if len(match) > 1 {
		return fmt.Sprintf(`SELECT TOP 1 %s`, match[1]), nil
	}

	// SELECT * FROM USER WHERE AGE>18 ORDER BY ID DESC LIMIT 100, 200
	patten := `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,\s*(\d+))`
	if gregex.IsMatchString(patten, toBeCommittedSql) == false {
		return toBeCommittedSql, nil
	}
	allMatch, err := gregex.MatchAllString(patten, toBeCommittedSql)
	if err != nil {
		return "", err
	}
	var index = 1
	// LIMIT statement checks.
	if len(allMatch) < 2 ||
		(strings.HasPrefix(allMatch[index][0], "LIMIT") == false &&
			strings.HasPrefix(allMatch[index][0], "limit") == false) {
		return toBeCommittedSql, nil
	}
	if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", toBeCommittedSql) == false {
		return toBeCommittedSql, nil
	}
	// ORDER BY statement checks.
	var (
		selectStr = ""
		orderStr  = ""
		haveOrder = gregex.IsMatchString("((?i)SELECT)(.+)((?i)ORDER BY)", toBeCommittedSql)
	)
	if haveOrder {
		queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)ORDER BY)", toBeCommittedSql)
		if len(queryExpr) != 4 ||
			strings.EqualFold(queryExpr[1], "SELECT") == false ||
			strings.EqualFold(queryExpr[3], "ORDER BY") == false {
			return toBeCommittedSql, nil
		}
		selectStr = queryExpr[2]
		orderExpr, _ := gregex.MatchString("((?i)ORDER BY)(.+)((?i)LIMIT)", toBeCommittedSql)
		if len(orderExpr) != 4 ||
			strings.EqualFold(orderExpr[1], "ORDER BY") == false ||
			strings.EqualFold(orderExpr[3], "LIMIT") == false {
			return toBeCommittedSql, nil
		}
		orderStr = orderExpr[2]
	} else {
		queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", toBeCommittedSql)
		if len(queryExpr) != 4 ||
			strings.EqualFold(queryExpr[1], "SELECT") == false ||
			strings.EqualFold(queryExpr[3], "LIMIT") == false {
			return toBeCommittedSql, nil
		}
		selectStr = queryExpr[2]
	}
	first, limit := 0, 0
	for i := 1; i < len(allMatch[index]); i++ {
		if len(strings.TrimSpace(allMatch[index][i])) == 0 {
			continue
		}
		if strings.HasPrefix(allMatch[index][i], "LIMIT") ||
			strings.HasPrefix(allMatch[index][i], "limit") {
			first, _ = strconv.Atoi(allMatch[index][i+1])
			limit, _ = strconv.Atoi(allMatch[index][i+2])
			break
		}
	}
	if haveOrder {
		toBeCommittedSql = fmt.Sprintf(
			selectWithOrderSqlTmp,
			orderStr, selectStr, first, first+limit,
		)
		return toBeCommittedSql, nil
	}

	if first == 0 {
		first = limit
	}
	toBeCommittedSql = fmt.Sprintf(
		selectSqlTmp,
		limit, first+limit, selectStr,
	)
	return toBeCommittedSql, nil
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}

	result, err = d.DoSelect(ctx, link, tablesSqlTmp)
	if err != nil {
		return
	}
	for _, m := range result {
		for _, v := range m {
			tables = append(tables, v.String())
		}
	}
	return
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
//
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result     gdb.Result
		link       gdb.Link
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	structureSql := fmt.Sprintf(tableFieldsSqlTmp, table)
	result, err = d.DoSelect(ctx, link, structureSql)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		fields[m["Field"].String()] = &gdb.TableField{
			Index:   i,
			Name:    m["Field"].String(),
			Type:    m["Type"].String(),
			Null:    m["Null"].Bool(),
			Key:     m["Key"].String(),
			Default: m["Default"].Val(),
			Extra:   m["Extra"].String(),
			Comment: m["Comment"].String(),
		}
	}
	return fields, nil
}

// DoInsert inserts or updates data forF given table.
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Save operation is not supported by mssql driver`,
		)

	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by mssql driver`,
		)

	default:
		return d.Core.DoInsert(ctx, link, table, list, option)
	}
}
