package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/saucerman/wepush/config"

	"github.com/gin-gonic/gin"
)

func Start(config *config.Config) {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(ConfigMiddleware(config))
	r.Use(TokenMiddleware(config))
	r.POST("/wepush", PushMsg)
	r.Run(":80") // listen and serve on 0.0.0.0:80
}

type PostMsgResult struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
}

func postJson(url string, content []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res PostMsgResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if res.Errcode != 0 {
		err = fmt.Errorf("发送消息错误，err: %s", res.Errmsg)
		return nil, err
	}
	return body, nil
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
