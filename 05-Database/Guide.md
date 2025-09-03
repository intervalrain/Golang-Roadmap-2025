# Chapter 5.1: SQL Basics

歡迎來到資料庫的第一站！在深入研究 Go 如何與資料庫互動之前，我們必須先掌握與資料庫溝通的語言——SQL。本節將介紹 SQL 的基礎知識，為後續的學習奠定堅實的基礎。

## 什麼是 SQL？

SQL (Structured Query Language) 是一種標準化的程式語言，專門用於管理**關聯式資料庫 (Relational Databases)** 並對其中的資料進行操作。無論您是想新增使用者、查詢訂單，還是更新文章，都離不開 SQL。

## 核心 SQL 指令

SQL 指令主要可以分為幾大類，我們來看看最常見的幾種。

### Data Query Language (DQL)

- **`SELECT`**: 用於從資料庫中查詢資料。
  - `SELECT name, age FROM users;` (查詢 `users` 表中的 `name` 和 `age` 欄位)
  - `SELECT * FROM users WHERE age > 30;` (查詢所有年齡大於 30 的使用者)

### Data Manipulation Language (DML)

- **`INSERT`**: 用於將新資料新增到資料庫表中。
  - `INSERT INTO users (name, email) VALUES ('John Doe', 'john.doe@example.com');`

- **`UPDATE`**: 用於修改資料庫表中已存在的資料。
  - `UPDATE users SET email = 'new.email@example.com' WHERE name = 'John Doe';`

- **`DELETE`**: 用於從資料庫表中刪除資料。
  - `DELETE FROM users WHERE name = 'John Doe';`

### Data Definition Language (DDL)

- **`CREATE TABLE`**: 用於建立新的資料庫表。
  ```sql
  CREATE TABLE users (
      id SERIAL PRIMARY KEY,
      name VARCHAR(100) NOT NULL,
      email VARCHAR(100) UNIQUE,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
  );
  ```

- **`ALTER TABLE`**: 用於修改現有資料庫表的結構，例如新增或刪除欄位。
  - `ALTER TABLE users ADD COLUMN age INT;`

- **`DROP TABLE`**: 用於刪除整個資料庫表。
  - `DROP TABLE users;

## 關鍵概念

- **`Primary Key` (主鍵)**: 表中每一筆資料的**唯一識別符**。例如 `users` 表中的 `id`。
- **`Foreign Key` (外鍵)**: 一個表中用於連結到另一個表主鍵的欄位，用於建立表與表之間的關聯。
- **`Index` (索引)**: 用於**加速查詢**操作的資料結構。可以把它想像成書本的目錄，讓資料庫能更快地找到所需的資料。

## 結語與提示

我們剛剛快速瀏覽了 SQL 的基礎。雖然看起來很簡單，但 SQL 的世界非常廣闊，組合查詢、JOIN 操作、交易處理等都是更進階的主題。

在下一個章節，我們將學習如何使用 Go 內建的 `database/sql` 套件來執行這些 SQL 指令。

---

# Chapter 5.2: database/sql Package

## 前言
Go 透過內建的 `database/sql` 套件提供了一個與 SQL 資料庫互動的標準介面。這個套件本身不包含任何特定資料庫的驅動程式，而是定義了一組通用的、抽象的介面。要連接到特定的資料庫（如 PostgreSQL 或 MySQL），您需要匯入對應的資料庫驅動程式。

## 核心介面
- **`sql.DB`**: 代表一個資料庫連線池，管理著與資料庫的底層連線。它是併發安全的，應該在應用程式的生命週期中被多個 Goroutine 共享。
- **`sql.Tx`**: 代表一個資料庫交易。
- **`sql.Stmt`**: 代表一個預備好的 SQL 語句 (prepared statement)，可以重複執行。
- **`sql.Rows`**: 代表查詢結果的迭代器。
- **`sql.Row`**: 代表單行查詢結果。

## 連接與查詢範例
首先，您需要匯入 `database/sql` 套件以及您選擇的資料庫驅動程式。注意，驅動程式通常是匿名匯入 (`_`)，因為我們只需要它向 `database/sql` 註冊自己。

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq" // 匯入 PostgreSQL 驅動
)

func main() {
    // 連接字串
    connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 測試連線
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    fmt.Println("Successfully connected!")

    // 執行查詢
    rows, err := db.Query("SELECT id, name FROM users WHERE id = $1", 1)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // 迭代結果
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("ID: %d, Name: %s\n", id, name)
    }
}
```

