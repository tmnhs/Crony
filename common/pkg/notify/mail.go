package notify

import (
	"fmt"
	"github.com/go-gomail/gomail"
)

type Noticer interface {
	Send(*Message)
}

type Message struct {
	Subject string
	Body    string
	To      []string
}

var _defaultMail *Mail

func Init(port int, from, host, secret, nickName string) {
	_defaultMail = &Mail{
		Port:     port,
		From:     from,
		Host:     host,
		Secret:   secret,
		Nickname: nickName,
		msgChan:  make(chan Message, 100),
	}
}

type Mail struct {
	Port     int
	From     string
	Host     string
	Secret   string
	Nickname string
	msgChan  chan Message
}

func Serve() {
	var err error
	m := gomail.NewMessage()
	for msg := range _defaultMail.msgChan {
		fmt.Println("dd")
		m.SetHeader("To", msg.To...)
		m.SetHeader("Subject", msg.Subject)
		m.SetBody("text/html", msg.Body)
		m.SetHeader("From", m.FormatAddress(_defaultMail.Nickname, "crony定时任务平台")) //这种方式可以添加别名，即“XX官方”
		//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
		//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果

		d := gomail.NewDialer(_defaultMail.Host, _defaultMail.Port, _defaultMail.From, _defaultMail.Secret)

		if err = d.DialAndSend(m); err != nil {
			//logger.GetLogger().Warn(fmt.Sprintf("smtp send msg[%+v] err: %s", msg, err.Error()))
			fmt.Printf("smtp send msg[%+v] err: %s", msg, err.Error())
		}
	}
}

func Send(msg Message) {
	_defaultMail.msgChan <- msg
}
