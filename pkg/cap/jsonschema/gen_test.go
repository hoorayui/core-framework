package jsonschema

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	fmt.Println(ToSnakeCase("userName"))
}
