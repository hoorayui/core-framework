package dbdriver

import (
	"context"
	"log"
	"testing"

	db "github.com/hoorayui/core-framework/pkg/cap/database/mysql"
	"github.com/hoorayui/core-framework/pkg/cap/test"
	"github.com/hoorayui/core-framework/pkg/table/data/driver"
	cap "github.com/hoorayui/core-framework/pkg/table/proto"
	"github.com/hoorayui/core-framework/pkg/table/registry"
	rt "github.com/hoorayui/core-framework/pkg/table/registry/test"
	_ "github.com/go-sql-driver/mysql"
)

var testDB *db.DB

func init() {
	var err error
	testDB, err = db.NewTestDBFromEnvVar()
	if err != nil {
		panic(err)
	}
}

func registerTableTemplate() {
}

func TestDBDataDriver_FindRows(t *testing.T) {
	// 注册选项
	registry.RegisterOptionFromProtoEnum(cap.FileAccessType(0))
	// 注册报表
	tmd, err := registry.LoadTMDFromStruct(&rt.TableTemplate{},
		func(t *cap.Template) *cap.Template {
			t.Body.Filter = &cap.FilterBody{}
			return t
		})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = registry.GlobalTableRegistry().TableMetaReg.Register(tmd)
	if err != nil {
		panic(err)
	}

	// 开始查询
	ss, err := testDB.NewSession()
	if err != nil {
		panic(err)
	}

	ddd := DBDataDriver{dbTableName: "table_template"}
	//	tpl := tmd.DefaultTpl()
	conditions := []*driver.Condition{}
	// tpl.Body.Filter.Conditions = append(tpl.Body.Filter.Conditions) // &cap.Condition{
	// 	ColumnId: "Id", OperatorId: "builtin.EQ",
	// 	Values: []*cap.FilterValue{
	// 		{Values: &cap.FilterValue_LiteralValues{LiteralValues: &cap.Value{V: &cap.Value_VString{VString: "54XX7"}}}}},
	// },
	// &cap.Condition{
	// 	ColumnId: "FCreateUser", OperatorId: "builtin.EQ",
	// 	Values: []*cap.FilterValue{
	// 		{Values: &cap.FilterValue_LiteralValues{LiteralValues: &cap.Value{V: &cap.Value_VInt{VInt: 1}}}}},
	// },

	ret, err := ddd.FindRows(context.Background(), ss, tmd, conditions, []string{"Id", "TName"},
		[]*driver.AggregateColumn{}, &cap.PageParam{Page: 1, PageSize: 2}, &cap.OrderParam{})
	test.DisplayObject(ret.Rows)
	test.DisplayObject(ret.PageInfo)
}
