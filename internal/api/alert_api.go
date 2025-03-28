package api

import (
"net/http"
"strconv"

"github.com/crypto-tracker/internal/models"
"github.com/crypto-tracker/internal/services"
"github.com/gin-gonic/gin"
)

// RegisterAlertRoutes 注册警报相关API路由
func RegisterAlertRoutes(router *gin.RouterGroup, alertService *services.AlertService) {
router.GET("/alerts", getAlerts(alertService))
router.GET("/alerts/:id", getAlertByID(alertService))
router.POST("/alerts", createAlert(alertService))
router.PUT("/alerts/:id", updateAlert(alertService))
router.DELETE("/alerts/:id", deleteAlert(alertService))
}

// getAlerts 获取所有警报
func getAlerts(alertService *services.AlertService) gin.HandlerFunc {
return func(c *gin.Context) {
alerts, err := alertService.GetAlerts()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

c.JSON(http.StatusOK, gin.H{
"success": true,
"data":    alerts,
})
}
}

// getAlertByID 通过ID获取警报
func getAlertByID(alertService *services.AlertService) gin.HandlerFunc {
return func(c *gin.Context) {
idStr := c.Param("id")
id, err := strconv.ParseInt(idStr, 10, 64)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{
"success": false,
"error":   "无效的警报ID",
})
return
}

alert, err := alertService.GetAlertByID(id)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

c.JSON(http.StatusOK, gin.H{
"success": true,
"data":    alert,
})
}
}

// createAlert 创建新警报
func createAlert(alertService *services.AlertService) gin.HandlerFunc {
return func(c *gin.Context) {
var req models.AlertRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

// 验证警报类型
if !models.ValidateAlertType(string(req.Type)) {
c.JSON(http.StatusBadRequest, gin.H{
"success": false,
"error":   "无效的警报类型",
})
return
}

alert, err := alertService.CreateAlert(req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

c.JSON(http.StatusCreated, gin.H{
"success": true,
"data":    alert,
})
}
}

// updateAlert 更新警报
func updateAlert(alertService *services.AlertService) gin.HandlerFunc {
return func(c *gin.Context) {
idStr := c.Param("id")
id, err := strconv.ParseInt(idStr, 10, 64)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{
"success": false,
"error":   "无效的警报ID",
})
return
}

var req models.AlertRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

// 验证警报类型
if !models.ValidateAlertType(string(req.Type)) {
c.JSON(http.StatusBadRequest, gin.H{
"success": false,
"error":   "无效的警报类型",
})
return
}

alert, err := alertService.UpdateAlert(id, req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

c.JSON(http.StatusOK, gin.H{
"success": true,
"data":    alert,
})
}
}

// deleteAlert 删除警报
func deleteAlert(alertService *services.AlertService) gin.HandlerFunc {
return func(c *gin.Context) {
idStr := c.Param("id")
id, err := strconv.ParseInt(idStr, 10, 64)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{
"success": false,
"error":   "无效的警报ID",
})
return
}

if err := alertService.DeleteAlert(id); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{
"success": false,
"error":   err.Error(),
})
return
}

c.JSON(http.StatusOK, gin.H{
"success": true,
"message": "警报已成功删除",
})
}
}
