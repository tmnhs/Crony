package httpclient

// http-client

import (
	"bytes"
	"fmt"
	"github.com/tmnhs/crony/common/pkg/logger"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(url string, timeout int64) (result string, err error) {
	var client = &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status code is not 200")
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("http get api url:%s send  err: %s", url, err.Error()))
		return
	}
	result = string(data)
	return
}

func PostParams(url string, params string, timeout int64) (result string, err error) {
	var client = &http.Client{}
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status code is not 200")
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("http post api url:%s send  err: %s", url, err.Error()))
		return
	}
	result = string(data)
	return
}

func PostJson(url string, body string, timeout int64) (result string, err error) {
	var client = &http.Client{}
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-type", "application/json")
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("http post api url:%s send  err: %s", url, err.Error()))
		return
	}
	result = string(data)
	return
}
