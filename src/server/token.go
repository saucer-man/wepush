package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/saucerman/wepush/config"
)

type GetTokenResult struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetToken(config *config.Config) (string, error) {
	if !config.TokenConfig.ExpiredAt.IsZero() && time.Now().Before(config.TokenConfig.ExpiredAt) {
		return config.TokenConfig.Token, nil
	}

	resp, err := http.Get("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + config.WechatWorkConfig.CorpId + "&corpsecret=" + config.WechatWorkConfig.CorpSecret)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var res GetTokenResult
	err = json.Unmarshal(data, &res)
	if err != nil {
		return "", err
	}
	if res.Errcode != 0 {
		return "", fmt.Errorf("获取accress_token失败：%s", res.Errmsg)
	}
	tokenStr := res.AccessToken
	tokenStr = strings.Replace(tokenStr, "\"", "", -1)
	config.TokenConfig.Token = tokenStr
	config.TokenConfig.ExpiredAt = time.Now().Add(1 * time.Hour)
	return tokenStr, err
}
