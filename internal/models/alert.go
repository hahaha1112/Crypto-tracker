package models

import (
"time"
)

// AlertType 警报类型枚举
type AlertType string

const (
PriceAbove AlertType = "price_above" // 价格高于阈值
PriceBelow AlertType = "price_below" // 价格低于阈值
ChangeAbove AlertType = "change_above" // 涨幅高于阈值
ChangeBelow AlertType = "change_below" // 跌幅低于阈值
)

// AlertStatus 警报状态枚举
type AlertStatus string

const (
Active   AlertStatus = "active"   // 活跃状态，将会触发
Triggered AlertStatus = "triggered" // 已触发状态
Disabled  AlertStatus = "disabled"  // 已禁用状态
)

// Alert 价格警报模型
type Alert struct {
ID          int64       `json:"id" db:"id"`
Coin        string      `json:"coin" db:"coin"`
Exchange    string      `json:"exchange" db:"exchange"`
Type        AlertType   `json:"type" db:"type"`
Threshold   float64     `json:"threshold" db:"threshold"`
Status      AlertStatus `json:"status" db:"status"`
Message     string      `json:"message" db:"message"`
CreatedAt   time.Time   `json:"createdAt" db:"created_at"`
TriggeredAt *time.Time  `json:"triggeredAt,omitempty" db:"triggered_at"`
}

// AlertRequest 创建/更新警报的请求
type AlertRequest struct {
Coin      string    `json:"coin" binding:"required"`
Exchange  string    `json:"exchange" binding:"required"`
Type      AlertType `json:"type" binding:"required"`
Threshold float64   `json:"threshold" binding:"required"`
Message   string    `json:"message"`
}

// ValidateAlertType 验证警报类型
func ValidateAlertType(alertType string) bool {
validTypes := []AlertType{PriceAbove, PriceBelow, ChangeAbove, ChangeBelow}
for _, t := range validTypes {
if AlertType(alertType) == t {
return true
}
}
return false
}
