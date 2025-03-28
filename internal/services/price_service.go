package services

import (
"fmt"
"sync"
"time"

"github.com/crypto-tracker/internal/config"
"github.com/crypto-tracker/internal/models"
)

// PriceService 管理加密货币价格数据
type PriceService struct {
exchanges []config.ExchangeConfig
prices    map[string]map[string]*models.Price // map[交易所][币种]价格
history   map[string]map[string][]*models.PriceHistory // map[交易所][币种]价格历史
mutex     sync.RWMutex
}

// NewPriceService 创建价格服务
func NewPriceService(exchanges []config.ExchangeConfig) *PriceService {
return &PriceService{
exchanges: exchanges,
prices:    make(map[string]map[string]*models.Price),
history:   make(map[string]map[string][]*models.PriceHistory),
mutex:     sync.RWMutex{},
}
}

// UpdateAllPrices 更新所有交易所的所有币种价格
func (s *PriceService) UpdateAllPrices() error {
var wg sync.WaitGroup
errorsChan := make(chan error, len(s.exchanges))

for _, exchange := range s.exchanges {
wg.Add(1)
go func(exchange config.ExchangeConfig) {
defer wg.Done()
if err := s.updateExchangePrices(exchange); err != nil {
errorsChan <- fmt.Errorf("更新%s价格失败: %v", exchange.Name, err)
}
}(exchange)
}

wg.Wait()
close(errorsChan)

// 收集错误
var errs []error
for err := range errorsChan {
errs = append(errs, err)
}

if len(errs) > 0 {
return fmt.Errorf("部分价格更新失败: %v", errs)
}

return nil
}

// updateExchangePrices 从指定交易所获取价格数据
func (s *PriceService) updateExchangePrices(exchange config.ExchangeConfig) error {
// 这里是模拟数据，实际应用中应该调用各交易所的API
// 为每个交易所的每个币种生成模拟价格数据
for _, coin := range exchange.Coins {
price := s.generateMockPrice(coin)

s.mutex.Lock()
// 确保交易所映射存在
if _, ok := s.prices[exchange.Name]; !ok {
s.prices[exchange.Name] = make(map[string]*models.Price)
}
// 保存当前价格
s.prices[exchange.Name][coin] = price

// 同时保存到历史记录
if _, ok := s.history[exchange.Name]; !ok {
s.history[exchange.Name] = make(map[string][]*models.PriceHistory)
}
if _, ok := s.history[exchange.Name][coin]; !ok {
s.history[exchange.Name][coin] = make([]*models.PriceHistory, 0)
}
// 添加历史价格记录
history := &models.PriceHistory{
Coin:      coin,
Exchange:  exchange.Name,
Price:     price.Price,
Timestamp: time.Now(),
}
// 限制历史记录数量，保留最近1000条记录
histories := s.history[exchange.Name][coin]
if len(histories) >= 1000 {
histories = histories[1:]
}
histories = append(histories, history)
s.history[exchange.Name][coin] = histories
s.mutex.Unlock()
}

return nil
}

