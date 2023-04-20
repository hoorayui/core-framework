package utils

import "github.com/hoorayui/core-framework/pkg/cap/msg/errors"

// ErrFailedMapValue ...
var ErrFailedMapValue = errors.New("failed to map value (%v) from dt(%s) to vt(%s)")

// ErrFailedParseConditionValue ...
var ErrFailedParseConditionValue = errors.New("failed to parse condition value for(%s): %v")
