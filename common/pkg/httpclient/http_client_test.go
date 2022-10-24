package httpclient

import (
	"strings"
	"testing"
)

// http-client

func TestPostJson(t *testing.T) {
	cmd := "http://www.tmnhs.top/ping?{\"name\":\"test\"}"
	url := strings.Split(cmd, "?")
	result, err := PostJson(url[0], `{"name":"test"}`, 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