## 防止 SQL 注入
**永遠不要**使用 `fmt.Sprintf` 或字串拼接來組合 SQL 查詢。請一律使用參數化查詢（如上例中的 `$1`），`database/sql` 套件會為您安全地處理參數，從而防止 SQL 注入攻擊。

---

# Chapter 5.3: PostgreSQL

## 前言
**PostgreSQL** 是一個功能強大、開源的物件關聯式資料庫系統，以其可靠性、功能完整性和效能而聞名。它是許多大型企業和新創公司的首選資料庫。

## Go 驅動程式
- **`github.com/lib/pq`**: 純 Go 實現的 PostgreSQL 驅動，非常流行且穩定。
- **`github.com/jackc/pgx`**: 一個效能更高、功能更豐富的驅動程式，提供了更多 PostgreSQL 的原生功能。對於新專案，`pgx` 通常是更好的選擇。

## 連接字串
一個典型的 PostgreSQL 連接字串 (Connection String) 包含以下部分：
`"postgres://user:password@host:port/dbname?sslmode=disable"`

---

# Chapter 5.4: MySQL

## 前言
**MySQL** 是世界上最流行的開源資料庫之一，特別是在 Web 開發領域。它由 Oracle 公司維護，以其易用性和高效能而著稱。

## Go 驅動程式
- **`github.com/go-sql-driver/mysql`**: 這是最常用、最受推薦的 MySQL 驅動程式。它完全相容 `database/sql` 介面。

## 連接字串 (DSN)
MySQL 的連接字串通常被稱為 DSN (Data Source Name)，格式如下：
`"user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"`

- `parseTime=True`: 讓驅動程式能將資料庫中的 `TIME`, `DATE`, `DATETIME` 類型自動轉換為 Go 的 `time.Time`。
- `charset=utf8mb4`: 確保能正確處理包含表情符號等特殊字元的字串。

---

# Chapter 5.5: SQLite

## 前言
**SQLite** 是一個無伺服器、零設定、交易性的 SQL 資料庫引擎。它將整個資料庫（包含定義、表、索引和資料）儲存在一個單一的檔案中。這使得它非常適合用於行動應用、桌面應用或簡單的 Web 應用。

## Go 驅動程式
- **`github.com/mattn/go-sqlite3`**: 一個基於 CGO 的驅動程式，是目前最主流的選擇。

## 使用場景
- **開發與測試:** 在開發和測試階段，使用 SQLite 可以快速啟動，無需設定一個完整的資料庫伺服器。
- **嵌入式應用:** 當您的應用程式需要一個輕量級的內嵌資料庫時。
- **資料分析:** 用於處理和查詢本地的資料集檔案。

---

# Chapter 5.6: NoSQL

## 前言
**NoSQL** (通常解釋為 "Not Only SQL") 泛指所有非關聯式的資料庫。它們提供了與傳統關聯式資料庫（如 PostgreSQL, MySQL）不同的資料儲存和檢索模型。NoSQL 資料庫在處理大量資料、高併發讀寫和非結構化資料方面表現出色。

## 主要類型
- **文件資料庫 (Document Databases):** 將資料儲存為類似 JSON 的文件。非常適合儲存非結構化或半結構化的資料。
  - **MongoDB:** 最流行的文件資料庫。Go 官方驅動為 `go.mongodb.org/mongo-driver`。
- **鍵/值儲存 (Key-Value Stores):** 資料以簡單的鍵/值對形式儲存。讀寫速度極快。
  - **Redis:** 一個高效能的記憶體內鍵/值儲存，常用作快取、訊息代理和即時排行榜。Go 驅動有 `github.com/redis/go-redis`。
- **寬列儲存 (Wide-Column Stores):** 例如 Cassandra, HBase。
- **圖形資料庫 (Graph Databases):** 例如 Neo4j。

## 何時選擇 NoSQL？
- **彈性的資料模型:** 當您的資料結構經常變動或沒有固定結構時。
- **高擴展性:** 當您需要水平擴展以處理大量流量時。
- **高效能:** 當您需要極低的讀寫延遲時（特別是鍵/值儲存）。

## 結語
恭喜您完成了資料庫章節的學習！您現在了解了 SQL 的基礎，掌握了如何使用 Go 的 `database/sql` 套件與 PostgreSQL、MySQL 和 SQLite 等關聯式資料庫互動，並對 MongoDB 和 Redis 等 NoSQL 資料庫有了初步的認識。

下一步，我們將探索 **ORM (Object-Relational Mapping)**，學習如何使用 GORM 等工具來簡化資料庫操作。