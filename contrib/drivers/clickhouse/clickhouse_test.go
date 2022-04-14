package clickhouse

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

const (
	sqlVisitsDDL = `
CREATE TABLE IF NOT EXISTS visits (
id UInt64,
duration Float64,
url String,
created DateTime
) ENGINE = MergeTree() 
PRIMARY KEY id 
ORDER BY id
`
	dimSqlDDL = `
CREATE TABLE IF NOT EXISTS dim (
"code" String COMMENT '编码',
"translation" String COMMENT '译文',
"superior" UInt64 COMMENT '上级ID',
"row_number" UInt16 COMMENT '行号',
"is_active" UInt8 COMMENT '是否激活',
"is_preset" UInt8 COMMENT '是否预置',
"category" String COMMENT '类别',
"tree_path" Array(String) COMMENT '树路径',
"id" UInt64 COMMENT '代理主键ID',
"scd" UInt64 COMMENT '缓慢变化维ID',
"version" UInt64 COMMENT 'Merge版本ID',
"sign" Int8 COMMENT '标识位',
"created_by" UInt64 COMMENT '创建者ID',
"created_at" DateTime64(3,'Asia/Shanghai') COMMENT '创建时间',
"updated_by" UInt64 COMMENT '最后修改者ID',
"updated_at" DateTime64(3,'Asia/Shanghai') COMMENT '最后修改时间',
"updated_tick" UInt16 COMMENT '累计修改次数'
) ENGINE = ReplacingMergeTree("version")
ORDER BY ("id","scd")
COMMENT '会计准则';
`
	dimSqlDML = `
insert into dim (code, translation, superior, row_number, is_active, is_preset, category, tree_path, id, scd, version, sign, created_by, created_at, updated_by, updated_at, updated_tick)
values  ('CN', '{"zh_CN":"中国大陆会计准则","en_US":"Chinese mainland accounting legislation"}', 0, 1, 1, 1, 1, '[''CN'']', 607972403489804288, 0, 0, 0, 607536279118155777, '2017-09-06 00:00:00', 607536279118155777, '2017-09-06 00:00:00', 0),
        ('HK', '{"zh_CN":"中国香港会计准则","en_US":"Chinese Hong Kong accounting legislation"}', 0, 2, 1, 1, 1, '[''HK'']', 607972558544834566, 0, 0, 0, 607536279118155777, '2017-09-06 00:00:00', 607536279118155777, '2017-09-06 00:00:00', 0);
`
	factSqlDDL = `
CREATE TABLE IF NOT EXISTS fact (
"adjustment_level" UInt64 COMMENT '调整层ID',
"data_version" UInt64 COMMENT '数据版本ID',
"accounting_legislation" UInt64 COMMENT '会计准则ID',
"fiscal_year" UInt16 COMMENT '会计年度',
"fiscal_period" UInt8 COMMENT '会计期间',
"fiscal_year_period" UInt32 COMMENT '会计年度期间',
"legal_entity" UInt64 COMMENT '法人主体ID',
"cost_center" UInt64 COMMENT '成本中心ID',
"legal_entity_partner" UInt64 COMMENT '内部关联方ID',
"financial_posting" UInt64 COMMENT '凭证头ID',
"line" UInt16 COMMENT '行号',
"general_ledger_account" UInt64 COMMENT '总账科目ID',
"debit" Decimal64(9) COMMENT '借方金额',
"credit" Decimal64(9) COMMENT '贷方金额',
"transaction_currency" UInt64 COMMENT '交易币种ID',
"debit_tc" Decimal64(9) COMMENT '借方金额（交易币种）',
"credit_tc" Decimal64(9) COMMENT '贷方金额（交易币种）',
"posting_date" Date32 COMMENT '过账日期',
"gc_year" UInt16 COMMENT '公历年',
"gc_quarter" UInt8 COMMENT '公历季',
"gc_month" UInt8 COMMENT '公历月',
"gc_week" UInt8 COMMENT '公历周',
"raw_info" String COMMENT '源信息',
"summary" String COMMENT '摘要',
"id" UInt64 COMMENT '代理主键ID',
"version" UInt64 COMMENT 'Merge版本ID',
"sign" Int8 COMMENT '标识位'
) ENGINE = ReplacingMergeTree("version")
ORDER BY ("adjustment_level","data_version","legal_entity","fiscal_year","fiscal_period","financial_posting","line")
PARTITION BY ("adjustment_level","data_version","legal_entity","fiscal_year","fiscal_period")
COMMENT '数据主表';
`
	factSqlDML = `
insert into fact (adjustment_level, data_version, accounting_legislation, fiscal_year, fiscal_period, fiscal_year_period, legal_entity, cost_center, legal_entity_partner, financial_posting, line, general_ledger_account, debit, credit, transaction_currency, debit_tc, credit_tc, posting_date, gc_year, gc_quarter, gc_month, gc_week, raw_info, summary, id, version, sign)
values  (607970943242866688, 607973669943119880, 607972403489804288, 2022, 3, 202203, 607974511316307985, 0, 607976190010986520, 607996702456025136, 1, 607985607569838111, 8674.39, 0, 607974898261823505, 8674.39, 0, '2022-03-05', 2022, 1, 3, 11, '{}', '摘要', 607992882741121073, 0, 0),
        (607970943242866688, 607973669943119880, 607972403489804288, 2022, 4, 202204, 607974511316307985, 0, 607976190010986520, 607993586419503145, 1, 607985607569838111, 9999.88, 0, 607974898261823505, 9999.88, 0, '2022-04-10', 2022, 2, 4, 18, '{}', '摘要', 607996939140599857, 0, 0);
`
)

