package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Mode    string `yaml:"mode"` // dev, prod
	} `yaml:"server"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Storage struct {
		Type     string `yaml:"type"`     // local, oss, s3
		BasePath string `yaml:"basePath"` // 本地存储的基础路径
		BaseURL  string `yaml:"baseUrl"`  // 文件访问的 URL 前缀
	} `yaml:"storage"`
	JWT struct {
		Secret               string `yaml:"secret"`
		AccessTokenDuration  int    `yaml:"accessTokenDuration"`  // 秒
		RefreshTokenDuration int    `yaml:"refreshTokenDuration"` // 秒
	} `yaml:"jwt"`
	Admin struct {
		Username string `yaml:"username"` // 管理员用户名
		Password string `yaml:"password"` // 管理员密码
	} `yaml:"admin"`
}

func Load() *Config {
	// 1. 先加载默认配置
	config := loadDefaults()

	// 2. 尝试从配置文件加载
	configFile := getConfigFile()
	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err == nil {
			if err := yaml.Unmarshal(data, config); err != nil {
				fmt.Printf("Warning: failed to parse config file %s: %v, using defaults\n", configFile, err)
			}
		}
	}

	// 3. 环境变量覆盖（优先级最高）
	applyEnvOverrides(config)

	return config
}

func getConfigFile() string {
	// 1. 优先使用环境变量指定的配置文件
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			return configFile
		}
	}

	// 2. 使用默认 config.yaml
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}

	// 3. 不存在配置文件，使用代码默认值
	return ""
}

func loadDefaults() *Config {
	config := &Config{}
	config.Server.Address = ":9050"
	config.Server.Mode = "dev"
	config.Database.DSN = "nebula.db"
	config.Storage.Type = "local"
	config.Storage.BasePath = "./uploads"
	config.Storage.BaseURL = "http://localhost:9050/files"
	config.JWT.Secret = "dev-secret-key-change-in-production"
	config.JWT.AccessTokenDuration = 7200    // 2 hours
	config.JWT.RefreshTokenDuration = 604800 // 7 days
	config.Admin.Username = "admin"
	config.Admin.Password = "admin123" // 默认密码，生产环境必须修改
	return config
}

func applyEnvOverrides(config *Config) {
	if val := os.Getenv("SERVER_ADDRESS"); val != "" {
		config.Server.Address = val
	}
	if val := os.Getenv("SERVER_MODE"); val != "" {
		config.Server.Mode = val
	}
	if val := os.Getenv("DATABASE_DSN"); val != "" {
		config.Database.DSN = val
	}
	if val := os.Getenv("STORAGE_TYPE"); val != "" {
		config.Storage.Type = val
	}
	if val := os.Getenv("STORAGE_BASE_PATH"); val != "" {
		config.Storage.BasePath = val
	}
	if val := os.Getenv("STORAGE_BASE_URL"); val != "" {
		config.Storage.BaseURL = val
	}
	if val := os.Getenv("JWT_SECRET"); val != "" {
		config.JWT.Secret = val
	}
	if val := os.Getenv("ADMIN_USERNAME"); val != "" {
		config.Admin.Username = val
	}
	if val := os.Getenv("ADMIN_PASSWORD"); val != "" {
		config.Admin.Password = val
	}
}

// GetAccessTokenDuration 返回 Access Token 有效期
func (c *Config) GetAccessTokenDuration() time.Duration {
	return time.Duration(c.JWT.AccessTokenDuration) * time.Second
}

// GetRefreshTokenDuration 返回 Refresh Token 有效期
func (c *Config) GetRefreshTokenDuration() time.Duration {
	return time.Duration(c.JWT.RefreshTokenDuration) * time.Second
}
