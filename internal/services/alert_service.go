package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-tracker/internal/models"
)

// AlertService 处理价格警报
type AlertService struct {
	db *sql.DB
}

// NewAlertService 创建一个新的警报服务
func NewAlertService(db *sql.DB) *AlertService {
	return &AlertService{
		db: db,
	}
}

// CreateAlert 创建新警报
func (s *AlertService) CreateAlert(req models.AlertRequest) (*models.Alert, error) {
	alert := &models.Alert{
		Coin:      req.Coin,
		Exchange:  req.Exchange,
		Type:      req.Type,
		Threshold: req.Threshold,
		Status:    models.Active,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	// 在实际的实现中，这里应该将警报保存到数据库
	// 这里为了简单起见，我们只返回创建的警报对象
	// 正常情况下应当执行SQL插入操作
	alert.ID = time.Now().Unix() // 生成一个模拟ID

	return alert, nil
}

// GetAlerts 获取所有警报
func (s *AlertService) GetAlerts() ([]*models.Alert, error) {
	// 在实际实现中，应该从数据库中查询所有警报
	// 这里仅做示例返回一些模拟数据
	alerts := []*models.Alert{
		{
			ID:        1,
			Coin:      "BTC",
			Exchange:  "Binance",
			Type:      models.PriceAbove,
			Threshold: 50000.0,
			Status:    models.Active,
			Message:   "比特币价格超过50000美元",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        2,
			Coin:      "ETH",
			Exchange:  "Coinbase",
			Type:      models.PriceBelow,
			Threshold: 2000.0,
			Status:    models.Active,
			Message:   "以太坊价格低于2000美元",
			CreatedAt: time.Now().Add(-12 * time.Hour),
		},
	}

	return alerts, nil
}

// GetAlertByID 通过ID获取警报
func (s *AlertService) GetAlertByID(id int64) (*models.Alert, error) {
	// 模拟数据，实际应该从数据库查询
	if id == 1 {
		return &models.Alert{
			ID:        1,
			Coin:      "BTC",
			Exchange:  "Binance",
			Type:      models.PriceAbove,
			Threshold: 50000.0,
			Status:    models.Active,
			Message:   "比特币价格超过50000美元",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		}, nil
	}

	return nil, fmt.Errorf("警报ID %d 不存在", id)
}

// UpdateAlert 更新警报
func (s *AlertService) UpdateAlert(id int64, req models.AlertRequest) (*models.Alert, error) {
	alert, err := s.GetAlertByID(id)
	if err != nil {
		return nil, err
	}

	// 更新警报属性
	alert.Coin = req.Coin
	alert.Exchange = req.Exchange
	alert.Type = req.Type
	alert.Threshold = req.Threshold
	alert.Message = req.Message

	// 在实际实现中，应该执行数据库更新操作
	return alert, nil
}

// DeleteAlert 删除警报
func (s *AlertService) DeleteAlert(id int64) error {
	// 确认警报存在
	_, err := s.GetAlertByID(id)
	if err != nil {
		return err
	}

	// 在实际实现中，应该执行数据库删除操作
	return nil
}

// CheckAlerts 检查所有警报是否触发
func (s *AlertService) CheckAlerts(priceService *PriceService) ([]*models.Alert, error) {
	alerts, err := s.GetAlerts()
	if err != nil {
		return nil, err
	}

	var triggeredAlerts []*models.Alert

	for _, alert := range alerts {
		// 跳过非活跃警报
		if alert.Status != models.Active {
			continue
		}

		// 获取当前价格
		price, err := priceService.GetPrice(alert.Exchange, alert.Coin)
		if err != nil {
			continue // 忽略无法获取价格的币种
		}

		// 检查警报条件
		var triggered bool
		switch alert.Type {
		case models.PriceAbove:
			triggered = price.Price > alert.Threshold
		case models.PriceBelow:
			triggered = price.Price < alert.Threshold
		case models.ChangeAbove:
			triggered = price.Change24h > alert.Threshold
		case models.ChangeBelow:
			triggered = price.Change24h < alert.Threshold
		}

		// 如果警报触发
		if triggered {
			now := time.Now()
			alert.Status = models.Triggered
			alert.TriggeredAt = &now
			triggeredAlerts = append(triggeredAlerts, alert)

			// 在实际实现中，应该更新数据库中的警报状态
		}
	}

	return triggeredAlerts, nil
}
