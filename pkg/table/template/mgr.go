package template

import (
	"fmt"
	"github.com/hoorayui/core-framework/util"
	"log"
	"time"

	db "github.com/hoorayui/core-framework/pkg/cap/database/mysql"
	"github.com/hoorayui/core-framework/pkg/cap/msg/errors"
	"github.com/hoorayui/core-framework/pkg/cap/utils/idgen"
	"github.com/hoorayui/core-framework/pkg/table/mysql"
	cap "github.com/hoorayui/core-framework/pkg/table/proto"
	"github.com/hoorayui/core-framework/pkg/table/registry"
	"github.com/golang/protobuf/proto"
)

// AccoutInfoProvider ...
type AccoutInfoProvider func(accountID string) (userName, displayName string, err error)

// AIP ...
var AIP AccoutInfoProvider

// Manager ...
type Manager struct{}

func (tm *Manager) genTemplateID() (string, error) {
	gen := idgen.NewUUIDGeneratorV1()
	return gen.Generate()
}

// CreateTemplate ...
func (tm *Manager) CreateTemplate(ss *db.Session, tableID, name string, body *cap.TemplateBody, accessType cap.FileAccessType,
	shareList []string, createUser string,
) (tplID string, err error) {
	mapper := mysql.NewTableTemplateMapper(ss)
	id, err := tm.genTemplateID()
	if err != nil {
		return id, errors.Wrap(err).Log()
	}
	templateBody, err := proto.Marshal(body)
	if err != nil {
		return "", errors.Wrap(err).Log()
	}
	tpl := &mysql.TableTpl{
		TableTemplate: mysql.TableTemplate{
			Id:          id,
			Name:        name,
			TableId:     tableID,
			FAccess:     int64(accessType),
			FCreateUser: createUser,
			FCreateTime: util.Now(),
			FModTime:    util.Now(),
			Body:        templateBody,
		},
	}
	if accessType == cap.FileAccessType_TA_SHARED {
		for _, s := range shareList {
			tpl.ShareList = append(tpl.ShareList, mysql.TableTemplateShare{UserId: s})
		}
	}
	err = mapper.CreateTemplate(tpl)
	if err != nil {
		return "", errors.Wrap(err).Log()
	}
	return id, nil
}

// DeleteTemplate ...
func (tm *Manager) DeleteTemplate(ss *db.Session, id, userID string) error {
	mapper := mysql.NewTableTemplateMapper(ss)
	tpl, err := mapper.FindTemplate(id, true)
	if err != nil {
		return errors.Wrap(err).Log()
	}
	// TODO. 其他高级权限校验？
	if tpl.FCreateUser != userID {
		return errors.Wrap(ErrOperatePermissionDenied).Log()
	}
	_, err = mapper.DeleteTemplates(mysql.FilterTemplateIDEquals(id))
	if err != nil {
		return errors.Wrap(err).Log()
	}
	return nil
}

// DeleteTemplateByCreateUser ...
func (tm *Manager) DeleteTemplateByCreateUser(ss *db.Session, currentUserID string) error {
	mapper := mysql.NewTableTemplateMapper(ss)
	_, err := mapper.DeleteTemplates(mysql.FilterCreateUserEquals(currentUserID))
	if err != nil {
		return errors.Wrap(err).Log()
	}
	return nil
}

// UpdateTemplate ...
func (tm *Manager) UpdateTemplate(ss *db.Session, tpl *cap.Template, currentUserID string) (*cap.Template, error) {
	mapper := mysql.NewTableTemplateMapper(ss)
	oldTpl, err := mapper.FindTemplate(tpl.Id, true)
	if err != nil {
		return nil, errors.Wrap(err).Log()
	}
	if oldTpl.FCreateUser != currentUserID {
		return nil, errors.Wrap(ErrOperatePermissionDenied).Log()
	}
	shareList := []mysql.TableTemplateShare{}
	templateBody, err := proto.Marshal(tpl.Body)
	if err != nil {
		return nil, errors.Wrap(err).Log()
	}
	if tpl.FileInfo.Access == cap.FileAccessType_TA_SHARED {
		for _, s := range tpl.FileInfo.ShareList {
			shareList = append(shareList, mysql.TableTemplateShare{UserId: s})
		}
	} else if tpl.FileInfo.Access == cap.FileAccessType_TA_PUBLIC {
		// XSG-3148
		return nil, errors.Wrap(ErrDisableSharingFormToAllUsers).Log()
	}
	err = mapper.UpdateTableTemplate(tpl.Id, tpl.Name, int(tpl.FileInfo.Access), templateBody, shareList)
	if err != nil {
		return nil, errors.Wrap(err).Log()
	}
	return tm.FindTemplate(ss, tpl.Id)
}

