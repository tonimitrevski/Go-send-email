package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/smtp"
	"bytes"
	"html/template"
	"strconv"
)

var auth smtp.Auth

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	auth = smtp.PlainAuth("", "d9eaf1c8f53e3a", "241e6d8be43f10", "smtp.mailtrap.io")
	templateData := struct {
		Name string
		URL  string
	}{
		Name: "Toni",
		URL:  "http://stativa.space",
	}
	r := NewRequest([]string{"mitrevski@mail.com"}, "Hello Toni!", "Hello, World!")
	err := r.ParseTemplate("bin/templates/template.html", templateData)
	if err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
		return events.APIGatewayProxyResponse{Body: strconv.FormatBool(ok), StatusCode: 200}, nil
	}
	return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
}

func main() {
	lambda.Start(Handler)
}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.mailtrap.io:465"

	if err := smtp.SendMail(addr, auth, "toni@stativa.com.mk", r.to, msg); err != nil {
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