// model 包定义 API 往返使用的请求/响应数据结构。
package model

import (
	"fmt"
	"strings"

	"github.com/W0n9/BUCTNet-Login/utils"
)

// ChallengeResp challenge 接口返回：加密用 token 与客户端 IP。
type ChallengeResp struct {
	Challenge string `json:"challenge"`
	ClientIp  string `json:"client_ip"`
}

// ActionResp 登录、注销等操作接口的通用响应。
type ActionResp struct {
	Res      string      `json:"res"`
	Error    string      `json:"error"`
	Ecode    interface{} `json:"ecode"`
	ErrorMsg string      `json:"error_msg"`
	ClientIp string      `json:"client_ip"`
}

// InfoResp 用户信息查询结果；字节类字段用 uint64 承接网关可能返回的超大无符号数。
type InfoResp struct {
	ServerFlag    int64   `json:"ServerFlag"`
	AddTime       int64   `json:"add_time"`
	AllBytes      uint64  `json:"all_bytes"`
	BytesIn       uint64  `json:"bytes_in"`
	BytesOut      uint64  `json:"bytes_out"`
	CheckoutDate  int64   `json:"checkout_date"`
	Domain        string  `json:"domain"`
	Error         string  `json:"error"`
	GroupID       string  `json:"group_id"`
	KeepaliveTime int64   `json:"keepalive_time"`
	OnlineIP      string  `json:"online_ip"`
	ProductsName  string  `json:"products_name"`
	RealName      string  `json:"real_name"`
	RemainBytes   uint64  `json:"remain_bytes"`
	RemainSeconds int64   `json:"remain_seconds"`
	SumBytes      uint64  `json:"sum_bytes"`
	SumSeconds    int64   `json:"sum_seconds"`
	UserBalance   float64 `json:"user_balance"`
	UserCharge    float64 `json:"user_charge"`
	UserMac       string  `json:"user_mac"`
	UserName      string  `json:"user_name"`
	WalletBalance float64 `json:"wallet_balance"`
}

// String 返回 InfoResp 的字符串表示
func (r *InfoResp) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(" 在线IP: %s\n", r.OnlineIP))
	sb.WriteString(fmt.Sprintf("上网账号: %s\n", r.UserName))
	sb.WriteString(fmt.Sprintf("电子钱包: ￥%.2f\n", r.WalletBalance))
	sb.WriteString(fmt.Sprintf("套餐余额: ￥%.2f\n", r.UserBalance))
	sb.WriteString(fmt.Sprintf("已用流量: %s\n", utils.FormatFlux(r.SumBytes)))
	sb.WriteString(fmt.Sprintf("在线时长: %s\n", utils.FormatTime(r.SumSeconds)))

	return sb.String()
}

// InfoResult OAuth/门户回调中的 ac_id、用户名等（若有使用）。
type InfoResult struct {
	Acid        int    `json:"ac_id"`
	Username    string `json:"username"`
	ClientIp    string `json:"client_ip"`
	AccessToken string `json:"access_token"`
}
