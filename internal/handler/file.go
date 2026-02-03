package handler

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"N_m3u8DL-RE-WEB-UI/internal/config"
	"N_m3u8DL-RE-WEB-UI/internal/model"
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

// GetDownloadProgress 返回所有活跃任务的下载进度
func GetDownloadProgress(c *gin.Context) {
	tasks, err := service.GetActiveTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 清理已完成任务的 PID
	for _, task := range tasks {
		if task.Status == model.TaskStatusCompleted || task.Status == model.TaskStatusFailed {
			if task.PID > 0 {
				service.CleanupTaskPID(task.ID)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// GetTaskProgress 返回单个任务的下载进度
func GetTaskProgress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	task, err := service.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       task.ID,
		"status":   task.Status,
		"progress": task.Progress,
		"speed":    task.Speed,
	})
}
