package handler

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"N_m3u8DL-RE-WEB-UI/internal/config"
	"N_m3u8DL-RE-WEB-UI/internal/service"

	"github.com/gin-gonic/gin"
)

func ListFiles(c *gin.Context) {
	cfg := config.Load()
	files, err := service.ListFiles(cfg.DownloadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func DownloadFile(c *gin.Context) {
	filename := c.Query("name")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	// 防止路径遍历
	filename = filepath.Base(filename)

	cfg := config.Load()
	filePath := filepath.Join(cfg.DownloadDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+url.PathEscape(filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

func DeleteFile(c *gin.Context) {
	filename := c.Param("name")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	// 防止路径遍历
	filename = filepath.Base(filename)

	cfg := config.Load()
	filePath := filepath.Join(cfg.DownloadDir, filename)

	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
