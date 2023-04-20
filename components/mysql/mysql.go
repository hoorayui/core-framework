package mysql

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hoorayui/core-framework/util/flag"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

var (
	rawDB     *gorm.DB
	debugMode bool
)

type DB struct {
	db *gorm.DB
}

func NewDB() *DB {
	if debugMode {
		return &DB{
			db: rawDB.Debug(),
		}
	}
	return &DB{
		db: rawDB,
	}
}

func NewTX() *DB {
	if debugMode {
		return &DB{
			db: rawDB.Debug().Begin(),
		}
	}
	return &DB{
		db: rawDB.Begin(),
	}
}

func (b *DB) GetDB() *gorm.DB {
	if b.db == nil {
		return rawDB
	}
	return b.db
}

func (b *DB) SetDB(db *gorm.DB) {
	b.db = db
}

func (b *DB) Begin() *gorm.DB {
	b.db = rawDB.Begin()
	return b.db
}

func (b *DB) Commit() {
	if b.db != nil {
		b.db.Commit()
	}
}
func mysqlLogger() logger.Interface {
	// TODO 替换为应用名称
	path := fmt.Sprintf("%s/db.log", flag.LogDir)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if nil != err {
		logrus.Fatalf("打开日志文件[%s]失败1", path)
		return nil
	}
	writer := io.MultiWriter(os.Stdout, f)
	return logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Second / 5, // 慢 SQL 阈值
			LogLevel:                  logger.Info,     // 日志级别
			IgnoreRecordNotFoundError: true,            // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,            // 彩色打印
		},
	)
}
