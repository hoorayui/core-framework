package table

import "github.com/zhlicen/converter"

type TableConfig struct {
	DSN         string
	OFile       string
	PackageName string
	TableName   string
	TableTag    bool
}

func StartTableService(cfg TableConfig) {
	table2Go(cfg.DSN, cfg.OFile, cfg.PackageName, cfg.TableName, cfg.TableTag)
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
