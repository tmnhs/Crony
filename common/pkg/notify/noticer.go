package notify

import (
	"time"
)

type Noticer interface {
	SendMsg(*Message)
}

type Message struct {
	Type      int
	IP        string
	Subject   string
	Body      string
	To        []string
	OccurTime time.Time
}

var msgQueue chan *Message

func Init(mail *Mail, web *WebHook) {
	_defaultMail = &Mail{
		Port:     mail.Port,
		From:     mail.From,
		Host:     mail.Host,
		Secret:   mail.Secret,
		Nickname: mail.Nickname,
	}
	_defaultWebHook = &WebHook{
		Kind: web.Kind,
		Url:  web.Url,
	}
	msgQueue = make(chan *Message, 64)
}

func Send(msg *Message) {
	msgQueue <- msg
}

func Serve() {
	//for msg := range _defaultMail.msgChan {
	for {
		select {
		case msg := <-msgQueue:
			if msg == nil {
				continue
			}
			switch msg.Type {
			case 1:
				//发送邮件
				go _defaultMail.SendMsg(msg)
			case 2:
				//webhook
				go _defaultWebHook.SendMsg(msg)
			}
		}
	}

}
