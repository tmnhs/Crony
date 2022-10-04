package notify

import (
	"bytes"
	"fmt"
	"github.com/go-gomail/gomail"
	"html/template"
)

const (
	NotifyTypeMail    = 1
	NotifyTypeWebHook = 2
)

var _defaultMail *Mail

type Mail struct {
	Port     int
	From     string
	Host     string
	Secret   string
	Nickname string
}

func (mail *Mail) SendMsg(msg *Message) {
	m := gomail.NewMessage()

	//邮件
	m.SetHeader("From", m.FormatAddress(_defaultMail.From, _defaultMail.Nickname)) //这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	msgData := parseMailTemplate(msg)
	m.SetBody("text/html", msgData)

	d := gomail.NewDialer(_defaultMail.Host, _defaultMail.Port, _defaultMail.From, _defaultMail.Secret)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		//logger.GetLogger().Warn(fmt.Sprintf("smtp send msg[%+v] err: %s", msg, err.Error()))
	}
}

func parseMailTemplate(msg *Message) string {
	tmpl, err := template.ParseFiles("./notify.html")
	if err != nil {
		return fmt.Sprintf("解析通知模板失败 ParseFiles: %s", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, msg)
	if err != nil {
		return fmt.Sprintf("解析通知模板失败 Execute: %s", err)
	}
	return buf.String()
}
