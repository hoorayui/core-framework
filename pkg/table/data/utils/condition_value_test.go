package utils

import (
	"fmt"
	"reflect"
	"testing"

	cap "github.com/hoorayui/core-framework/pkg/table/proto"
	"github.com/hoorayui/core-framework/pkg/table/registry"
)

func TestParseConditionValue(t *testing.T) {
	v, err := ParseConditionValue(&registry.TableColumnDescriptor{
		DataType:    reflect.TypeOf(0),
		ValueType:   cap.ValueType_VT_STRING,
		ValueFormat: "%03d",
	}, &cap.FilterValue{
		LiteralValues: &cap.Value{V: &cap.Value_VString{VString: "001"}},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
