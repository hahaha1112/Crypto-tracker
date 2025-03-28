package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 应用配置结构
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Database struct {
		Path string `json:"path"`
	} `json:"database"`
	Exchanges   []ExchangeConfig `json:"exchanges"`
	UpdateInterval int           `json:"updateInterval"` // 价格更新间隔（秒）
}

// ExchangeConfig 交易所配置
type ExchangeConfig struct {
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	APIKey string   `json:"apiKey"`
	Coins  []string `json:"coins"`
}

// LoadConfig 从配置文件加载配置
func LoadConfig() (*Config, error) {
	// 默认配置
	cfg := &Config{}
	cfg.Server.Host = "localhost"
	cfg.Server.Port = 8080
	cfg.Database.Path = "data/crypto.db"
	cfg.UpdateInterval = 60 // 默认60秒更新一次

	// 尝试加载配置文件
	configPath := filepath.Join("configs", "config.json")
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			return nil, err
		}
	}

	// 确保目录存在
	dbDir := filepath.Dir(cfg.Database.Path)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
