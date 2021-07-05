package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"wepush/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type GetTokenResult struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

//...

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

func Start(config *config.Config) {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(ConfigMiddleware(config))
	r.Use(TokenMiddleware(config))
	// 与群控终端相关的接口
	// 1. 策略相关
	r.GET("/push", PushMsg)
	r.POST("/push", PushMsg)
	r.Run(":80") // listen and serve on 0.0.0.0:80
}

type PostData struct {
	Msg string `json:"msg"`
}

func PushMsg(c *gin.Context) {
	log.Debug("请求/push接口")
	// 获取配置
	config, ok := c.MustGet("config").(*config.Config)
	if !ok {
		log.Warn("获取配置出错")
		c.JSON(500, gin.H{
			"error": "服务端出错",
		})
		return
	}
	// 获取msg
	var msg string
	if c.Request.Method == "GET" {
		msg = c.Query("msg")
	} else {
		msg = c.PostForm("msg")
	}
	if msg == "" {
		json := PostData{}
		c.BindJSON(&json)
		msg = json.Msg
	}
	if msg == "" {
		c.JSON(400, gin.H{
			"error": "请求参数出错",
		})
		return
	}
	log.Debugf("获取到的msg为%s", msg)
	// 获取token
	token, err := GetToken(config)
	if err != nil {
		log.Warn(err)
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("服务端出错: %s", err),
		})
		return
	}
	log.Debugf("获取到的token为:%s", token)

	// 发送消息
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	m := TextMessage{
		Agentid: config.WechatWorkConfig.AgentId,
		Msgtype: "text",
		Touser:  "@all",
	}
	m.Text.Content = msg
	postJson(url, m)

}

type TextMessage struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Agentid string `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

type TextCardMessage struct {
	Touser   string `json:"touser"`
	Msgtype  string `json:"msgtype"`
	Agentid  string `json:"agentid"`
	TextCard struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
	} `json:"textcard"`
}

type PostMsgResult struct {
	// {
	// 	"errcode" : 0,
	// 	"errmsg" : "ok",
	// 	"invaliduser" : "userid1|userid2",
	// 	"invalidparty" : "partyid1|partyid2",
	// 	"invalidtag": "tagid1|tagid2"
	//   }
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
}

func postJson(url string, m TextMessage) (body []byte, err error) {
	jsonStr, _ := json.Marshal(m)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	var res PostMsgResult
	err = json.Unmarshal(body, &res)
	log.Debugf("post msg返回值为%v", res)
	if err != nil {
		log.Error("解析postmsg返回值失败")
	}

	if res.Errcode != 0 {
		log.Errorf("postJson errmsg:%s", res.Errmsg)
	}
	return body, err
}

// ApiMiddleware will add the db connection to the context
func ConfigMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", config)
		c.Next()
	}
}

func TokenMiddleware(config *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			token = c.Query("token")
		}
		if token != config.AuthToken {
			c.JSON(401, gin.H{
				"error": "认证失败",
			})
			c.Abort()
			return
		}
		//请求处理
		c.Next()
	}
}
