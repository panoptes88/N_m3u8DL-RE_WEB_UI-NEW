package handler

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"N_m3u8DL-RE-WEB-UI/internal/config"
	"N_m3u8DL-RE-WEB-UI/internal/service"

	"github.com/gin-gonic/gin"
)

// 视频 MIME 类型映射（包级别变量，避免重复创建）
var videoMimeTypes = map[string]string{
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".mkv":  "video/x-matroska",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
	".flv":  "video/x-flv",
	".wmv":  "video/x-ms-wmv",
	".m4v":  "video/x-m4v",
	".3gp":  "video/3gpp",
	".mpg":  "video/mpeg",
	".mpeg": "video/mpeg",
}

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
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", getVideoMimeType(filename))
	c.Header("Content-Disposition", "attachment; filename="+url.PathEscape(filename))
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

func getVideoMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if mime, ok := videoMimeTypes[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}