func clickhouseConfigDB() gdb.DB {
	connect, err := gdb.New(gdb.ConfigNode{
		Host:  "127.0.0.1",
		Port:  "9000",
		User:  "default",
		Pass:  "asdf_qwer",
		Name:  "default",
		Type:  "clickhouse",
		Debug: true,
	})
	gtest.AssertNil(err)
	gtest.AssertNE(connect, nil)
	return connect
}

func createClickhouseTableVisits(connect gdb.DB) error {
	_, err := connect.Exec(context.Background(), sqlVisitsDDL)
	return err
}

func createClickhouseTableDim(connect gdb.DB) error {
	_, err := connect.Exec(context.Background(), dimSqlDDL)
	return err
}

func createClickhouseTableFact(connect gdb.DB) error {
	_, err := connect.Exec(context.Background(), factSqlDDL)
	return err
}

func dropClickhouseTableVisits(conn gdb.DB) {
	sqlStr := fmt.Sprintf("DROP TABLE IF EXISTS `visits`")
	_, _ = conn.Exec(context.Background(), sqlStr)
}

func dropClickhouseTableDim(conn gdb.DB) {
	sqlStr := fmt.Sprintf("DROP TABLE IF EXISTS `dim`")
	_, _ = conn.Exec(context.Background(), sqlStr)
}

func dropClickhouseTableFact(conn gdb.DB) {
	sqlStr := fmt.Sprintf("DROP TABLE IF EXISTS `fact`")
	_, _ = conn.Exec(context.Background(), sqlStr)
}

func TestDriverClickhouse_Create(t *testing.T) {
	gtest.AssertNil(createClickhouseTableVisits(clickhouseConfigDB()))
}

func TestDriverClickhouse_New(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertNE(connect, nil)
	gtest.AssertNil(connect.PingMaster())
	gtest.AssertNil(connect.PingSlave())
}

func TestDriverClickhouse_OpenLink_Ping(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertNE(connect, nil)
	gtest.AssertNil(connect.PingMaster())
}

