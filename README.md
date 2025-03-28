# 加密货币价格追踪器

一个使用Go语言开发的实时加密货币价格追踪Web应用，可以监控多种加密货币的价格变动，并提供分析和提醒功能。

## 功能

- 实时价格监控：支持多种主流加密货币
- 数据可视化：直观的价格走势图表
- 警报系统：价格达到设定阈值时发送通知
- 历史数据：存储和分析历史价格数据

## 技术栈

- 后端：Go、Gin、SQLite
- 前端：HTML、CSS、JavaScript、Bootstrap、Chart.js

## 项目结构

```
crypto-tracker/
├── cmd/
│   └── server/         # 应用入口
├── internal/
│   ├── api/            # API处理函数
│   ├── models/         # 数据模型
│   ├── services/       # 业务逻辑
│   └── config/         # 配置管理
├── pkg/
│   ├── exchanges/      # 交易所API客户端
│   ├── database/       # 数据库连接和操作
│   └── utils/          # 通用工具函数
├── web/
│   ├── static/         # 静态资源
│   └── templates/      # HTML模板
└── configs/            # 配置文件
```

## 安装与运行

1. 确保安装了Go (版本 1.16+)
2. 克隆项目
3. 安装依赖：`go mod download`
4. 运行应用：`go run cmd/server/main.go`
5. 访问：http://localhost:8080

## API文档

应用提供以下API端点：

- `GET /api/v1/coins` - 获取支持的加密货币列表
- `GET /api/v1/prices` - 获取所有币种的当前价格
- `GET /api/v1/prices/{coin}` - 获取指定币种的价格
- `GET /api/v1/history/{coin}?period={period}` - 获取历史价格数据
- `POST /api/v1/alerts` - 创建新警报
- `GET /api/v1/alerts` - 获取所有警报
- `PUT /api/v1/alerts/{id}` - 更新警报
- `DELETE /api/v1/alerts/{id}` - 删除警报

# Crypto-tracker

