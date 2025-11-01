package main

import (
	"log"
	"os"

	"chinese-chess-ai/internal/api"
	"chinese-chess-ai/internal/config"
	"chinese-chess-ai/internal/game"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载.env文件: %v", err)
	}

	// 加载配置
	cfg := config.LoadConfig()
	log.Printf("配置加载成功: AppID=%s, BotAppKey=%s", maskKey(cfg.TencentCloud.AppID), maskKey(cfg.TencentCloud.BotAppKey))

	// 初始化游戏管理器（单例模式会自动初始化）
	game.GetManager(cfg) // 修复：移除未使用的变量赋值

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 静态文件服务
	r.Static("/static", "./web/static")
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/test_frontend.html", "./test_frontend.html")
	r.StaticFile("/test_canvas.html", "./test_canvas.html")
	r.StaticFile("/debug_frontend.html", "./debug_frontend.html")
	r.StaticFile("/frontend_debug.html", "./frontend_debug.html")

	// API路由
	apiGroup := r.Group("/api")
	{
		// 健康检查端点
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK", "message": "服务器运行正常"})
		})

		apiGroup.POST("/game/new", api.NewGame)
		apiGroup.POST("/game/move", api.PlayerMove)
		apiGroup.GET("/game/:id", api.GetGameState)
		apiGroup.POST("/game/:id/ai-move", api.AIMove)
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// maskKey 遮蔽密钥中间部分
func maskKey(key string) string {
	if key == "" || len(key) < 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
