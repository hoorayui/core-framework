package context

import (
	"github.com/gin-gonic/gin"
	"github.com/hoorayui/core-framework/components/config"
	"github.com/hoorayui/core-framework/components/log"
	"github.com/hoorayui/core-framework/components/mysql"
	"github.com/hoorayui/core-framework/types"
	"github.com/hoorayui/core-framework/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"jihulab.com/ycwjck/metal-backend/api/apimodel"
)

type Context struct {
	db                *gorm.DB      // 数据库实例
	transactionEnable bool          // 事务开启状态
	logger            *logrus.Entry // 日志
	Ctx               *gin.Context  // gin context
	*User                           // 当前用户状态
}
type User struct {
	TokenUUid   string
	IsLogin     bool
	UserID      string
	UserName    string
	DisplayName string
}

// DB 获取数据库实例
func (c *Context) DB() *gorm.DB {
	if c.db == nil {
		return mysql.NewDB().GetDB()
	}
	return c.db
}

// NewDB 获取数据库独立新实例，不受DB()实例的事务影响
func (c *Context) NewDB() *gorm.DB {
	return mysql.NewDB().GetDB()
}

// SetDB 设置数据库
func (c *Context) SetDB(db *gorm.DB) *Context {
	c.db = db
	return c
}

// Logger 日志
func (c *Context) Logger() *logrus.Entry {
	if c.logger == nil {
		c.logger = logrus.NewEntry(log.GetLogger().Logger)
	}
	return c.logger
}

// SetLogger 设置日志
func (c *Context) SetLogger(e *logrus.Entry) *Context {
	c.logger = e
	return c
}

// Begin 开启事务
func (c *Context) Begin() *Context {
	if c.transactionEnable {
		c.logger.Errorf("重复开启事务:%v", &c.db)
		return c
	}
	c.logger.Infof("开启事务:%v", &c.db)
	c.SetDB(c.db.Begin())
	types.TXCounter.Inc()
	c.transactionEnable = true
	return c
}

// Commit true/空 提交事务，若参数为false 则为回滚
func (c *Context) Commit(ifSuccess ...bool) *Context {
	if !c.transactionEnable {
		c.logger.Errorf("非事务提交:%v", &c.db)
		return c
	}
	c.transactionEnable = false
	var txSuccess = true
	var err error
	if len(ifSuccess) > 0 {
		txSuccess = ifSuccess[0]
	}
	if txSuccess {
		err = c.db.Commit().Error
		c.logger.Infof("事务提交:%v", err)
	} else {
		types.TXRollbackCounter.Inc()
		err = c.db.Rollback().Error
		c.logger.Infof("事务回滚:%v", &c.db)
	}
	if nil != err {
		types.TXFailedCounter.Inc()
		c.logger.Errorf("事务提交/回滚失败:%v", &c.db)
	} else {
		c.logger.Infof("事务结束:%v", &c.db)
	}
	// 设置新db，已回滚事务会影响下次事务执行
	c.SetDB(mysql.NewDB().GetDB())
	return c
}

// RollbackIfNotCommited 若为有效事务且未提交，直接回滚
func (c *Context) RollbackIfNotCommited() {
	if c.transactionEnable {
		c.Commit(false)
	}
}

// Set 上下文件设置值
func (c *Context) Set(key string, value any) {
	c.Ctx.Set(key, value)
}

// Get 上下文件获取值
func (c *Context) Get(key string) (value any, exists bool) {
	return c.Ctx.Get(key)
}

// IsAdmin 是否为管理员
func (c *Context) IsAdmin() bool {
	return c.UserID == "account_admin001"
}

// NewContext new api context
func NewContext(ctx ...*gin.Context) *Context {
	var c = &Context{
		User: &User{},
	}

	if len(ctx) > 0 {
		c.Ctx = ctx[0]
		value, exist := ctx[0].Get("profile")
		if exist && value != nil {
			AccessDetails := value.(*apimodel.AccessDetails)
			c.User = &User{
				IsLogin:     true,
				TokenUUid:   AccessDetails.TokenUuid,
				UserID:      AccessDetails.UserId,
				UserName:    AccessDetails.UserName,
				DisplayName: AccessDetails.DisplayName,
			}
		}
	} else {
		c.Ctx = &gin.Context{}
	}

	// debug模式下 使用管理员身份
	if !c.IsLogin && config.GetConfig().Default.Server.DebugMode {
		c.User = &User{
			UserID:   "account_admin001",
			UserName: "admin",
		}
	}
	c.SetLogger(log.GetLogger().Logger.
		WithField("request_id", c.Ctx.Params.ByName("request_id")).
		WithField("token", c.TokenUUid))
	c.SetDB(mysql.NewDB().GetDB())
	return c
}
func NewWebSocketContext(c ...*gin.Context) *Context {
	ctx := NewContext(c...)
	ctx.SetLogger(ctx.Logger().WithField("websocket-id", ctx.Ctx.Request.Header.Get("Sec-WebSocket-Key")))
	return ctx
}

// NewMockContext new mock api context
func NewMockContext(ctx ...*gin.Context) *Context {
	var c *Context
	if len(ctx) > 0 {
		c = &Context{Ctx: ctx[0]}
	} else {
		c = &Context{Ctx: &gin.Context{}}
	}
	uuid := util.NewUUIDString("")
	c.Ctx.Params = append(c.Ctx.Params, gin.Param{Key: "request_id", Value: uuid})
	c.SetLogger(log.GetLogger().Logger.WithField("request_id", uuid))
	c.SetDB(mysql.NewDB().GetDB())
	// mock 使用管理员身份
	c.User = &User{
		IsLogin:  true,
		UserID:   "account_admin001",
		UserName: "admin",
	}
	return c
}