func mapTpl(src *mysql.TableTpl) (dst *cap.Template) {
	dst = &cap.Template{
		Id:      src.Id,
		Name:    src.Name,
		TableId: src.TableId,
		FileInfo: &cap.FileInfo{
			Access: cap.FileAccessType(src.FAccess),
			// TODO map account info
			CreateUser: &cap.UserInfo{
				Id: src.FCreateUser,
			},
			CreateTime: src.FCreateTime.Format(time.RFC3339),
			ModifyTime: src.FModTime.Format(time.RFC3339),
		},
		Body: &cap.TemplateBody{},
	}
	if AIP != nil {
		userName, displayName, _ := AIP(src.FCreateUser)
		dst.FileInfo.CreateUser.UserName = userName
		dst.FileInfo.CreateUser.DisplayName = displayName
	}
	if dst.FileInfo.Access == cap.FileAccessType_TA_SHARED {
		for _, s := range src.ShareList {
			dst.FileInfo.ShareList = append(dst.FileInfo.ShareList, s.UserId)
		}
	}
	err := proto.Unmarshal(src.Body, dst.Body)
	if err != nil {
		log.Printf("error while unmarshal body: " + err.Error())
	}
	tmd, err := registry.GlobalTableRegistry().TableMetaReg.Find(src.TableId)
	if err != nil {
		log.Printf("error while TableMetaReg.Find(%s): "+err.Error(), src.TableId)
		return
	}
	// 对列进行过滤，可能表格已经没有这一列了
	colList := []*cap.TemplateColumn{}
	for _, col := range dst.Body.Output.VisibleColumns {
		if _, err := tmd.Columns().Find(col.ColumnId); err == nil {
			colList = append(colList, col)
		}
	}
	dst.Body.Output.VisibleColumns = colList
	visibleMap := make(map[string]interface{})
	for i := range dst.Body.Output.VisibleColumns {
		col, err := tmd.Columns().Find(dst.Body.Output.VisibleColumns[i].ColumnId)
		if err != nil {
			log.Printf("error while tmd.Columns().Find(%s.%s): "+err.Error(),
				src.TableId, dst.Body.Output.VisibleColumns[i].ColumnId)
			continue
		}
		dst.Body.Output.VisibleColumns[i].ColumnDetail = col.ToTableColumn()
		// XSG-2380
		// 旧的模板兼容
		if src.FCreateTime.Before(time.Date(2021, 7, 9, 0, 0, 0, 0, time.Local)) {
			dst.Body.Output.VisibleColumns[i].Visible = true
		}
		visibleMap[col.ID] = nil
	}
	// XSG-2380 把看不见的列放底部
	if len(dst.Body.Output.VisibleColumns) == len(tmd.Columns().List()) {
		return
	}
	for _, col := range tmd.Columns().List() {
		if _, ok := visibleMap[col.ID]; !ok {
			dst.Body.Output.VisibleColumns = append(dst.Body.Output.VisibleColumns,
				&cap.TemplateColumn{ColumnId: col.ID, ColumnDetail: col.ToTableColumn(), Visible: col.Required})
		}
	}
	return
}

// FindTemplate ...
func (tm *Manager) FindTemplate(ss *db.Session, id string) (*cap.Template, error) {
	mapper := mysql.NewTableTemplateMapper(ss)
	dbTpl, err := mapper.FindTemplate(id)
	if err != nil {
		return nil, errors.Wrap(err).Log()
	}
	return mapTpl(dbTpl), nil
}

// FindTemplatesByTableAndCreateUser ...
func (tm *Manager) FindTemplatesByTableAndCreateUser(ss *db.Session, tableID, createUser string) ([]*cap.Template, error) {
	mapper := mysql.NewTableTemplateMapper(ss)
	dbTplList, err := mapper.FindTemplates(mysql.FilterTableIDEquals(tableID), mysql.FilterCreateUserEquals(createUser))
	if err != nil {
		return nil, errors.Wrap(err).Log()
	}
	var tplList []*cap.Template
	for _, dbTpl := range dbTplList {
		tplList = append(tplList, mapTpl(dbTpl))
	}
	return tplList, nil
}

// FindTemplatesByTableAndShareUser ...
func (tm *Manager) FindTemplatesByTableAndShareUser(ss *db.Session, tableID, shareUser string) ([]*cap.Template, error) {
	mapper := mysql.NewTableTemplateMapper(ss)
	dbTplList, err := mapper.FindTemplatesByShareUserAndTableID(shareUser, tableID)
	if err != nil {
		return nil, errors.Wrap(err).Log()
	}
	var tplList []*cap.Template
	for _, dbTpl := range dbTplList {
		if dbTpl.FCreateUser == shareUser {
			continue
		}
		tpl := mapTpl(dbTpl)
		tpl.Name += fmt.Sprintf("[%s共享]", tpl.FileInfo.CreateUser.DisplayName)
		tplList = append(tplList, tpl)
	}
	return tplList, nil
}

// FindPublicTemplatesByTable ...
func (tm *Manager) FindPublicTemplatesByTable(ss *db.Session, tableID, currentUserID string) ([]*cap.Template, error) {
	mapper := mysql.NewTableTemplateMapper(ss)
	dbTplList, err := mapper.FindTemplates(mysql.FilterTableIDEquals(tableID),
		mysql.FilterTableAccessEquals(int(cap.FileAccessType_TA_PUBLIC)))
	if err != nil {
		return nil, errors.Wrap(err).Log()
	}
	var tplList []*cap.Template
	for _, dbTpl := range dbTplList {
		if dbTpl.FCreateUser == currentUserID {
			continue
		}
		tpl := mapTpl(dbTpl)
		tpl.Name += fmt.Sprintf("[%s共享]", tpl.FileInfo.CreateUser.DisplayName)
		tplList = append(tplList, tpl)
	}
	return tplList, nil
}

var theManager *Manager

// GlobalManager gets global template manager
func GlobalManager() *Manager {
	return theManager
}

func Init() {
	if theManager == nil {
		theManager = &Manager{}
	}
}
