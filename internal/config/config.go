package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port         int
	DBPath       string
	DownloadDir  string
	AdminPassword string
	BinDir       string
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
		dbPath = "./data.db"
	}

	return &Config{
		Port:         port,
		DBPath:       dbPath,
		DownloadDir:  downloadDir,
		AdminPassword: adminPassword,
		BinDir:       binDir,
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
