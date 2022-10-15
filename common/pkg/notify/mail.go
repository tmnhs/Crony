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

var mailTemplate = `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <title></title>
    <meta charset="utf-8"/>

</head>
<body>
<div class="cap" style="
            border: 2px solid black;
            background-color: whitesmoke;
            height: 500px"
>
    <div class="content" style="
            background-color: white;
            background-clip: padding-box;
            color:black;
            font-size: 13px;
            padding: 25px 25px;
            margin: 25px 25px;
        ">
        <div class="hello" style="text-align: center; color: #FF3333;font-size: 18px;font-weight: bolder">
            {{.Subject}}
        </div>
        <br>
        <div>
            <table border="1"  bordercolor="black" cellspacing="0px" cellpadding="4px" style="margin: 0 auto;">
                <tr >
                    <td>告警主机</td>
                    <td >{{.IP}}</td>
                </tr>

                <tr>
                    <td>告警时间</td>
                    <td>{{.OccurTime}}</td>
                </tr>

                <tr>
                    <td>告警信息</td>
                    <td>{{.Body}}</td>
                </tr>

            </table>
        </div>
        <br><br>
    </div>
</div>
<br>

</body>
</html>
`

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
	tmpl, err := template.New("notify").Parse(mailTemplate)
	if err != nil {
		return fmt.Sprintf("Failed to parse the notification template error: %s", err.Error())
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, msg)
	if err != nil {
		return fmt.Sprintf("Failed to parse the notification template execute error: %s", err.Error())
	}
	return buf.String()
}
