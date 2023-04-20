package service

import "framework/pkg/cap/msg/errors"

// ErrNotFormAction ...
var ErrNotFormAction = errors.New("action(%s) is not a form action")

// ErrTempTableNotSupportTemplateOp ...
var ErrTempTableNotSupportTemplateOp = errors.New("temp table not support template operation")