func TestDriverClickhouse_Tables(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertEQ(createClickhouseTableVisits(connect), nil)
	defer dropClickhouseTableVisits(connect)
	tables, err := connect.Tables(context.Background())
	gtest.AssertNil(err)
	gtest.AssertNE(len(tables), 0)
}

func TestDriverClickhouse_TableFields_Use_Config(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertNil(createClickhouseTableVisits(connect))
	defer dropClickhouseTableVisits(connect)
	field, err := connect.TableFields(context.Background(), "visits")
	gtest.AssertNil(err)
	gtest.AssertEQ(len(field), 4)
	gtest.AssertNQ(field, nil)
}

func TestDriverClickhouse_TableFields_Use_Link(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertNil(createClickhouseTableVisits(connect))
	defer dropClickhouseTableVisits(connect)
	field, err := connect.TableFields(context.Background(), "visits")
	gtest.AssertNil(err)
	gtest.AssertEQ(len(field), 4)
	gtest.AssertNQ(field, nil)
}

func TestDriverClickhouse_Transaction(t *testing.T) {
	connect := clickhouseConfigDB()
	defer dropClickhouseTableVisits(connect)
	gtest.AssertNE(connect.Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
		return nil
	}), nil)
}

func TestDriverClickhouse_InsertIgnore(t *testing.T) {
	connect := clickhouseConfigDB()
	_, err := connect.InsertIgnore(context.Background(), "", nil)
	gtest.AssertEQ(err, errUnsupportedInsertIgnore)
}

func TestDriverClickhouse_InsertAndGetId(t *testing.T) {
	connect := clickhouseConfigDB()
	_, err := connect.InsertAndGetId(context.Background(), "", nil)
	gtest.AssertEQ(err, errUnsupportedInsertGetId)
}

func TestDriverClickhouse_InsertOne(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertEQ(createClickhouseTableVisits(connect), nil)
	defer dropClickhouseTableVisits(connect)
	_, err := connect.Model("visits").Data(g.Map{
		"duration": float64(grand.Intn(999)),
		"url":      gconv.String(grand.Intn(999)),
		"created":  time.Now(),
	}).Insert()
	gtest.AssertNil(err)
}

func TestDriverClickhouse_InsertMany(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertEQ(createClickhouseTableVisits(connect), nil)
	defer dropClickhouseTableVisits(connect)
	tx, err := connect.Begin(context.Background())
	gtest.AssertEQ(err, errUnsupportedBegin)
	gtest.AssertNil(tx)
}

func TestDriverClickhouse_Insert(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertEQ(createClickhouseTableVisits(connect), nil)
	defer dropClickhouseTableVisits(connect)
	type insertItem struct {
		Id       uint64    `orm:"id"`
		Duration float64   `orm:"duration"`
		Url      string    `orm:"url"`
		Created  time.Time `orm:"created"`
	}
	var (
		insertUrl = "https://goframe.org"
		total     = 0
		item      = insertItem{
			Duration: 1,
			Url:      insertUrl,
			Created:  time.Now(),
		}
	)
	_, err := connect.Model("visits").Data(item).Insert()
	gtest.AssertNil(err)
	_, err = connect.Model("visits").Data(item).Save()
	gtest.AssertNil(err)
	total, err = connect.Model("visits").Count()
	gtest.AssertNil(err)
	gtest.AssertEQ(total, 2)
	var list []*insertItem
	for i := 0; i < 50; i++ {
		list = append(list, &insertItem{
			Duration: float64(grand.Intn(999)),
			Url:      insertUrl,
			Created:  time.Now(),
		})
	}
	_, err = connect.Model("visits").Data(list).Insert()
	gtest.AssertNil(err)
	_, err = connect.Model("visits").Data(list).Save()
	gtest.AssertNil(err)
	total, err = connect.Model("visits").Count()
	gtest.AssertNil(err)
	gtest.AssertEQ(total, 102)
}

