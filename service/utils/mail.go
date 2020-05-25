package utils

import (
	"lchat/service/conf"
	"net/smtp"
	"regexp"
	"strings"
)

// 校验邮箱格式
func VerifyEmailFormat(email string) bool {
	pattern := "\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*"
	b, _ := regexp.Match(pattern, []byte(email))
	return b
}

// 发送电子邮件
func SendToMail(user, password, addr, to, subject, body, mailtype string) error {
	host := strings.Split(addr, ":")[0]
	auth := smtp.PlainAuth("", user, password, host)
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/html; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	return smtp.SendMail(addr, auth, user, sendTo, msg)
}

func SendEmailRegisterCode(email, code string) error {
	subject := "lchat 邮箱号注册"
	body := "lchat <" + email + ">邮箱号注册验证码: " + code
	return SendToMail(conf.Get().Mail.User, conf.Get().Mail.Password, conf.Get().Mail.Addr, email, subject, body, "")
}
