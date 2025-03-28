package api

import (
"net/http"
"strconv"

"github.com/crypto-tracker/internal/services"
"github.com/gin-gonic/gin"
)

// RegisterPriceRoutes 注册价格相关API路由
func RegisterPriceRoutes(router *gin.RouterGroup, priceService *services.PriceService) {
router.GET("/prices", getAllPrices(priceService))
router.GET("/prices/:coin", getCoinPrice(priceService))
router.GET("/history/:coin", getPriceHistory(priceService))
}

// getAllPrices 获取所有币种的价格
func getAllPrices(priceService *services.PriceService) gin.HandlerFunc {
return func(c *gin.Context) {
prices := priceService.GetAllPrices()
c.JSON(http.StatusOK, gin.H{
"success": true,
"data":    prices,
})
}
}

// getCoinPrice 获取指定币种的价格
func getCoinPrice(priceService *services.PriceService) gin.HandlerFunc {
return func(c *gin.Context) {
coin := c.Param("coin")
exchange := c.DefaultQuery("exchange", "模拟交易所")

price, err := priceService.GetPrice(exchange, coin)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

c.JSON(http.StatusOK, gin.H{
"success": true,
"data":    price,
})
}
}

// getPriceHistory 获取指定币种的价格历史
func getPriceHistory(priceService *services.PriceService) gin.HandlerFunc {
return func(c *gin.Context) {
coin := c.Param("coin")
exchange := c.DefaultQuery("exchange", "模拟交易所")
limitStr := c.DefaultQuery("limit", "100")
limit, err := strconv.Atoi(limitStr)
if err != nil {
limit = 100 // 默认限制100条记录
}

history, err := priceService.GetHistory(exchange, coin, limit)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

c.JSON(http.StatusOK, gin.H{
"success": true,
"data":    history,
})
}
}
