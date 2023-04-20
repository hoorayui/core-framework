package util

import "time"

//Now 生成东八区时间
func Now() time.Time {
	// 东八区
	cstZone := time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone)
}
