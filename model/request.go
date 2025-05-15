package model

import (
	"fmt"
	"net/url"
)

// Challenge 生成 challenge 请求参数
func Challenge(username string) url.Values {
	return url.Values{
		"username": {username},
		"ip":       {""},
	}
}

// Login 生成 login 请求参数
func Login(username, password string, acid int) url.Values {
	return url.Values{
		"action":   {"login"},
		"username": {username},
		"password": {password},
		"ac_id":    {fmt.Sprint(acid)},
		"ip":       {""},
		"info":     {},
		"chksum":   {},
		"n":        {"200"},
		"type":     {"1"},
	}
}

// Logout 生成 logout 请求参数
func Logout(username string) url.Values {
	return url.Values{
		"action":   {"logout"},
		"username": {username},
	}
}
