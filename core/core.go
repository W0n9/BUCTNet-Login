package core

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/W0n9/BUCTNet-Login/config"
	"github.com/W0n9/BUCTNet-Login/hash"
	"github.com/W0n9/BUCTNet-Login/model"
	"github.com/W0n9/BUCTNet-Login/resp"
	"github.com/W0n9/BUCTNet-Login/utils"
	"go.uber.org/zap"
)

var (
	log = zap.S()
	cfg = config.GetConfig()
)

// Prepare 获取 acid 等参数
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

// Login 登录 API，处理登录流程
// step 1: prepare & get acid
// step 2: get challenge
// step 3: do login
func Login(account *model.Account) (err error) {
	log.Debugw("Login", "username", account.Username)

	// 先获取acid
	acid, err := Prepare()
	if err != nil {
		log.Debugw("prepare error", "err", err)
		return
	}

	username := account.Username

	// 创建登录表单
	formLogin := model.Login(username, account.Password, acid)

	// get token
	rc, err := getChallenge(username)
	if err != nil {
		log.Debugw("get challenge error", "err", err)
		return
	}

	token := rc.Challenge
	ip := rc.ClientIp

	formLogin.Set("ip", ip)
	formLogin.Set("info", hash.GenInfo(formLogin, token))
	formLogin.Set("password", hash.PwdHmd5("", token))
	formLogin.Set("chksum", hash.Checksum(formLogin, token))

	// response
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

// Info 查询用户信息 API
func Info() (info *model.InfoResp, err error) {
	// 余量查询
	err = utils.GetJson(cfg.BaseURL+cfg.SucceedURL, url.Values{}, &info)
	if err != nil {
		return nil, err
	}
	return
}

// Logout 注销 API
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

// getChallenge 获取 challenge
func getChallenge(username string) (res resp.ChallengeResp, err error) {
	qc := model.Challenge(username)
	err = utils.GetJson(cfg.BaseURL+cfg.ChallengeURL, qc, &res)
	return
}
