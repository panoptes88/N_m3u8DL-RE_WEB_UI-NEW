package handler

import (
	"net/http"
	"strconv"

	"N_m3u8DL-RE-WEB-UI/internal/model"
	"N_m3u8DL-RE-WEB-UI/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	URL              string `json:"url" binding:"required"`
	OutputName       string `json:"output_name" binding:"required"`
	ThreadCount      int    `json:"thread_count"`
	RetryCount       int    `json:"retry_count"`
	Headers          string `json:"headers"`
	BaseURL          string `json:"base_url"`
	DelAfterDone     bool   `json:"del_after_done"`
	BinaryMerge      bool   `json:"binary_merge"`
	AutoSelect       bool   `json:"auto_select"`
	Key              string `json:"key"`
	DecryptionEngine string `json:"decryption_engine"`
	CustomArgs       string `json:"custom_args"`
	CustomProxy      string `json:"custom_proxy"`
}

func ListTasks(c *gin.Context) {
	var tasks []model.Task
	status := c.Query("status")

	query := model.GetDB().Model(&model.Task{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query = query.Order("created_at DESC").Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 设置默认值
	if req.ThreadCount <= 0 {
		req.ThreadCount = 32
	}
	if req.RetryCount <= 0 {
		req.RetryCount = 15
	}
	if req.DecryptionEngine == "" {
		req.DecryptionEngine = "MP4DECRYPT"
	}

	// 转换为 service 层的请求类型
	serviceReq := &service.CreateTaskRequest{
		URL:              req.URL,
		OutputName:       req.OutputName,
		ThreadCount:      req.ThreadCount,
		RetryCount:       req.RetryCount,
		Headers:          req.Headers,
		BaseURL:          req.BaseURL,
		DelAfterDone:     req.DelAfterDone,
		BinaryMerge:      req.BinaryMerge,
		AutoSelect:       req.AutoSelect,
		Key:              req.Key,
		DecryptionEngine: req.DecryptionEngine,
		CustomArgs:       req.CustomArgs,
		CustomProxy:      req.CustomProxy,
	}

	task, err := service.CreateTask(serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetTask(c *gin.Context) {
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

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := service.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetTaskLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	logContent, err := service.GetTaskLog(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": logContent})
}
