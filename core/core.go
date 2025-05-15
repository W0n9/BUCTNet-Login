package core

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vouv/srun/hash"
	"github.com/vouv/srun/model"
	"github.com/vouv/srun/resp"
	"github.com/vouv/srun/utils"
	"go.uber.org/zap"
)

const (
	baseAddr = "http://202.4.130.95"

	challengeUrl = "/cgi-bin/get_challenge"
	portalUrl    = "/cgi-bin/srun_portal"

	succeedUrl = "/cgi-bin/rad_user_info"
)

var log = zap.S()

// Prepare 获取 acid 等参数
func Prepare() (int, error) {
	first, err := get(baseAddr)
	if err != nil {
		return 1, err
	}
	second, err := get(first.Header.Get("Location"))
	if err != nil {
		return 1, err
	}
	target := second.Header.Get("location")
	query, _ := url.Parse(baseAddr + target)
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

	if err = utils.GetJson(baseAddr+portalUrl, formLogin, &ra); err != nil {
		log.Debugw("request error", "err", err)
		return
	}

	if ra.Res != "ok" {
		log.Debugw("response msg is not 'ok'", "msg", ra.ErrorMsg)
		if strings.Contains(ra.ErrorMsg, "Arrearage users") {
			err = errors.New("已欠费")
		} else {
			err = errors.New(fmt.Sprint(ra))
		}
		return
	}

	account.AccessToken = token
	account.Acid = acid
	return
}

// Info 查询用户信息 API
func Info() (info *model.InfoResp, err error) {

	// 余量查询
	err = utils.GetJson(baseAddr+succeedUrl, url.Values{}, &info)
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
	if err = utils.GetJson(baseAddr+portalUrl, q, &ra); err != nil {
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
	err = utils.GetJson(baseAddr+challengeUrl, qc, &res)
	return
}