func TestDriverClickhouse_Insert_Use_Exec(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertEQ(createClickhouseTableFact(connect), nil)
	defer dropClickhouseTableFact(connect)
	_, err := connect.Exec(context.Background(), factSqlDML)
	gtest.AssertNil(err)
}

func TestDriverClickhouse_Delete(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertEQ(createClickhouseTableVisits(connect), nil)
	defer dropClickhouseTableVisits(connect)
	_, err := connect.Model("visits").Where("created >", "2021-01-01 00:00:00").Delete()
	gtest.AssertNil(err)
}

func TestDriverClickhouse_Update(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertEQ(createClickhouseTableVisits(connect), nil)
	defer dropClickhouseTableVisits(connect)
	_, err := connect.Model("visits").Where("created > ", "2021-01-01 15:15:15").Data(g.Map{
		"created": time.Now().Format("2006-01-02 15:04:05"),
	}).Update()
	gtest.AssertNil(err)
}

func TestDriverClickhouse_Replace(t *testing.T) {
	connect := clickhouseConfigDB()
	_, err := connect.Replace(context.Background(), "", nil)
	gtest.AssertEQ(err, errUnsupportedReplace)
}

func TestDriverClickhouse_DoFilter(t *testing.T) {
	rawSQL := "select * from visits where 1 = 1"
	this := Driver{}
	replaceSQL, _, err := this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	gtest.AssertEQ(rawSQL, replaceSQL)
	rawSQL = "update visit set url = '1'"
	replaceSQL, _, err = this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	// this SQL can't run ,clickhouse will report an error because there is no WHERE statement
	gtest.AssertEQ(replaceSQL, "ALTER TABLE visit update url = '1'")
	rawSQL = "delete from visit"
	replaceSQL, _, err = this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	// this SQL can't run ,clickhouse will report an error because there is no WHERE statement
	gtest.AssertEQ(replaceSQL, "ALTER TABLE visit delete")

	rawSQL = "update visit set url = '1' where url = '0'"
	replaceSQL, _, err = this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	// this SQL can't run ,clickhouse will report an error because there is no WHERE statement
	gtest.AssertEQ(replaceSQL, "ALTER TABLE visit update url = '1' where url = '0'")
	rawSQL = "delete from visit where url='0'"
	replaceSQL, _, err = this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	// this SQL can't run ,clickhouse will report an error because there is no WHERE statement
	gtest.AssertEQ(replaceSQL, "ALTER TABLE visit delete where url='0'")
}

func TestDriverClickhouse_Select(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertNil(createClickhouseTableVisits(connect))
	defer dropClickhouseTableVisits(connect)
	_, err := connect.Model("visits").Data(g.Map{
		"url":      "goframe.org",
		"duration": float64(1),
	}).Insert()
	gtest.AssertNil(err)
	temp, err := connect.Model("visits").Where("url", "goframe.org").Where("duration >= ", 1).One()
	gtest.AssertNil(err)
	gtest.AssertEQ(temp.IsEmpty(), false)
	_, err = connect.Model("visits").Data(g.Map{
		"url":      "goframe.org",
		"duration": float64(2),
	}).Insert()
	gtest.AssertNil(err)
	data, err := connect.Model("visits").Where("url", "goframe.org").Where("duration >= ", 1).All()
	gtest.AssertNil(err)
	gtest.AssertEQ(len(data), 2)
}

func TestDriverClickhouse_Exec_OPTIMIZE(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertNil(createClickhouseTableVisits(connect))
	defer dropClickhouseTableVisits(connect)
	sqlStr := "OPTIMIZE table visits"
	_, err := connect.Exec(context.Background(), sqlStr)
	gtest.AssertNil(err)
}

func TestDriverClickhouse_ExecInsert(t *testing.T) {
	connect := clickhouseConfigDB()
	gtest.AssertEQ(createClickhouseTableDim(connect), nil)
	defer dropClickhouseTableDim(connect)
	_, err := connect.Exec(context.Background(), dimSqlDML)
	gtest.AssertNil(err)
}
