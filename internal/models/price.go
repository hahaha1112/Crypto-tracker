package models

import (
"time"
)

// Price 加密货币价格模型
type Price struct {
ID        int64     `json:"id" db:"id"`
Coin      string    `json:"coin" db:"coin"`
Exchange  string    `json:"exchange" db:"exchange"`
Price     float64   `json:"price" db:"price"`
Change24h float64   `json:"change24h" db:"change24h"`
Volume24h float64   `json:"volume24h" db:"volume24h"`
UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// PriceHistory 价格历史记录
type PriceHistory struct {
ID        int64     `json:"id" db:"id"`
Coin      string    `json:"coin" db:"coin"`
Exchange  string    `json:"exchange" db:"exchange"`
Price     float64   `json:"price" db:"price"`
Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// CoinInfo 币种信息
type CoinInfo struct {
Symbol      string `json:"symbol"`
Name        string `json:"name"`
LogoURL     string `json:"logoUrl"`
Description string `json:"description"`
}

// SupportedCoins 返回支持的加密货币列表
func SupportedCoins() []CoinInfo {
return []CoinInfo{
{Symbol: "BTC", Name: "比特币", LogoURL: "/static/images/btc.png", Description: "世界上第一个也是最知名的加密货币"},
{Symbol: "ETH", Name: "以太坊", LogoURL: "/static/images/eth.png", Description: "支持智能合约的去中心化应用平台"},
{Symbol: "BNB", Name: "币安币", LogoURL: "/static/images/bnb.png", Description: "币安交易所的原生代币"},
{Symbol: "XRP", Name: "瑞波币", LogoURL: "/static/images/xrp.png", Description: "专注于跨境支付系统的加密货币"},
{Symbol: "ADA", Name: "艾达币", LogoURL: "/static/images/ada.png", Description: "一个基于科学哲学和研究驱动的区块链平台"},
{Symbol: "SOL", Name: "索拉纳", LogoURL: "/static/images/sol.png", Description: "高性能区块链，支持快速、低成本交易"},
{Symbol: "DOGE", Name: "狗狗币", LogoURL: "/static/images/doge.png", Description: "源于互联网模因的加密货币"},
{Symbol: "DOT", Name: "波卡", LogoURL: "/static/images/dot.png", Description: "实现跨链互操作性的区块链网络"},
}
}
