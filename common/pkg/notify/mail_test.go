package notify

import (
	"testing"
	"time"
)

//https://open.feishu.cn/open-apis/bot/v2/hook/184aa608-e62c-4d49-9621-d8fb22ed1254
func TestSend(t *testing.T) {
	//Init(465,"1685290935@qq.com","smtp.qq.com","wqkncnzhdvbpdabj","test")
	go Serve()
	Send(&Message{
		Subject: "test subject",
		Body:    "test body select",
		To:      []string{"mananhai@highlight.mobi"},
	})
	time.Sleep(3 * time.Second)
}

//报警信息
func TestFeiShuSend(t *testing.T) {
	//sendMsg("https://open.feishu.cn/open-apis/bot/v2/hook/4df45fb8-6458-4af9-a2d5-6c0396895743","test")
}

func TestParse(t *testing.T) {
	t.Log(parseMailTemplate(&Message{
		IP:        "127.0.0.1",
		Subject:   "test subject",
		Body:      "test body select",
		To:        []string{"mananhai@highlight.mobi"},
		OccurTime: time.Now(),
	}))
}
