package Utility

import (
	"OBPkg/Config"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ymzuiku/hit"
	"net/smtp"
	"strings"
)

func EmailSendConfirmation(target string, token string, expiry string) error {
	subject := "Hmmsim 线路库 账号激活验证信息"
	linkA := "https://reg.bvecs.zbx1425.cn/#/user/activate/" + token
	body := fmt.Sprintf(
		"<h2>多谢您使用 Hmmsim 线路库!</h2>"+
			"<p>%s :</p>"+
			"<p>请用以下链接激活您的账号:<br/>"+
			"<a href=\"%s\">%s</a></p>"+
			"<p>请注意此链接仅在 %s 前有效。</p>"+
			"<p>顺颂时祺</p>", target, linkA, linkA, expiry)
	return send(target, subject, body)
}

func EmailSendFileUpdate(target string, validated bool, fileID uint, pkgID uint, pkgName string, reason string) error {
	result := hit.If(validated, "通过", "退回").(string)
	subject := "Hmmsim 线路库 文件" + result + "通知"
	body := fmt.Sprintf(
		"<h2>%s</h2>"+
			"<p>您好：<br>非常抱歉让您久等了。管理员已验证处理您的以下封包：</p>"+
			"<p>%s: <a href='https://reg.bvecs.zbx1425.cn/#/package/edit/%d'>https://reg.bvecs.zbx1425.cn/#/package/edit/%d</a></p>"+
			"<p>"+
			"处理结果: %s<br>"+
			"文件编号: #%d<br>"+
			"所属封包: %s<br>"+
			"原因: %s"+
			"</p>"+
			"<p>如有任何问询事项，请与管理人员联系。本邮件自动发送，请勿原地址回信。</p>"+
			"<p>顺颂时祺</p>", subject, pkgName, pkgID, pkgID, result, fileID, pkgName, reason)
	return send(target, subject, body)
}

type loginAuth struct {
	username, password string
}

// LoginAuth is used for smtp login auth
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func add76crlf(msg string) string {
	var buffer bytes.Buffer
	for k, c := range strings.Split(msg, "") {
		buffer.WriteString(c)
		if k%76 == 75 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}

func utf8Split(utf8string string, length int) []string {
	resultString := []string{}
	var buffer bytes.Buffer
	for k, c := range strings.Split(utf8string, "") {
		buffer.WriteString(c)
		if k%length == length-1 {
			resultString = append(resultString, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		resultString = append(resultString, buffer.String())
	}
	return resultString
}

func encodeSubject(subject string) string {
	var buffer bytes.Buffer
	buffer.WriteString("Subject:")
	for _, line := range utf8Split(subject, 13) {
		buffer.WriteString(" =?utf-8?B?")
		buffer.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		buffer.WriteString("?=\r\n")
	}
	return buffer.String()
}

func send(target string, subject string, body string) error {
	auth := LoginAuth(
		Config.CurrentConfig.SMTP.Username,
		Config.CurrentConfig.SMTP.Password)

	var header bytes.Buffer
	header.WriteString("From: " + Config.CurrentConfig.SMTP.Sender + "\r\n")
	header.WriteString("To: " + target + "\r\n")
	header.WriteString(encodeSubject(subject))
	header.WriteString("MIME-Version: 1.0\r\n")
	header.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")

	var message bytes.Buffer
	message = header
	message.WriteString("\r\n")
	message.WriteString(body)

	fmt.Println(message.String())
	err := smtp.SendMail(Config.CurrentConfig.SMTP.Host, auth,
		Config.CurrentConfig.SMTP.Username, []string{target}, []byte(message.String()))
	return err
}
