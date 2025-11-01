package config

import (
	"os"
)

// Config 应用配置
type Config struct {
	TencentCloud TencentCloudConfig
}

// TencentCloudConfig 腾讯云配置
type TencentCloudConfig struct {
	SecretID  string
	SecretKey string
	Region    string
	BotBizID  string // LKE智能体业务ID
	AppID     string // LKE应用ID
	BotAppKey string // LKE Bot应用密钥
	Endpoint  string // LKE API端点
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		TencentCloud: TencentCloudConfig{
			SecretID:  getEnv("TENCENT_SECRET_ID", ""),
			SecretKey: getEnv("TENCENT_SECRET_KEY", ""),
			Region:    getEnv("TENCENT_REGION", "ap-guangzhou"),
			BotBizID:  getEnv("LKE_BOT_BIZ_ID", ""),
			AppID:     getEnv("LKE_APP_ID", ""),
			BotAppKey: getEnv("LKE_APP_KEY", ""),
			Endpoint:  getEnv("LKE_ENDPOINT", "https://wss.lke.cloud.tencent.com"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
