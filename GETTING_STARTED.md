# 加密货币价格追踪器 - 使用指南

## 项目简介

这是一个使用Go语言开发的实时加密货币价格追踪Web应用，可以监控多种加密货币的价格变动，并提供分析和提醒功能。

## 前提条件

1. 安装 Go (版本 1.16+)
2. 安装 GCC 编译器（SQLite3驱动需要）
   - Windows: 安装 [MinGW](https://sourceforge.net/projects/mingw-w64/)
   - macOS: `xcode-select --install`
   - Linux: `sudo apt-get install build-essential`

## 安装步骤

1. 克隆或下载项目到本地

2. 进入项目目录
   ```bash
   cd crypto-tracker
   ```

3. 安装依赖
   ```bash
   go mod download
   ```
   
   或者手动添加依赖：
   ```bash
   go get github.com/gin-gonic/gin
   go get github.com/mattn/go-sqlite3
   ```

## 运行应用

1. 从项目根目录启动应用
   ```bash
   go run cmd/server/main.go
   ```

2. 打开浏览器，访问以下地址：
   ```
   http://localhost:8080
   ```

## 功能说明

应用包括以下主要功能页面：

1. **仪表盘**：显示实时价格和市场概览
   - 实时价格表格，显示主要加密货币的当前价格
   - 价格趋势图表，支持多个时间范围选择

2. **价格警报**：设置价格警报
   - 创建、编辑和删除价格警报
   - 支持多种警报条件（价格高于/低于、涨跌幅高于/低于）

3. **历史数据**：查看价格历史
   - 选择币种和时间范围
   - 图表显示历史价格走势

## 注意事项

- 当前版本使用模拟数据，未连接真实交易所API
- 数据库存储在 `data/crypto.db`，首次运行时会自动创建
- 配置文件位于 `configs/config.json`，可根据需要修改设置

## 后续开发

1. 集成真实交易所API
2. 添加用户账户和认证
3. 实现邮件/短信警报通知
4. 改进数据分析功能
5. 添加移动端适配
