package model

import (
	"encoding/json"
	"fmt"
)

type Account struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	AccessToken string `json:"access_token"`
	Acid        int    `json:"acid"`
}

// JSONString 返回 Account 的 JSON 字符串
func (a *Account) JSONString() (jsonStr string, err error) {
	jsonData, err := json.Marshal(a)
	if err != nil {
		return
	}
	jsonStr = string(jsonData)
	return
}

// JSONBytes 返回 Account 的 JSON 字节数组
func (a *Account) JSONBytes() (jsonData []byte, err error) {
	return json.Marshal(a)
}

// String 返回 Account 的字符串表示
func (a *Account) String() string {
	return fmt.Sprintln("用户名:", a.Username)
}
