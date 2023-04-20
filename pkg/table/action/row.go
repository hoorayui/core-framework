package action

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"framework/pkg/cap/database/mysql"
	"framework/pkg/cap/msg/errors"
	cap "framework/pkg/table/proto"
)

// ID creates id of an action
func ID(tableName, actionID string) string {
	return tableName + "." + actionID
}

// RowAction 行操作
type RowAction interface {
	ID() string
	Type() cap.RowActionType
	Name() string
}

// CustomRowAction 自定义行操作
type CustomRowAction struct {
	id   string
	name string
}

// NewCustomRowAction creates custom row action
func NewCustomRowAction(id, name string) RowAction {
	return &CustomRowAction{id: id, name: name}
}

// ID ...
func (a *CustomRowAction) ID() string {
	return a.id
}

// Type ...
func (a *CustomRowAction) Type() cap.RowActionType {
	return cap.RowActionType_RAT_CUSTOM
}

// Name ...
func (a *CustomRowAction) Name() string {
	return a.name
}

type SchemaForm interface{}

// FormRowAction 表单行操作
type FormRowAction struct {
	id       string
	name     string
	executor RowFormActionExecutor
	form     RowActionForm
}

// HrefProvider ...
type HrefProvider func(row interface{}) string

// HrefRowAction 超链接action
type HrefRowAction struct {
	id        string
	name      string
	hp        HrefProvider
	hrefStyle cap.HrefStyle
}

// ID ...
func (a *HrefRowAction) ID() string {
	return a.id
}

// Type ...
func (a *HrefRowAction) Type() cap.RowActionType {
	return cap.RowActionType_RAT_HREF
}

// Name ...
func (a *HrefRowAction) Name() string {
	return a.name
}

// Href ...
func (a *HrefRowAction) Href(row interface{}) string {
	return a.hp(row)
}

// HrefStyle ...
func (a *HrefRowAction) HrefStyle() cap.HrefStyle {
	return a.hrefStyle
}

// NewFormRowAction creates form row action
func NewFormRowAction(id, name string, executor RowFormActionExecutor, form RowActionForm) RowAction {
	return &FormRowAction{
		id:       id,
		name:     name,
		executor: executor,
		form:     form,
	}
}

// NewHrefRowAction creates form row action
func NewHrefRowAction(id, name string, hp HrefProvider, hrefStyle cap.HrefStyle) RowAction {
	return &HrefRowAction{
		id:        id,
		name:      name,
		hp:        hp,
		hrefStyle: hrefStyle,
	}
}

// RowActionForm 行操作表单
type RowActionForm interface {
	Schema(ctx context.Context, rowData interface{}) ([]byte, error)
}

// JSONSchemaActionProvider ...
type JSONSchemaActionProvider struct {
	obj interface{}
}

// RowFormActionExecutor 表单行操作执行器
type RowFormActionExecutor func(grpcCtx context.Context, ss *mysql.Session, formJSON []byte) error

// ID ...
func (a *FormRowAction) ID() string {
	return a.id
}

// Type ...
func (a *FormRowAction) Type() cap.RowActionType {
	return cap.RowActionType_RAT_JSON_FORM
}

// Name ...
func (a *FormRowAction) Name() string {
	return a.name
}

// ErrFormActionNotSupported 不支持的操作，Schema方法返回该错误，不显示操作按钮
var ErrFormActionNotSupported = errors.New("action not supported")

// Schema ...
func (a *FormRowAction) Schema(ctx context.Context, rowData interface{}) ([]byte, error) {
	return a.form.Schema(ctx, rowData)
}

// Execute ...
func (a *FormRowAction) Execute(grpcCtx context.Context, ss *mysql.Session, formJSON []byte) error {
	// TODO.. validate
	return a.executor(grpcCtx, ss, formJSON)
}

// RowFormSQLExecutor 行操作sql执行器
type RowFormSQLExecutor struct {
	tpl string
}

// NewRowFormSQLExecutor sqlTpl
func NewRowFormSQLExecutor(sqlTpl string) RowFormActionExecutor {
	return func(grpcCtx context.Context, ss *mysql.Session, formJSON []byte) error {
		tpl := template.New("RowFormSQLExecutor")
		data := map[string]interface{}{}
		err := json.Unmarshal(formJSON, &data)
		if err != nil {
			return err
		}
		tpl.Parse(sqlTpl)
		buf := bytes.NewBuffer([]byte{})
		err = tpl.Execute(buf, data)
		if err != nil {
			return errors.Wrap(err).Log()
		}
		_, err = ss.Exec(buf.String())
		if err != nil {
			return errors.Wrap(err).Log()
		}
		return nil
	}
}
