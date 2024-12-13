package prometheus

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
	"testing"
)

func TestSendEmail(t *testing.T) {
	err := sendEmail("332192810@qq.com", "liming2300705@163.com", "邮件测试", "hello world")

	if err != nil {
		t.Fatal(err)
	}
}

const (
	Host     = "smtp.qq.com"
	Port     = 587
	UserName = "332192810@qq.com"
	PSW      = "qwjjgtnypfhmbidj"
)

func sendEmail(from, to, subject, body string) error {
	// 准备邮件内容
	msg := buildMessage(from, to, subject, body)

	auth := smtp.PlainAuth("", UserName, PSW, Host)

	// 发送邮件
	err := smtp.SendMail(fmt.Sprintf("%s:%d", Host, Port), auth, from, []string{to}, []byte(msg))
	return err
}

func buildMessage(from, to, subject, body string) string {

	// 邮件内容
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\n%s\n\n%s",
		from, to, subject, "Content-Type: text/plain; charset=UTF-8", body)

	return msg
}

// 编码邮件头部
func encodeAddress(addr string) string {
	// 使用 MIME 编码 RFC2047
	parsedAddr, err := mail.ParseAddress(addr)
	if err != nil {
		log.Fatal("地址解析失败: ", err)
	}
	// 如果地址中有非 ASCII 字符，需要用 base64 编码
	if strings.ContainsAny(parsedAddr.Name, "非ASCII字符") {
		encodedName := base64Encode(parsedAddr.Name)
		return fmt.Sprintf("=?UTF-8?B?%s?= <%s>", encodedName, parsedAddr.Address)
	}
	return addr
}

// 编码邮件主题
func encodeSubject(subject string) string {
	// 使用 RFC 2047 编码主题
	encodedSubject := base64Encode(subject)
	return fmt.Sprintf("=?UTF-8?B?%s?=", encodedSubject)
}

// Base64 编码
func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
