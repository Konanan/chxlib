package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/logs"
	"sort"
)

func Sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func Hmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

/**
	method:GET/POST  请求类型
	host：gm-moshi.chuxinhd.com  请求的域名
    uri : /api/login  请求的路由
	secretId： 签名的ID
	secretKey：签名的秘钥
	algorithm：签名类型 CX-HMAC-SHA256  固定写法  CX-HMAC-SHA256
	timestamp：时间戳
**/
func Sign(method, host, uri, payload, secretId, secretKey, algorithm string, timestamp int64) string {

	hashedRequestPayload := Sha256hex(payload)
	canonicalRequest := fmt.Sprintf("%s%s%s%s%s",
		method,
		host,
		uri,
		hashedRequestPayload)

	scope := fmt.Sprintf("%s/%s/chuxin_request", ToString(timestamp), secretId)
	hashedCanonicalRequest := Sha256hex(canonicalRequest)

	string2sign := fmt.Sprintf("%s%s%s%s",
		algorithm,
		ToString(timestamp),
		scope,
		hashedCanonicalRequest)

	secretTime := Hmacsha256(ToString(timestamp), "chuxin"+secretKey)
	secretUri := Hmacsha256(uri, secretTime)
	secretSigning := Hmacsha256("chuxin_request", secretUri)
	signature := hex.EncodeToString([]byte(Hmacsha256(string2sign, secretSigning)))
	return signature
}

/**
	secretId： 签名的ID
	secretKey：签名的秘钥
	payload: 请求参数的json数据
	algorithm：签名类型 CX-HMAC-SHA256  固定写法  CX-HMAC-SHA256
	timestamp：时间戳
**/
func Sign2(secretId, secretKey, algorithm string, timestamp int64, params map[string]interface{}) string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(params))
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i := range keys {
		k := keys[i]
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(ToString(params[k]))
		buf.WriteString("&")
	}
	buf.Truncate(buf.Len() - 1)

	buf.WriteString(secretId + algorithm + ToString(timestamp))
	hashed := hmac.New(sha1.New, []byte(secretKey))
	hashed.Write(buf.Bytes())
	sign := base64.StdEncoding.EncodeToString(hashed.Sum(nil))
	logs.Info("生成签名 buf:%s sign:%s", buf.String(), sign)
	return sign
}
