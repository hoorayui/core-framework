package errors

import (
	"framework/pkg/cap/msg/errors/handle"
	"framework/pkg/table/data"
	"framework/pkg/table/data/driver/dbdriver"
	"framework/pkg/table/data/utils"
	"framework/pkg/table/mysql"
	"framework/pkg/table/registry"
	"framework/pkg/table/service"
	"framework/pkg/table/template"
)

func init() {
	handle.RegisterDefaultHandler(template.ErrOperatePermissionDenied, ErrOperatePermissionDenied)
	// XSG-3148
	handle.RegisterDefaultHandler(template.ErrDisableSharingFormToAllUsers, ErrDisableSharingFormToAllUsers)
	handle.RegisterDefaultHandler(service.ErrNotFormAction, ErrNotFormAction)
	handle.RegisterDefaultHandler(registry.ErrNodeNotExist, ErrNodeNotExist)
	handle.RegisterDefaultHandler(registry.ErrInvalidColumnID, ErrInvalidColumnID)
	handle.RegisterDefaultHandler(registry.ErrOptionNotFound, ErrOptionNotFound)
	handle.RegisterDefaultHandler(registry.ErrOperatorNotSupported, ErrOperatorNotSupported)
	handle.RegisterDefaultHandler(registry.ErrInvalidConditionValue, ErrInvalidConditionValue)
	handle.RegisterDefaultHandler(registry.ErrInvalidConditionValueType, ErrInvalidConditionValueType)
	handle.RegisterDefaultHandler(registry.ErrAggregateNotSupported, ErrAggregateNotSupported)
	handle.RegisterDefaultHandler(registry.ErrTplNullOutput, ErrTplNullOutput)
	handle.RegisterDefaultHandler(registry.ErrInvalidRowActionID, ErrInvalidRowActionID)
	handle.RegisterDefaultHandler(registry.ErrIllegalArguments, ErrIllegalArguments)
	handle.RegisterDefaultHandler(mysql.ErrRowsAffectedZero, ErrRowsAffectedZero)
	handle.RegisterDefaultHandler(mysql.ErrDuplicateEntry, ErrDuplicateEntry)
	handle.RegisterDefaultHandler(mysql.ErrDataTooLong, ErrDataTooLong)
	handle.RegisterDefaultHandler(dbdriver.ErrInvalidValueForCondition, ErrInvalidValueForCondition)
	handle.RegisterDefaultHandler(dbdriver.ErrUnknownOperator, ErrUnknownOperator)
	handle.RegisterDefaultHandler(dbdriver.ErrOperatorNotSupportedForColumn, ErrOperatorNotSupportedForColumn)
	handle.RegisterDefaultHandler(dbdriver.ErrResultExceedMaxLimit, ErrResultExceedMaxLimit)
	handle.RegisterDefaultHandler(utils.ErrFailedMapValue, ErrFailedMapValue)
	handle.RegisterDefaultHandler(utils.ErrFailedParseConditionValue, ErrFailedParseConditionValue)
	handle.RegisterDefaultHandler(data.ErrDriverNotFoundForTable, ErrDriverNotFoundForTable)
	handle.RegisterDefaultHandler(data.ErrInvalidColumnFromStruct, ErrInvalidColumnFromStruct)
	handle.RegisterDefaultHandler(data.ErrNotResultForID, ErrNotResultForID)
}
