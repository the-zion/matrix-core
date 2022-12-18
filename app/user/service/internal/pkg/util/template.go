package util

import "fmt"

func GetPhoneTemplate(template string) string {
	var id string
	switch template {
	case "1":
		id = "SMS_264780342"
	case "2":
		id = "SMS_264740382"
	case "3":
		id = "SMS_264785340"
	case "4":
		id = "SMS_264775391"

	}
	return id
}

func GetEmailTemplate(template, code string) string {
	var content string
	switch template {
	case "1":
		content = code + " 为您的登录验证码，请于2分钟内填写，如非本人操作，请忽略本信息。"
	case "2":
		content = fmt.Sprintf("您的动态验证码为：%s， 2分钟内有效！您正在进行密码重置操作，如非本人操作，请忽略本信息！", code)
	case "3":
		content = fmt.Sprintf("您正在进行账号绑定，验证码为：%s，2分钟有效，为保障帐户安全，请勿向任何人提供此验证码。", code)
	case "4":
		content = fmt.Sprintf("您正在进行账号解绑操作，验证码为：%s，2分钟有效，为保障帐户安全，请勿向任何人提供此验证码", code)
	case "5":
		content = code + " 为您的邮箱注册验证码，请于2分钟内填写，如非本人操作，请忽略本信息。"
	}
	return content
}
