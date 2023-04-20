package mysql

import (
	"encoding/json"
	"github.com/hoorayui/core-framework/util"
	"log"
	"testing"

	db "github.com/hoorayui/core-framework/pkg/cap/database/mysql"
	"github.com/hoorayui/core-framework/pkg/cap/test"
	"github.com/hoorayui/core-framework/pkg/cap/utils/idgen"
)

var testDB *db.DB

func init() {
	var err error
	testDB, err = db.NewTestDBFromEnvVar()
	if err != nil {
		panic(err)
	}
}

func TestTableTemplateMapper_CreateTemplate(t *testing.T) {
	ss, err := testDB.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	rand := idgen.NewRandomIDGenerator(5)
	defer func(ss *db.Session, err error) {
		_ = ss.Close(err)
	}(ss, err)
	mapper := NewTableTemplateMapper(ss)
	for i := 0; i < 10; i++ {
		id, _ := rand.Generate()
		err = mapper.CreateTemplate(&TableTpl{
			TableTemplate: TableTemplate{
				Id:          id,
				Name:        id + "_n",
				TableId:     "111",
				FAccess:     0,
				FCreateUser: "aaa",
				FCreateTime: util.Now(),
				FModTime:    util.Now(),
				Body:        []byte(`{"a":1}`),
			},
			ShareList: []TableTemplateShare{
				{UserId: "aaa"},
				{UserId: "bbb"},
				{UserId: "ccc"},
			},
		})
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestTableTemplateMapper_FindTemplates(t *testing.T) {
	ss, err := testDB.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer func(ss *db.Session, err error) {
		_ = ss.Close(err)
	}(ss, err)
	mapper := NewTableTemplateMapper(ss)
	templates, err := mapper.FindTemplates(FilterTemplateIDEquals("392109bf-96b4-11eb-ab69-005056afd813"))
	if err != nil {
		t.Fatal(err)
	}
	test.DisplayObject(templates)
}

func TestTableTemplateMapper_FindTemplate(t *testing.T) {
	ss, err := testDB.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer func(ss *db.Session, err error) {
		_ = ss.Close(err)
	}(ss, err)
	mapper := NewTableTemplateMapper(ss)
	template, err := mapper.FindTemplate("STgf0", true)
	if err != nil {
		t.Fatal(err)
	}
	test.DisplayObject(template)
}

func TestTableTemplateMapper_DeleteTemplates(t *testing.T) {
	ss, err := testDB.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer func(ss *db.Session, err error) {
		_ = ss.Close(err)
	}(ss, err)
	mapper := NewTableTemplateMapper(ss)
	affected, err := mapper.DeleteTemplates(FilterTemplateIDEquals("STgf0"))
	if err != nil {
		t.Fatal(err)
	}
	log.Println("affected =", affected)
}

func TestTableTemplateMapper_UpdateTableTemplate(t *testing.T) {
	ss, err := testDB.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer func(ss *db.Session, err error) {
		_ = ss.Close(err)
	}(ss, err)
	mapper := NewTableTemplateMapper(ss)
	err = mapper.UpdateTableTemplate("p7Cj9", "测试", 2, json.RawMessage(`{"test":77777777}`),
		[]TableTemplateShare{
			{UserId: "aaaa"},
			{UserId: "bbbb"},
			{UserId: "cccc"},
		})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTableTemplateMapper_FindTemplatesByShareUser(t *testing.T) {
	ss, err := testDB.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer func(ss *db.Session, err error) {
		_ = ss.Close(err)
	}(ss, err)
	mapper := NewTableTemplateMapper(ss)
	templates, err := mapper.FindTemplatesByShareUserAndTableID("111", "111")
	if err != nil {
		t.Fatal(err)
	}
	test.DisplayObject(templates)
}
