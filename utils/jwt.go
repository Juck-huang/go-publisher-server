package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"hy.juck.com/go-publisher-server/config"
	"time"
)

var (
	G = config.G
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var MySecret = []byte(G.C.Jwt.Token.Secret)

// GenToken 生成token
func GenToken(username string) (string, error) {
	c := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(G.C.Jwt.Token.Expire) * time.Second).Unix(), // 过期时间
			Issuer:    username,                                                                 // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (any, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, err
	}
	return nil, errors.New("请先登录或者token已经失效")
}
