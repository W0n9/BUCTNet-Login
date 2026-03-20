// utils 包提供 HTTP JSONP 请求、响应体排错日志，以及流量/时长等格式化工具。
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/W0n9/BUCTNet-Login/logger"
)

var errParse = errors.New("error-parse")
var log = logger.GetLogger()

// maxLogBody 失败日志中单次打印响应体的最大长度，防止撑满磁盘
const maxLogBody = 2048

// truncateForLog 裁剪过长文本，便于排错时人工阅读
func truncateForLog(s string) string {
	if len(s) <= maxLogBody {
		return s
	}
	return s[:maxLogBody] + "…(truncated)"
}

// logGetJSONFailure：Debug 打印截断 body，Error 打印原因，便于对照网关返回排错
func logGetJSONFailure(reqURL, rawBody string, cause error) {
	log.Debugw("getjson raw response", "url", reqURL, "body", truncateForLog(rawBody))
	log.Errorw("getjson failed", "url", reqURL, "err", cause)
}

// genCallback 生成 JSONP 所需的 callback 参数名
func genCallback() string {
	return fmt.Sprintf("jsonp%d", int(time.Now().Unix()))
}

// DoRequest 发送 GET，自动附加 callback 与时间戳参数（校园网接口多为 JSONP）
func DoRequest(url string, params url.Values) (*http.Response, error) {

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

// GetJson 请求 JSONP，剥去外层 `callback(...)` 后将内部 JSON 解到 res（res 为指向目标结构体的指针）。
func GetJson(url string, data url.Values, res interface{}) (err error) {
	resp, err := DoRequest(url, data)
	if err != nil {
		log.Errorw("getjson request failed", "url", url, "err", err)
		return err
	}
	defer resp.Body.Close()
	reqURL := resp.Request.URL.String()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorw("getjson read body failed", "url", reqURL, "err", err)
		return err
	}
	rawStr := string(raw)

	// 从 `xxx({...})` 中截取花括号 JSON 片段
	start := strings.Index(rawStr, "(")
	end := strings.LastIndex(rawStr, ")")
	if start == -1 || end == -1 || end <= start {
		logGetJSONFailure(reqURL, rawStr, errParse)
		return errParse
	}
	dt := rawStr[start+1 : end]

	if err := json.Unmarshal([]byte(dt), res); err != nil {
		logGetJSONFailure(reqURL, rawStr, err)
		return err
	}
	return nil
}
