package core

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

var (
	// 创建全局 resty 客户端
	httpClient *resty.Client
)

func init() {
	// 初始化 resty 客户端
	httpClient = resty.New().
		SetTimeout(cfg.APITimeout).
		SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			// 默认跟随重定向但不执行任何操作，让响应中仍然包含重定向信息
			return http.ErrUseLastResponse
		}))
}

// get 发送 GET 请求
func get(addr string) (*http.Response, error) {
	resp, err := httpClient.R().Get(addr)
	return resp.RawResponse, err
}
