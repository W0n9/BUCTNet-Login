package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var errParse = errors.New("error-parse")

// genCallback 生成 callback 字符串
func genCallback() string {
	return fmt.Sprintf("jsonp%d", int(time.Now().Unix()))
}

// DoRequest 发送带 callback 的 GET 请求
func DoRequest(url string, params url.Values) (*http.Response, error) {

	// add callback
	params.Add("callback", genCallback())
	params.Add("_", fmt.Sprint(time.Now().UnixNano()))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = params.Encode()
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetJson 发送请求并解析 JSON 响应
func GetJson(url string, data url.Values, res interface{}) (err error) {
	resp, err := DoRequest(url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	rawStr := string(raw)

	// cut jsonp
	start := strings.Index(rawStr, "(")
	end := strings.LastIndex(rawStr, ")")
	if start == -1 && end == -1 {
		log.Debug("raw response:", rawStr)
		return errParse
	}
	dt := string(raw)[start+1 : end]

	return json.Unmarshal([]byte(dt), &res)
}
