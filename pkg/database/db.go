package database

import (
"database/sql"
"fmt"
"os"
"path/filepath"

_ "github.com/mattn/go-sqlite3"
)

// 数据库表定义常量
const (
createPricesTable = `
CREATE TABLE IF NOT EXISTS prices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    coin TEXT NOT NULL,
    exchange TEXT NOT NULL,
    price REAL NOT NULL,
    change24h REAL NOT NULL,
    volume24h REAL NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(coin, exchange)
);`

createPriceHistoryTable = `
CREATE TABLE IF NOT EXISTS price_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    coin TEXT NOT NULL,
    exchange TEXT NOT NULL,
    price REAL NOT NULL,
    timestamp TIMESTAMP NOT NULL
);`

createAlertsTable = `
CREATE TABLE IF NOT EXISTS alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    coin TEXT NOT NULL,
    exchange TEXT NOT NULL,
    type TEXT NOT NULL,
    threshold REAL NOT NULL,
    status TEXT NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL,
    triggered_at TIMESTAMP
);`
)

// InitDB 初始化数据库连接和表结构
func InitDB(dbPath string) (*sql.DB, error) {
// 确保数据库目录存在
dbDir := filepath.Dir(dbPath)
if _, err := os.Stat(dbDir); os.IsNotExist(err) {
if err := os.MkdirAll(dbDir, 0755); err != nil {
return nil, fmt.Errorf("创建数据库目录失败: %v", err)
}
}

// 连接数据库
db, err := sql.Open("sqlite3", dbPath)
if err != nil {
return nil, fmt.Errorf("连接数据库失败: %v", err)
}

// 测试连接
if err := db.Ping(); err != nil {
return nil, fmt.Errorf("测试数据库连接失败: %v", err)
}

// 创建表
if err := createTables(db); err != nil {
return nil, fmt.Errorf("创建数据库表失败: %v", err)
}

return db, nil
}

// createTables 创建所需的数据库表
func createTables(db *sql.DB) error {
statements := []string{
createPricesTable,
createPriceHistoryTable,
createAlertsTable,
}

for _, stmt := range statements {
_, err := db.Exec(stmt)
if err != nil {
return err
}
}

return nil
}
