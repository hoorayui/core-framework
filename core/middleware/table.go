package middleware

import (
	cap "github.com/hoorayui/core-framework/pkg/table/proto"
)

// 表格白名单，以下表格都不控制权限
var tableWhiteList = map[string]interface{}{
	noPermissionControlTableID: nil,
	"tables.QuickLink":         nil,
	"tables.SiteMsg":           nil,
	"tables.MyWInventoryLog":   nil,
}

const noPermissionControlTableID = "no-pc-table"

func tableID(req interface{}) string {
	if r, ok := req.(*cap.GetTableInfoReq); ok {
		return r.TableId
	} else if r, ok := req.(*cap.GetTableTemplatesReq); ok {
		return r.TableId
	} else if r, ok := req.(*cap.CreateTableTemplateReq); ok {
		return r.Template.TableId
	} else if _, ok := req.(*cap.DeleteTableTemplateReq); ok {
		return noPermissionControlTableID
	} else if r, ok := req.(*cap.GetTableColumnsReq); ok {
		return r.TableId
	} else if r, ok := req.(*cap.GetTableRowsReq); ok {
		return r.TableId
	} else if r, ok := req.(*cap.GetTableRowByIDReq); ok {
		return r.TableId
	} else if r, ok := req.(*cap.GetTableColumnOptionsReq); ok {
		return r.TableId
	} else if r, ok := req.(*cap.GetTableRowsLiteReq); ok {
		return r.TableId
	} else if _, ok := req.(*cap.GetOptionsReq); ok {
		return noPermissionControlTableID
	} else if r, ok := req.(*cap.DoRowFormActionReq); ok {
		return r.TableId
	} else {
		return ""
	}
}

func isWhiteListTable(req interface{}) bool {
	tid := tableID(req)
	if _, ok := tableWhiteList[tid]; ok {
		return true
	}
	return false
}
