package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-monolith-frame/internal/config"
	"go-monolith-frame/internal/handler"
	"go-monolith-frame/internal/middleware"
	"go-monolith-frame/internal/model/entity"
	"go-monolith-frame/internal/repository"
	"go-monolith-frame/internal/service"
	"go-monolith-frame/pkg/cache"
	"go-monolith-frame/pkg/logger"
	"go-monolith-frame/pkg/mysql"
	"go-monolith-frame/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	env string
)

func init() {
	flag.StringVar(&env, "env", "dev", "运行环境: dev, prod")
	flag.Parse()
}

func main() {
	// 0. 加载 .env 文件（如果存在）
	if err := godotenv.Load(); err == nil {
		fmt.Println("Loaded .env file")
	}

	// 1. 加载配置
	cfg, err := config.Load(env)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting server...",
		zap.String("env", env),
		zap.String("addr", cfg.Server.Addr),
	)

	// 3. 初始化数据库
	if err := mysql.Init(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}
	defer mysql.Close()

	// 4. 自动迁移数据库表
	if err := mysql.AutoMigrate(&entity.User{}); err != nil {
		logger.Fatal("Failed to auto migrate database", zap.Error(err))
	}
	logger.Info("Database migration completed")

	// 5. 初始化缓存
	cacheClient, err := cache.New(&cfg.Cache)
	if err != nil {
		logger.Fatal("Failed to init cache", zap.Error(err))
	}
	defer cacheClient.Close()

	// 6. 初始化组件
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo, cacheClient)
	userHandler := handler.NewUserHandler(userService)

	// 7. 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 8. 创建路由
	r := gin.New()

	// 注册中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 公开路由（无需认证）
		users := api.Group("/users")
		{
			users.POST("/login", userHandler.Login)
			users.POST("", userHandler.CreateUser)
		}

		// 需要认证的路由
		authGroup := api.Group("")
		authGroup.Use(middleware.Auth())
		{
			users := authGroup.Group("/users")
			{
				users.GET("", userHandler.ListUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
				users.PATCH("/:id/password", userHandler.ChangePassword)
			}
		}
	}

	// 9. 创建HTTP服务器
	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 10. 优雅启动和关闭
	go func() {
		logger.Info("Server is running", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭，等待5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
