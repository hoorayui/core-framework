package jwt

import "testing"

func TestGenToken(t *testing.T) {
	got, err := GenToken(1, "张三")
	if err != nil {
		t.Logf("获取JWT Token失败")
	}
	t.Logf(got)
}

func TestParseToken(t *testing.T) {
	str, err := GenToken(1, "张三")
	claims, err := ParseToken(str)
	if err != nil {
		t.Logf("解析JWT Token失败")
	}
	t.Logf("\n用户ID:[%v]\n用户名:[%v]\nStandardClaims:[%+v]\n", claims.UserID, claims.Username, claims.StandardClaims)
}

func TestCheckToken(t *testing.T) {
	str, err := GenToken(1, "张三")
	if err != nil {
		t.Logf("获取JWT Token失败")
	}

	userName, valid := CheckToken(str)
	t.Logf("用户名:[%v]\n是否有效:[%v]\n", userName, valid)
}
