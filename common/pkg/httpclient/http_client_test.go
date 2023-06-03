package httpclient

import (
	"testing"
)

// http-client

func TestPostJson(t *testing.T) {
	//cmd := "http://www.tmnhs.top/ping?{\"name\":\"test\"}"
	cmd := "http://localhost:8089/ping"
	//url := strings.Split(cmd, "?")
	result, err := Get(cmd, 300)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
