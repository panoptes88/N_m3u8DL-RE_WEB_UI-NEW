package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"N_m3u8DL-RE-WEB-UI/internal/config"
	"N_m3u8DL-RE-WEB-UI/internal/handler"
	"N_m3u8DL-RE-WEB-UI/internal/middleware"
	"N_m3u8DL-RE-WEB-UI/internal/model"
	"N_m3u8DL-RE-WEB-UI/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	model.InitDB(cfg.DBPath)

	// 初始化下载目录
	if err := os.MkdirAll(cfg.DownloadDir, 0755); err != nil {
		log.Fatalf("创建下载目录失败: %v", err)
	}

	// 初始化默认管理员用户
	service.InitAdminUser(cfg.AdminPassword)

	// 启动下载任务轮询
	go service.StartTaskPolling(cfg)

	// Gin 设置
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	// 静态文件服务（前端构建产物）
	r.Static("/static", "./web/dist/static")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	// SPA 路由支持：所有未匹配的路径返回 index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	// API 路由
	api := r.Group("/api")
	{
		// 公开路由
		api.POST("/auth/login", handler.Login)
		api.POST("/auth/logout", handler.Logout)

		// 需要认证的路由
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// 用户
			protected.GET("/user", handler.GetUser)

			// 任务管理
			protected.GET("/tasks", handler.ListTasks)
			protected.POST("/tasks", handler.CreateTask)
			protected.GET("/tasks/:id", handler.GetTask)
			protected.DELETE("/tasks/:id", handler.DeleteTask)
			protected.GET("/tasks/:id/log", handler.GetTaskLog)

			// 文件管理
			protected.GET("/files", handler.ListFiles)
			protected.GET("/files/download", handler.DownloadFile)
			protected.DELETE("/files/:name", handler.DeleteFile)
		}
	}

	// 启动服务
	log.Printf("服务器启动在 :%d", cfg.Port)
	if err := r.Run(":" + fmt.Sprintf("%d", cfg.Port)); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
