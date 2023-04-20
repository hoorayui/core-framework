package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zhlicen/converter"
)

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func main() {
	t2g := flag.Bool("t2g", false, "table to go struct, usage: cap -t2g -dsn \"root:root@tcp(localhost:3306)/test?charset=utf8\" -p packageName -of output.go")
	dsn := flag.String("dsn", "", "use with t2g")
	oFile := flag.String("of", "", "use with t2g, output file")
	packageName := flag.String("p", "", " package name")
	tableTag := flag.Bool("tts", false, " generator table tags")
	tableName := flag.String("tn", "", "use with t2g, table name, convert all tables if no tn is specified")
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "cap [OPTION] [ARGUMENTS] \n")
		flag.PrintDefaults()
	}
	// table to go struct
	if *t2g {
		table2Go(*dsn, *oFile, *packageName, *tableName, *tableTag)
		return
	}

	flag.Usage()
}

func table2Go(dsn, file, packageName, tableName string, tableTags bool) {
	if dsn == "" || file == "" {
		panic("invalid dsn or f args")
	}
	cov := converter.NewTable2Struct().Dsn(dsn)
	if packageName != "" {
		cov.PackageName(packageName)
	}
	if tableName != "" {
		cov.Table(tableName)
	}
	cov.TagKey("db")
	cov.RealNameMethod("TableName")
	cov.SavePath("./" + file)
	cov.DateToTime(true)
	cov.EnableTableTags(tableTags)
	cov.Config(&converter.T2tConfig{StructNameToHump: true, GenNullableType: true})
	err := cov.Run()
	if err != nil {
		panic(err)
	}
}
