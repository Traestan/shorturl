package endpoint

import (
	"net/smtp"
	"os"

	"bytes"
	"fmt"
	"html/template"

	"github.com/Traestan/shorturl/repository"
	"github.com/go-kit/kit/log"
)

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func MailSend(to, subject, body string) {

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger.Log("msg", "hello")

	templateData := struct {
		Name string
		Body string
	}{
		Name: "Test",
		Body: body,
	}

	r := NewRequest([]string{to}, subject, body)
	err := r.ParseTemplate("template.html", templateData)
	if err != nil {
		logger.Log("err", err)
	}
	if err := r.ParseTemplate("template.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
	}

}

func MailSendRegister(user *repository.User, subject, body string) {

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger.Log("msg", "hello")

	templateData := struct {
		Name string
		Body string
		//Pass string
	}{
		Name: user.Email,
		Body: body,
		//Pass: user.Password,
	}

	r := NewRequest([]string{user.Email}, subject, body)
	err := r.ParseTemplate("emails/register.html", templateData)
	if err != nil {
		logger.Log("err", err)
	}
	if err := r.ParseTemplate("emails/register.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		logger.Log("msg", fmt.Sprintf("Emain send %v", ok))
	}
}

func MailSendForgot(user *repository.User, subject, body string) {

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger.Log("msg", "hello")

	templateData := struct {
		Name string
		Body string
		//Pass string
	}{
		Name: user.Email,
		Body: body,
		//Pass: user.Password,
	}

	r := NewRequest([]string{user.Email}, subject, body)
	err := r.ParseTemplate("emails/forgot.html", templateData)
	if err != nil {
		logger.Log("err", err)
	}
	if err := r.ParseTemplate("emails/forgot.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		logger.Log("msg", fmt.Sprintf("Emain forgot send %v", ok))
	}
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	var auth smtp.Auth

	addr := "smtp.yandex.ru:587"
	from := "info@cloza.ru"
	pass := "Cloza25"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	frommime := "From: " + from + "\n" + "To: " + r.to[0] + "\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + frommime + mime + "\n" + r.body)

	auth = smtp.PlainAuth("", from, pass, "smtp.yandex.ru")

	err := smtp.SendMail(addr, auth, "info@cloza.ru", r.to, msg)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
