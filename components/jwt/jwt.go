package jwt

import (
	"errors"
	"framework/util"
	"time"

	"github.com/golang-jwt/jwt"
)

const TokenExpireDuration = time.Minute * 15 // 两小时过期

var privateKey = []byte("M@O8VWUb81YvmtWLHGB2I_V7di5-@0p(MF*GrE!sIws23F")

type myClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (*myClaims) Valid() error {
	return nil
}

// GetToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己的声明数据
	c := &myClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: util.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "fivegold",                                 // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(privateKey)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*myClaims, error) {
	mc := new(myClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

func CheckToken(tokenString string) (string, bool) {
	if tokenString == "" || len(tokenString) <= 0 {
		return "", false
	}

	mc := new(myClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		return "", false
	}

	if !token.Valid {
		return "", false
	}

	now := util.Now()
	if mc.ExpiresAt <= now.Unix() {
		return mc.Username, false
	} // 过期

	return mc.Username, true
}
