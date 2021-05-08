package gtools

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
)

//md5方法
func Md5(s string) string {
	return Md5Bytes([]byte(s))
}

func Md5Bytes(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(data string) string {
	shaFn := sha1.New()
	shaFn.Write([]byte(data))
	return hex.EncodeToString(shaFn.Sum([]byte("")))
}

func HmacSha256(date, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(date))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func Sha256(data string) string {
	shaFn := sha256.New()
	shaFn.Write([]byte(data))
	return hex.EncodeToString(shaFn.Sum([]byte("")))
}

func HmacSHA1(key string, data string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// Guid 方法
func Guid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

func PasswordEncrypt(password string, salt int64) string {
	saltStr := strconv.FormatInt(salt, 10)
	md51 := Md5(password + "$@H0me" + saltStr)
	md52 := Md5(saltStr)

	result1 := SubString(md52, 24, 8)
	result2 := SubString(md51, 0, 24)

	return Md5(result1 + result2)
}

func Stringify(obj interface{}) string {
	type stringer interface {
		String() string
	}

	switch obj := obj.(type) {
	case stringer:
		return obj.String()

	case string:
		return obj

	case int:
		return Int2Str(obj)

	case int64:
		return Int642Str(obj)

	case float64:
		return Float642Str(obj)

	case bool:
		if obj {
			return "true"
		}
		return "false"

	default:
		log.Println(fmt.Sprintf("[Stringify] no match, obj: %#v", obj))
		return ""
	}
}

// BaiduSignature 空参数值的参数名参与签名,来源:百度开放平台
func BaiduSignature(params map[string]interface{}, secret string) string {
	paramLen := len(params)
	if paramLen <= 0 {
		log.Println(fmt.Sprintf("[BaiduSignature] params len is 0, params: %v", params))
		return ""
	}

	condBox := make([]string, paramLen)
	var i int = 0
	for k, _ := range params {
		condBox[i] = k
		i++
	}

	// 按字典序列排序
	sort.Strings(condBox)

	str := "" // 待签名字符串
	for i = 0; i < paramLen; i++ {
		key := condBox[i]
		str += fmt.Sprintf("%s=%s&", key, Stringify(params[key]))
	}
	str += secret
	//log.Println("[BaiduSignature] need signature str:", str)

	return Md5(str)
}
