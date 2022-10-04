package notify

import (
	"testing"
	"time"
)

//https://open.feishu.cn/open-apis/bot/v2/hook/184aa608-e62c-4d49-9621-d8fb22ed1254
func TestSend(t *testing.T) {
	Init(&Mail{
		Port:     465,
		From:     "1685290935@qq.com",
		Host:     "smtp.qq.com",
		Secret:   "wqkncnzhdvbpdabj",
		Nickname: "test",
	}, &WebHook{
		Url:  "https://open.feishu.cn/open-apis/bot/v2/hook/4df45fb8-6458-4af9-a2d5-6c0396895743",
		Kind: "feishu",
	})
	go Serve()
	Send(&Message{
		IP:      "127.0.0.1",
		Type:    1,
		Subject: "test subject",
		Body:    "test body select",
		To:      []string{"mananhai@highlight.mobi"},
	})
	Send(&Message{
		IP:      "127.0.0.1",
		Type:    2,
		Subject: "test2 subject",
		Body:    "test2 body select",
		To:      []string{"mananhai@highlight.mobi"},
	})
	time.Sleep(5 * time.Second)
}
func TestParse(t *testing.T) {
	t.Log(parseMailTemplate(&Message{
		IP:      "127.0.0.1",
		Subject: "test subject",
		Body:    "test body select",
		To:      []string{"mananhai@highlight.mobi"},
	}))
}
