package handler

import (
	"net/http"

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
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	if !model.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 设置 cookie，有效期 24 小时
	c.SetCookie("auth_token", user.Username, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func GetUser(c *gin.Context) {
	username, err := c.Cookie("auth_token")
	if err != nil || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
	})
}

func IsLoggedIn(c *gin.Context) bool {
	username, err := c.Cookie("auth_token")
	return err == nil && username != ""
}

func GetCurrentUserID(c *gin.Context) uint {
	return 0
}

func GetCurrentUsername(c *gin.Context) string {
	username, _ := c.Cookie("auth_token")
	return username
}

func RefreshSession(c *gin.Context) {
	// Cookie 自动刷新
}
