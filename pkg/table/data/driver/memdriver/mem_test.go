package memdriver

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/hoorayui/core-framework/pkg/table/data/driver"
	cap "github.com/hoorayui/core-framework/pkg/table/proto"
	"github.com/hoorayui/core-framework/pkg/table/registry"
	rt "github.com/hoorayui/core-framework/pkg/table/registry/test"
	_ "github.com/araddon/qlbridge/qlbdriver"
)

type TableTemplate struct {
	Id          string             `db:"id" t_key:"true" t_fm:"SSTR" t_name:"ID" t_vt:"STRING"`            // ID|唯一id
	TName       string             `db:"name" t_fm:"SSTR" t_name:"模板名" t_vt:"STRING"`                      // 模板名|table_id+模板名需唯一
	TableId     string             `db:"table_id" t_fm:"SSTR" t_name:"表ID" t_vt:"STRING"`                  // 表ID|表ID
	FAccess     cap.FileAccessType `db:"f_access" t_fm:"SINT" t_name:"访问权限" t_vt:"OPTION"`                 // 访问权限|0-PRIVATE, 1-PUBLIC, 2-SHARED
	FCreateUser int64              `db:"f_create_user" t_agg:"SUM" t_fm:"SINT" t_name:"创建用户ID" t_vt:"INT"` // 创建用户ID|
}

func getDataList() []interface{} {
	dataList := []interface{}{
		&TableTemplate{Id: "1", TName: "11", TableId: "", FCreateUser: 1},
		&TableTemplate{Id: "2", TName: "string", TableId: "", FCreateUser: 1},
		&TableTemplate{Id: "3", TName: "11", TableId: "", FCreateUser: 1},
		&TableTemplate{Id: "4", TName: "11", TableId: "", FCreateUser: 1},
		&TableTemplate{Id: "5", TName: "string", TableId: "", FCreateUser: 1},
		&TableTemplate{Id: "6", TName: "11", TableId: "", FCreateUser: 1},
	}
	return dataList
}

func TestMemDriver_FindRows(t *testing.T) {
	d := NewMemDriver(reflect.TypeOf(TableTemplate{}), getDataList)
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

	_, err = d.FindRows(context.Background(), nil, tmd, conditions, []string{"Id", "TName"},
		[]*driver.AggregateColumn{}, &cap.PageParam{Page: 1, PageSize: 2}, &cap.OrderParam{})
	fmt.Println(err)
	// test.DisplayObject(ret.Rows)
	// test.DisplayObject(ret.PageInfo)
}
