// core 封装校园网门户 HTTP 流程：准备 ac_id、challenge、登录/注销/查询。
package core

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/W0n9/BUCTNet-Login/config"
	"github.com/W0n9/BUCTNet-Login/hash"
	"github.com/W0n9/BUCTNet-Login/logger"
	"github.com/W0n9/BUCTNet-Login/model"
	"github.com/W0n9/BUCTNet-Login/resp"
	"github.com/W0n9/BUCTNet-Login/utils"
)

var (
	log = logger.GetLogger()
	cfg = config.GetConfig()
)

// Prepare 跟随门户重定向，从最终 URL 解析出 ac_id。
func Prepare() (int, error) {
	first, err := get(cfg.BaseURL)
	if err != nil {
		return 1, err
	}
	second, err := get(first.Header.Get("Location"))
	if err != nil {
		return 1, err
	}
	target := second.Header.Get("location")
	query, _ := url.Parse(cfg.BaseURL + target)
	return strconv.Atoi(query.Query().Get("ac_id"))
}

// Login 门户登录：解析 ac_id → 取 challenge/token → 带加密 info/chksum 提交表单。
func Login(account *model.Account) (err error) {
	log.Debugw("Login", "username", account.Username)

	// 跟随门户跳转，得到当前线路的 ac_id
	acid, err := Prepare()
	if err != nil {
		log.Debugw("prepare error", "err", err)
		return
	}

	username := account.Username

	// 创建登录表单
	formLogin := model.Login(username, account.Password, acid)

	rc, err := getChallenge(username)
	if err != nil {
		log.Debugw("get challenge error", "err", err)
		return
	}

	// challenge 同时作为后续 HMAC / XEncode 的密钥
	token := rc.Challenge
	ip := rc.ClientIp
	if strings.TrimSpace(token) == "" {
		return fmt.Errorf("获取 challenge 失败: token 为空")
	}

	formLogin.Set("ip", ip)
	formLogin.Set("info", hash.GenInfo(formLogin, token))
	formLogin.Set("password", hash.PwdHmd5("", token))
	formLogin.Set("chksum", hash.Checksum(formLogin, token))

	ra := resp.ActionResp{}

	if err = utils.GetJson(cfg.BaseURL+cfg.PortalURL, formLogin, &ra); err != nil {
		log.Debugw("request error", "err", err)
		return
	}

	if ra.Res != "ok" {
		log.Debugw("login failed", "response", ra)
		// 检查已知错误类型
		switch {
		case strings.Contains(ra.ErrorMsg, "Arrearage users"):
			return errors.New("已欠费")
		case strings.Contains(ra.ErrorMsg, "Password is error"):
			return errors.New("密码错误")
		case strings.Contains(ra.ErrorMsg, "Username is error"):
			return errors.New("用户名错误")
		default:
			// 未知错误返回原始信息
			return fmt.Errorf("登录失败: %s", ra.ErrorMsg)
		}
	}

	account.AccessToken = token
	account.Acid = acid
	return
}

// Info 查询在线用户信息（成功页 JSONP），解析为 InfoResp。
func Info() (*model.InfoResp, error) {
	var v model.InfoResp
	if err := utils.GetJson(cfg.BaseURL+cfg.SucceedURL, url.Values{}, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// Logout 调用注销接口并清空本地 Account 中的 token/acid。
func Logout(account *model.Account) (err error) {
	defer func() {
		account.AccessToken = ""
		account.Acid = 0
	}()

	q := model.Logout(account.Username)
	ra := resp.ActionResp{}
	if err = utils.GetJson(cfg.BaseURL+cfg.PortalURL, q, &ra); err != nil {
		log.Debugw("logout error", "err", err)
		err = ErrRequest
		return
	}
	if ra.Error != "ok" {
		log.Debugw("logout response", "resp", ra)
		err = ErrRequest
	}
	return
}

// getChallenge 请求网关下发 challenge 与客户端 IP。
func getChallenge(username string) (res resp.ChallengeResp, err error) {
	qc := model.Challenge(username)
	err = utils.GetJson(cfg.BaseURL+cfg.ChallengeURL, qc, &res)
	return
}
