package registry

import (
	"fmt"
	"testing"

	cap "framework/pkg/table/proto"
)

func TestLoadOptionFromProtoEnum(t *testing.T) {
	LoadOptionFromProtoEnum(cap.ValueType(0))
}

func TestRegisterOptionFromProtoEnum(t *testing.T) {
	err := RegisterOptionFromProtoEnum(cap.ValueType(0))
	if err != nil {
		t.Fatal(err)
	}
	opts, err := GlobalTableRegistry().OptionReg.GetOptions("cap.ValueType")
	if err != nil {
		t.Fatal(err)
	}
	for _, opt := range opts {
		fmt.Println(opt.Id, opt.Name)
	}
}
