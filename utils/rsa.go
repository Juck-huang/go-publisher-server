package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type Rsa struct {
	privateKey []byte
}

func NewRsa(privateKey string) *Rsa {
	var privateKeyAll = []byte(fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----", privateKey))
	return &Rsa{
		privateKey: privateKeyAll,
	}
}

// Decrypt rsa解密
func (o *Rsa) Decrypt(ciphertext []byte) ([]byte, error) {
	code, _ := base64.StdEncoding.DecodeString(string(ciphertext))
	block, _ := pem.Decode(o.privateKey)
	if block == nil {
		return nil, errors.New("私钥key错误")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("x509解析错误:" + err.Error())
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, code)
}
