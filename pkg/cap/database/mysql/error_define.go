package mysql

import "errors"

// ErrDeleteMustContainFilters ...
var ErrDeleteMustContainFilters = errors.New("delete operation MUST contain filters")

// ErrNoSessionInCtx no session in context
var ErrNoSessionInCtx = errors.New("no session in context")

// ErrSessionTimeout session timeout
var ErrSessionTimeout = errors.New("session timeout")