// generateMockPrice 生成模拟价格数据
func (s *PriceService) generateMockPrice(coin string) *models.Price {
// 模拟数据生成逻辑
var basePrice float64
var change24h float64
var volume24h float64

switch coin {
case "BTC":
basePrice = 50000.0 + float64(time.Now().Unix() % 5000)
change24h = -2.5 + (float64(time.Now().Unix()%500) / 100.0)
volume24h = 10000000000 + float64(time.Now().Unix()%1000000000)
case "ETH":
basePrice = 3000.0 + float64(time.Now().Unix() % 300)
change24h = 1.2 + (float64(time.Now().Unix()%300) / 100.0)
volume24h = 5000000000 + float64(time.Now().Unix()%500000000)
case "BNB":
basePrice = 500.0 + float64(time.Now().Unix() % 50)
change24h = 0.8 + (float64(time.Now().Unix()%200) / 100.0)
volume24h = 1000000000 + float64(time.Now().Unix()%100000000)
case "XRP":
basePrice = 1.0 + (float64(time.Now().Unix()%20) / 100.0)
change24h = -0.5 + (float64(time.Now().Unix()%100) / 100.0)
volume24h = 500000000 + float64(time.Now().Unix()%50000000)
case "ADA":
basePrice = 1.5 + (float64(time.Now().Unix()%30) / 100.0)
change24h = 3.0 + (float64(time.Now().Unix()%200) / 100.0)
volume24h = 300000000 + float64(time.Now().Unix()%30000000)
case "SOL":
basePrice = 100.0 + float64(time.Now().Unix() % 10)
change24h = 5.0 + (float64(time.Now().Unix()%400) / 100.0)
volume24h = 800000000 + float64(time.Now().Unix()%80000000)
case "DOGE":
basePrice = 0.1 + (float64(time.Now().Unix()%5) / 100.0)
change24h = -1.0 + (float64(time.Now().Unix()%200) / 100.0)
volume24h = 200000000 + float64(time.Now().Unix()%20000000)
case "DOT":
basePrice = 20.0 + float64(time.Now().Unix() % 3)
change24h = 0.3 + (float64(time.Now().Unix()%150) / 100.0)
volume24h = 400000000 + float64(time.Now().Unix()%40000000)
default:
basePrice = 10.0 + (float64(time.Now().Unix()%100) / 10.0)
change24h = 0.0 + (float64(time.Now().Unix()%500) / 100.0) - 2.5
volume24h = 100000000 + float64(time.Now().Unix()%10000000)
}

return &models.Price{
Coin:      coin,
Exchange:  "模拟交易所",
Price:     basePrice,
Change24h: change24h,
Volume24h: volume24h,
UpdatedAt: time.Now(),
}
}

// GetPrice 获取指定交易所的指定币种价格
func (s *PriceService) GetPrice(exchange, coin string) (*models.Price, error) {
s.mutex.RLock()
defer s.mutex.RUnlock()

if _, ok := s.prices[exchange]; !ok {
return nil, fmt.Errorf("交易所 %s 不存在", exchange)
}

if price, ok := s.prices[exchange][coin]; ok {
return price, nil
}

return nil, fmt.Errorf("币种 %s 在交易所 %s 上不存在", coin, exchange)
}

// GetAllPrices 获取所有价格
func (s *PriceService) GetAllPrices() map[string]map[string]*models.Price {
s.mutex.RLock()
defer s.mutex.RUnlock()

// 深拷贝价格数据以避免并发问题
result := make(map[string]map[string]*models.Price)
for exchange, coins := range s.prices {
result[exchange] = make(map[string]*models.Price)
for coin, price := range coins {
// 创建价格对象的副本
priceCopy := *price
result[exchange][coin] = &priceCopy
}
}

return result
}

// GetHistory 获取指定交易所指定币种的价格历史
func (s *PriceService) GetHistory(exchange, coin string, limit int) ([]*models.PriceHistory, error) {
s.mutex.RLock()
defer s.mutex.RUnlock()

if _, ok := s.history[exchange]; !ok {
return nil, fmt.Errorf("交易所 %s 不存在", exchange)
}

if history, ok := s.history[exchange][coin]; ok {
// 如果请求的历史记录数量超过实际记录数量，返回所有记录
if limit <= 0 || limit > len(history) {
limit = len(history)
}
// 返回最近的limit条记录
result := make([]*models.PriceHistory, limit)
for i := 0; i < limit; i++ {
// 创建历史记录的副本
historyCopy := *history[len(history)-i-1]
result[i] = &historyCopy
}
return result, nil
}

return nil, fmt.Errorf("币种 %s 在交易所 %s 上没有历史记录", coin, exchange)
}

// GetSupportedCoins 获取支持的币种列表
func (s *PriceService) GetSupportedCoins() []models.CoinInfo {
return models.SupportedCoins()
}

// GetSupportedExchanges 获取支持的交易所列表
func (s *PriceService) GetSupportedExchanges() []string {
s.mutex.RLock()
defer s.mutex.RUnlock()

exchanges := make([]string, 0, len(s.exchanges))
for _, exchange := range s.exchanges {
exchanges = append(exchanges, exchange.Name)
}
return exchanges
}
