package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            int
	DBPath          string
	DownloadDir     string
	AdminPassword   string
	BinDir          string
	AllowOrigins    string
	AllowInsecure   bool   // 是否允许非HTTPS环境（开发模式）
	DownloadTimeout int    // 下载超时时间（秒），0表示不限制
}

func Load() *Config {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	downloadDir := os.Getenv("DOWNLOAD_DIR")
	if downloadDir == "" {
		downloadDir = "./downloads"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	binDir := os.Getenv("BIN_DIR")
	if binDir == "" {
		binDir = "./bin"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./db/data.db"
	}

	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:8080,http://127.0.0.1:8080"
	}

	allowInsecure, _ := strconv.ParseBool(os.Getenv("ALLOW_INSECURE"))

	downloadTimeout, _ := strconv.Atoi(os.Getenv("DOWNLOAD_TIMEOUT"))
	if downloadTimeout <= 0 {
		downloadTimeout = 0 // 0 表示不限制超时
	}

	return &Config{
		Port:            port,
		DBPath:          dbPath,
		DownloadDir:     downloadDir,
		AdminPassword:   adminPassword,
		BinDir:          binDir,
		AllowOrigins:    allowOrigins,
		AllowInsecure:   allowInsecure,
		DownloadTimeout: downloadTimeout,
	}
}

func (c *Config) GetN_m3u8DLREPath() string {
	return fmt.Sprintf("%s/N_m3u8DL-RE", c.BinDir)
}

func (c *Config) GetFFmpegPath() string {
	return fmt.Sprintf("%s/ffmpeg", c.BinDir)
}

func (c *Config) GetMp4decryptPath() string {
	return fmt.Sprintf("%s/mp4decrypt", c.BinDir)
}

// GetDownloadTimeout 返回下载超时时间
func (c *Config) GetDownloadTimeout() time.Duration {
	if c.DownloadTimeout <= 0 {
		return 0 // 不限制超时
	}
	return time.Duration(c.DownloadTimeout) * time.Second
}
