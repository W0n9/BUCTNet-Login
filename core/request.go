package core

import (
	"bufio"
	"net"
	"net/http"
	"time"
)

// get 发送 GET 请求
func get(addr string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, addr, nil)
	return request(req)
}

// request 发送 HTTP 请求
func request(req *http.Request) (*http.Response, error) {
	conn, err := net.DialTimeout("tcp", req.URL.Hostname()+":http", time.Second)
	if err != nil {
		return nil, err
	}
	_ = req.Write(conn)
	return http.ReadResponse(bufio.NewReader(conn), req)
}
