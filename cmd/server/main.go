package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/crypto-tracker/internal/api"
	"github.com/crypto-tracker/internal/config"
	"github.com/crypto-tracker/internal/services"
	"github.com/crypto-tracker/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	db, err := database.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close()

	// 初始化服务
	priceService := services.NewPriceService(cfg.Exchanges)
	alertService := services.NewAlertService(db)

	// 启动价格更新任务
	go startPriceUpdateTask(priceService, cfg.UpdateInterval)

	// 设置Gin路由
	router := setupRouter(priceService, alertService)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("服务器启动于 http://%s", addr)
	log.Fatal(router.Run(addr))
}

// 设置Gin路由
func setupRouter(priceService *services.PriceService, alertService *services.AlertService) *gin.Engine {
	router := gin.Default()

	// 加载HTML模板
	router.LoadHTMLGlob("web/templates/*")

	// 静态资源
	router.Static("/static", "./web/static")

	// 首页路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "加密货币价格追踪器",
		})
	})

	// API路由
	apiV1 := router.Group("/api/v1")
	{
		api.RegisterCoinRoutes(apiV1, priceService)
		api.RegisterPriceRoutes(apiV1, priceService)
		api.RegisterAlertRoutes(apiV1, alertService)
	}

	return router
}

// 启动价格更新定时任务
func startPriceUpdateTask(priceService *services.PriceService, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// 立即执行一次更新
	if err := priceService.UpdateAllPrices(); err != nil {
		log.Printf("更新价格失败: %v", err)
	}

	// 定时执行更新
	for range ticker.C {
		if err := priceService.UpdateAllPrices(); err != nil {
			log.Printf("更新价格失败: %v", err)
		}
	}
}
