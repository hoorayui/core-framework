package main

import "testing"

func Test_table2Go(t *testing.T) {
	table2Go("root:AAbb1234@yc@tcp(121.40.118.206:10001)/5metal?charset=utf8&parseTime=True&loc=Local", "test.go", "mysql", "", true)
}
