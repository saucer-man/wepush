package config

import (
	"os"
	"time"
)

type Config struct {
	AuthToken        string
	LogLevel         string
	LogPath          string
	WechatWorkConfig struct {
		DefaultReceiverUserId string
		CorpId                string
		CorpSecret            string
		AgentId               string
	}
	TokenConfig struct {
		Token     string
		ExpiredAt time.Time
	}
}

func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

func LoadConfig() Config {
	var config Config
	config.AuthToken = GetEnvDefault("AuthToken", "123456")
	config.LogLevel = GetEnvDefault("LogLevel", "debug")
	config.LogPath = GetEnvDefault("LogPath", "/var/log/wepush")
	config.WechatWorkConfig.CorpSecret = GetEnvDefault("CorpSecret", "xxx")
	config.WechatWorkConfig.CorpId = GetEnvDefault("CorpId", "xxx")
	config.WechatWorkConfig.DefaultReceiverUserId = GetEnvDefault("DefaultReceiverUserId", "@all")
	config.WechatWorkConfig.AgentId = GetEnvDefault("AgentId", "xxx")
	return config
}
