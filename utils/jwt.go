package utils

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"hy.juck.com/go-publisher-server/config"
)

var (
	G = config.G
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var MySecret []byte

// 获取服务器唯一标识
func GetServerUUid() (string, error) {
	// linux下获取唯一标识命令
	var command string
	var cmd *exec.Cmd
	sysType := runtime.GOOS
	switch sysType {
	case "linux":
		command = `dmidecode -s system-uuid | tr 'A-Z' 'a-z'`
		cmd = exec.Command("bash", "-c", command)
	case "windows":
		command = "csproduct get UUID"
		cmd = exec.Command(command)
	default:
		return "jwt_token", nil
	}
	// 开始执行脚本
	out, err := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(out))
	if err != nil {
		return "", err
	}
	return outStr, nil
}

// GenToken 生成token
func GenToken(username string) (string, error) {
	uuId, err := GetServerUUid()
	if err != nil {
		return "", err
	}
	MySecret = []byte(uuId)
	c := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(G.C.Jwt.Token.Expire) * time.Second).Unix(), // 过期时间
			Issuer:    username,                                                                 // 签发人
			Subject:   username,                                                                 // 签发对象
			NotBefore: time.Now().Unix(),                                                        // 最早使用时间
			IssuedAt:  time.Now().Unix(),                                                        // 签发时间
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	uuId, err := GetServerUUid()
	if err != nil {
		return nil, err
	}

	MySecret = []byte(uuId)
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
