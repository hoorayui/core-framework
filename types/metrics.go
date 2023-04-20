package types

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ops
	OpsProcessedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	// 事务数量
	TXCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_dbtx_total",
		Help: "db transaction rollback total number",
	})
	// 事务回滚数量
	TXRollbackCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_dbtx_rollback_total",
		Help: "db transaction rollback number",
	})
	// 事务失败数量
	TXFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_dbtx_failed_total",
		Help: "db transaction failed number",
	})
	ResponseTimeHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "myapp_response_time_milliseconds",
		Help: "api response time(ms)",
	})
	// promauto.NewCounterFunc(prometheus.CounterOpts{
	// 	Name: "myapp_processed_ops_total",
	// 	Help: "The total number of processed events",
	// }, func() float64 {
	// 	return 1
	// })
)
