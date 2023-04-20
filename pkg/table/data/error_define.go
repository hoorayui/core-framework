package data

import "framework/pkg/cap/msg/errors"

// ErrDupplicateDriverForTable ...
var ErrDupplicateDriverForTable = errors.New("dupplicate driver for table(%s)")

// ErrDriverNotFoundForTable ...
var ErrDriverNotFoundForTable = errors.New("driver not found for table(%s)")

// ErrInvalidColumnFromStruct ...
var ErrInvalidColumnFromStruct = errors.New("invalid column(%s) from struct(%v)")

// ErrTableColumnNotLinkable ...
var ErrTableColumnNotLinkable = errors.New("table(%s.%s) can not be linked, column must support builtin.EQ or builtin.IN")

// ErrNotResultForID ...
var ErrNotResultForID = errors.New("no result for id(%s)")
