package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Generator struct {
	Nano int
	Rand int
}

var (
	gen   Generator
	mutex sync.Mutex
)

func init() {
	rand.Seed(Now().Unix())
	gen = Generator{}
}

func DD(v interface{}) {
	j, _ := json.MarshalIndent(v, "", "	")
	println(string(j))
}

const (
	DefaultPrefixUUIDString = "default-73a61d26-8c61-4086-899c-9a5959f145fa"
	DefaultUUIDString       = "73a61d26-8c61-4086-899c-9a5959f145fa"
	ShortIDDigits           = "abcdefghijkmnpqrstuvwxyz0123456789"
)

// OldUUIDString : create an uuid string, return <prefix>-<uuid> if prefix is not empty, or just return uuid
func OldUUIDString(prefix string) string {
	needPrefix := prefix != ""
	uuidStr := DefaultUUIDString
	if needPrefix {
		uuidStr = DefaultPrefixUUIDString
	}

	newID := uuid.NewV1()

	uuidStr = newID.String()
	if needPrefix {
		uuidStr = prefix + "-" + uuidStr
	}
	return strings.Replace(uuidStr, "-", "", -1)
}

func UUIDToShortID(UUID string) string {
	// 32uuid -> 32md5 hex
	data := []byte(UUID)
	hash := md5.Sum(data)
	md5str := fmt.Sprintf("%x", hash)

	var result []byte
	for i := 0; i < 16; i++ {
		// parse 2bit char from 16base to 10base
		index, _ := strconv.ParseUint(md5str[2*i:2*i+2], 16, 32)
		result = append(result, ShortIDDigits[index%34])
	}
	return string(result)
}

func (g *Generator) String() string {
	mutex.Lock()
	defer mutex.Unlock()
	tInt := Now().UnixNano()
	if tInt != int64(gen.Nano) {
		gen.Nano = int(tInt)
		g.Rand = rand.Intn(59999)
	} else {
		g.Rand += 1
	}
	return fmt.Sprintf("%d%05d", gen.Nano, gen.Rand)
}

func NewUUIDString(prefix string) string {
	return prefix + gen.String()
}

// ParseBaseUrl 解析访问地址
func ParseBaseUrl(ctx *gin.Context) string {
	if ctx.Request.Referer() != "" {
		referer := ctx.Request.Referer()
		spList := strings.Split(referer, "/")
		if len(spList) >= 3 {
			return strings.Join(spList[0:3], "/")
		}
	} else {
		return "http://" + ctx.Request.Host
	}
	return ""
}
