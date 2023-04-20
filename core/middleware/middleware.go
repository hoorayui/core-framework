package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hoorayui/core-framework/types"
	util2 "github.com/hoorayui/core-framework/util"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// responseWriter 存储响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 读取响应数据
func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GetAllBaseMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		LogMiddleware,
		PromeOpsMiddleware,
		// CorsMiddleware,
	}
}

// PromeOpsMiddleware ops计数
func PromeOpsMiddleware(c *gin.Context) {
	types.OpsProcessedCounter.Inc()
}
func LogMiddleware(c *gin.Context) {
	// request_id
	uuid := util2.NewUUIDString("")
	c.Params = append(c.Params, gin.Param{Key: "request_id", Value: uuid})
	logger := logrus.WithField("request_id", uuid)

	// 响应起始时间
	start := util2.Now()
	// 图片接口不记录请求体
	if c.Request.URL.Path != "/upload-picture" {
		// 获取响应体
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		logger.WithFields(logrus.Fields{
			"request": string(reqBody),
			"url":     c.Request.Host + c.Request.URL.Path,
		}).Infof("request start")
	} else {
		logger.WithFields(logrus.Fields{
			"url": c.Request.Host + c.Request.URL.Path,
		}).Infof("request start")
	}
	responsWriter := &responseWriter{
		body:           bytes.NewBufferString(""),
		ResponseWriter: c.Writer,
	}
	c.Writer = responsWriter

	// 执行业务
	c.Next()

	// 计算响应时间，记录响应结束日志
	cost := time.Since(start).Microseconds()
	// 收集响应时间
	types.ResponseTimeHistogram.Observe(float64(cost) / 1000)
	if cost > 1000000 {
		logger.WithFields(logrus.Fields{
			"cost":     fmt.Sprintf("%.3f s", float32(cost)/1000000),
			"response": responsWriter.body.String(),
		}).Info("request finish")
	} else {
		logger.WithFields(logrus.Fields{
			"cost":     fmt.Sprintf("%.3f ms", float32(cost)/1000),
			"response": responsWriter.body.String(),
		}).Info("request finish")
	}
}

func CorsMiddleware(c *gin.Context) {
	method := c.Request.Method

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
	}
	// 处理请求
	c.Next()
}
