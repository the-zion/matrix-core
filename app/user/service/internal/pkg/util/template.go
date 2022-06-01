package util

func GetPhoneTemplate(template string) string {
	var id string
	switch template {
	case "1":
		id = "1365683"
	case "2":
		id = "1365681"
	case "3":
		id = "1178113"
	case "4":
		id = "1372638"

	}
	return id
}

func GetEmailTemplate(template, code string) string {
	var content string
	switch template {
	case "1":
		content = code + " 为您的登录验证码，请于2分钟内填写，如非本人操作，请忽略本信息。"
	case "2":
		content = "您的动态验证码为：" + code + "， 2分钟内有效！您正在进行密码重置操作，如非本人操作，请忽略本信息！"
	case "3":
		content = "您正在修改注册手机号码，验证码为：" + code + "，2分钟有效，为保障帐户安全，请勿向任何人提供此验证码"
	case "4":
		content = "您正在修改注册邮箱号码，验证码为：" + code + "，2分钟有效，为保障帐户安全，请勿向任何人提供此验证码"
	}
	return content
}
