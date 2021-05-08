package gtools

import (
	"regexp"
)

// VerifyEmail
func VerifyEmail(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyMobile,国内手机号
func VerifyMobile(mobile string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}

func VerifyBirthday(birthday string) bool {
	patter := `^\d{4}-\d{2}-\d{2}`
	reg := regexp.MustCompile(patter)
	return reg.MatchString(birthday)
}
