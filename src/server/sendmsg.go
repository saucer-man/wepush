package server

import (
	"encoding/json"
	"fmt"

	"github.com/saucerman/wepush/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type MarkdownMessage struct {
	Touser   string `json:"touser"`
	Msgtype  string `json:"msgtype"`
	Agentid  string `json:"agentid"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

type TextMessage struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Agentid string `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

type PostData struct {
	Msg    string `json:"msg"`
	Touser string `json:"touser"`
	Type   string `json:"type"`
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
	// 解析postjson
	postdata := PostData{}
	c.BindJSON(&postdata)
	log.Debugf("获取到的post json为 %+v", postdata)
	if postdata.Msg == "" {
		c.JSON(400, gin.H{
			"error": "请求参数出错",
		})
		return
	}

	if postdata.Touser == "" {
		postdata.Touser = "@all"
	}
	if postdata.Type == "" {
		postdata.Type = "text"
	}
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
	var contentByte []byte
	if postdata.Type == "markdown" {
		m := MarkdownMessage{
			Agentid: config.WechatWorkConfig.AgentId,
			Msgtype: postdata.Type,
			Touser:  postdata.Touser,
		}
		m.Markdown.Content = postdata.Msg
		contentByte, err = json.Marshal(m)
		if err != nil {
			log.Errorf("Marshal 错误:%v", err)
			c.JSON(500, gin.H{
				"error": fmt.Sprintf("服务端错误: %v", err),
			})
			return
		}

	} else if postdata.Type == "text" {
		m := TextMessage{
			Agentid: config.WechatWorkConfig.AgentId,
			Msgtype: postdata.Type,
			Touser:  postdata.Touser,
		}
		m.Text.Content = postdata.Msg
		contentByte, err = json.Marshal(m)
		if err != nil {
			log.Errorf("Marshal 错误:%v", err)
			c.JSON(500, gin.H{
				"error": fmt.Sprintf("服务端错误: %v", err),
			})
			return
		}

	} else {
		log.Errorf("错误的消息类型: %s", postdata.Type)
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("错误的消息类型: %s", postdata.Type),
		})
		return
	}

	_, err = postJson(url, contentByte)
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("发送消息错误: %v", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": nil,
	})

}
