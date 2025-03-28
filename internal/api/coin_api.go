package api

import (
"net/http"

"github.com/crypto-tracker/internal/services"
"github.com/gin-gonic/gin"
)

// RegisterCoinRoutes 注册币种相关API路由
func RegisterCoinRoutes(router *gin.RouterGroup, priceService *services.PriceService) {
router.GET("/coins", listCoins(priceService))
router.GET("/exchanges", listExchanges(priceService))
}

// listCoins 获取支持的加密货币列表
func listCoins(priceService *services.PriceService) gin.HandlerFunc {
return func(c *gin.Context) {
coins := priceService.GetSupportedCoins()
c.JSON(http.StatusOK, gin.H{
"success": true,
"data":    coins,
})
}
}

// listExchanges 获取支持的交易所列表
func listExchanges(priceService *services.PriceService) gin.HandlerFunc {
return func(c *gin.Context) {
exchanges := priceService.GetSupportedExchanges()
c.JSON(http.StatusOK, gin.H{
"success": true,
"data":    exchanges,
})
}
}
