# Chapter 6.1: GORM

歡迎來到 ORM 的世界！在本節中，我們將學習 Go 生態系中最受歡迎的 ORM 函式庫——GORM。使用 GORM，您可以用更 Go-like 的方式與資料庫互動，而不需要手動編寫大量的 SQL 語句。

## 什麼是 ORM？

**ORM (Object-Relational Mapping)** 是一種程式設計技術，用於在「物件」和「關聯式資料庫」之間進行概念和資料的轉換。簡單來說，它允許我們用程式語言中的物件（例如 Go 的 struct）來對應資料庫中的表 (table)，並透過操作物件來達到操作資料表的目的。

- **優點**: 提高開發效率、減少重複的 SQL 程式碼、提供一定程度的資料庫抽象化。
- **缺點**: 可能有效能損耗、對於複雜查詢的掌握度不如原生 SQL 靈活。

## 開始使用 GORM

首先，需要安裝 GORM 以及對應的資料庫驅動。這裡我們以 SQLite 為例，因為它最簡單，不需要額外架設資料庫伺服器。

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```

## GORM 核心操作

### 1. 定義 Model

在 GORM 中，`Model` 是一個 Go struct，對應到資料庫中的一張表。通常我們會嵌入 `gorm.Model`，它已經包含了 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` 這幾個常用欄位。

```go
package main

import "gorm.io/gorm"

// Product 對應到 products 資料表
type Product struct {
  gorm.Model
  Code  string
  Price uint
}
```

### 2. 連線資料庫

```go
package main

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

func main() {
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // 自動遷移 schema，保持資料庫結構最新
  db.AutoMigrate(&Product{})
}
```

### 3. 核心 CRUD 操作

- **Create (新增)**
  `db.Create(&Product{Code: "D42", Price: 100})`

- **Read (查詢)**
  `var product Product`
  `db.First(&product, 1) // 根據 integer primary key 查詢`
  `db.First(&product, "code = ?", "D42") // 查詢 code 為 D42 的紀錄`

- **Update (更新)**
  `db.Model(&product).Update("Price", 200)`

- **Delete (刪除)**
  `db.Delete(&product, 1)`

## 結語與提示

GORM 提供了強大而便利的方式來操作資料庫，讓開發者能更專注於業務邏輯。它還有更多進階功能，例如 Preload/Joins、Transaction、Hooks 等，等待您進一步探索。

在下一個章節，我們將深入比較使用原生 SQL 和 ORM 的各自優劣，幫助您在不同情境下做出最合適的技術選擇。
