package Utility

import (
	"OBPkg/Config"
	"errors"
	"fmt"
	"net/smtp"
)

func EmailSendConfirmation(target string, token string, expiry string) error {
	auth := LoginAuth(
		Config.CurrentConfig.SMTP.Username,
		Config.CurrentConfig.SMTP.Password)
	subject := "BVE Content Service - Email Confirmation"
	linkA := "https://bvecs.tk/#/user/activate/" + token
	body := fmt.Sprintf(
		"<h2>Thank you for joining the BVE Content Service platform!</h2>"+
			"<p>%s :</p>"+
			"<p>Please use the following link to activate your account:<br/>"+
			"<a href=\"%s\">%s</a></p>"+
			"<p>Please note that this link is only valid before %s.</p>"+
			"<p>Best Regards</p>", target, linkA, linkA, expiry)
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s",
		target,
		Config.CurrentConfig.SMTP.Username,
		subject,
		"Content-Type: text/html; charset=UTF-8",
		body))
	fmt.Println(fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s",
		target,
		Config.CurrentConfig.SMTP.Username,
		subject,
		"Content-Type: text/html; charset=UTF-8",
		body))
	err := smtp.SendMail(Config.CurrentConfig.SMTP.Host, auth,
		Config.CurrentConfig.SMTP.Username, []string{target}, msg)
	return err
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
