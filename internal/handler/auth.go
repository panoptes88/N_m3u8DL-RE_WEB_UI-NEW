package handler

import (
	"net/http"

	"N_m3u8DL-RE-WEB-UI/internal/config"
	"N_m3u8DL-RE-WEB-UI/internal/model"
	"N_m3u8DL-RE-WEB-UI/internal/service"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	user, err := service.GetUserByUsername(req.Username)
	if err != nil || !model.CheckPassword(req.Password, user.Password) {
		// 统一错误信息，防止用户名枚举
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 创建会话，生成随机 token
	session, err := model.CreateSession(model.GetDB(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
		return
	}

	// 设置 cookie，有效期 24 小时
	cfg := config.Load()
	c.SetCookie("auth_token", session.Token, 86400, "/", "", !cfg.AllowInsecure, true)

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func Logout(c *gin.Context) {
	token, err := c.Cookie("auth_token")
	if err == nil && token != "" {
		model.DeleteSession(model.GetDB(), token)
	}
	c.SetCookie("auth_token", "", -1, "/", "", !config.Load().AllowInsecure, true)
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func GetUser(c *gin.Context) {
	token, err := c.Cookie("auth_token")
	if err != nil || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	user, err := model.GetUserByToken(model.GetDB(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "登录已过期"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
	})
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}

func ChangePassword(c *gin.Context) {
	username, err := c.Cookie("auth_token")
	if err != nil || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := service.ChangePassword(username, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
