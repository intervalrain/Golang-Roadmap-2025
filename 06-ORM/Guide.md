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

---

# Chapter 6.2: Raw SQL vs ORM

## 前言
在選擇資料存取方式時，開發者經常面臨一個抉擇：是使用原生 SQL (`database/sql`) 還是採用 ORM (如 GORM)？這兩者各有優劣，了解它們的差異有助於您根據專案需求做出明智的決定。

## ORM 的優點
1.  **開發效率高:** ORM 將資料庫操作封裝成物件方法，減少了大量的 CRUD 樣板程式碼，讓開發者能更專注於業務邏輯。
2.  **可讀性與維護性:** 物件導向的語法通常比 SQL 字串更容易閱讀和維護。
3.  **資料庫遷移性:** 好的 ORM 提供了對多種資料庫的支援，理論上可以在不更改業務邏輯程式碼的情況下更換底層資料庫。
4.  **內建功能:** 許多 ORM 內建了如軟刪除 (Soft Deletes)、時間戳自動更新、關聯預載入 (Preloading) 等便利功能。

## 原生 SQL 的優點
1.  **極致的效能與靈活性:** 您可以完全控制執行的 SQL 語句，進行細緻的效能優化。對於非常複雜的查詢（例如多重 JOIN、子查詢、資料庫特有函式），原生 SQL 更具優勢。
2.  **無額外抽象層:** 所見即所得，沒有 ORM 這種中間層，更容易預測行為和進行偵錯。
3.  **學習曲線平緩:** 對於已經熟悉 SQL 的開發者來說，`database/sql` 的學習成本非常低。
4.  **輕量:** `database/sql` 是 Go 的標準函式庫，無需引入額外的第三方依賴。

## 如何選擇？
| 情境 | 建議選擇 | 原因 |
| :--- | :--- | :--- |
| **快速原型開發/新創公司** | ORM | 開發速度是關鍵，ORM 能快速建構功能。 |
| **標準 CRUD 密集的應用** | ORM | 大部分操作都是簡單的增刪改查，ORM 能顯著提升效率。 |
| **高效能要求的核心系統** | 原生 SQL | 需要對查詢進行極致優化，榨乾每一分效能。 |
| **複雜的資料分析/報表** | 原生 SQL | 查詢邏輯複雜多變，ORM 可能難以表達或效率低下。 |
| **團隊 SQL 掌握度高** | 原生 SQL | 團隊成員都是 SQL 專家，直接寫 SQL 更直接高效。 |

**混合使用也是一個常見且有效的策略**：在應用程式的大部分地方使用 ORM 來提高開發效率，對於少數需要高效能或複雜查詢的關鍵路徑，則使用原生 SQL 來處理。

---

# Chapter 6.3: Migrations

## 前言
**Database Migrations (資料庫遷移)** 是一種管理資料庫結構 (Schema) 變更的版本控制方法。隨著應用程式的迭代，您可能需要新增資料表、修改欄位、新增索引等。Migrations 讓您可以用程式碼的形式來描述這些變更，並在不同環境（開發、測試、生產）中以可預測、可重複的方式來應用它們。

## 為何需要 Migrations？
- **版本控制:** 將資料庫結構的變更納入 Git 等版本控制系統。
- **團隊協作:** 讓團隊中的每個成員都能輕鬆地將資料庫更新到最新結構。
- **自動化部署:** 在 CI/CD 流程中自動執行資料庫結構的更新。
- **可追溯與回滾:** 清楚地記錄了每一次變更，並在需要時可以撤銷 (回滾) 變更。

## GORM 的 AutoMigrate
GORM 提供了 `db.AutoMigrate(&YourModel{})` 功能。它會檢查您的 Go struct model，並自動在資料庫中建立或更新資料表、欄位、索引等，以使其與 model 保持一致。

- **優點:** 非常方便，適合在開發或原型階段快速迭代。
- **缺點:**
    - 它**只能新增**，不能安全地修改或刪除欄位（例如更改欄位類型或刪除欄位），因為這可能導致資料遺失。
    - 對於生產環境，這種「隱式」的變更不夠明確和可控。

## 專業的 Migration 工具
對於正式專案，特別是生產環境，建議使用專門的 migration 工具。
- **`golang-migrate/migrate`**: Go 生態系中最流行的 migration 工具之一。它支援從檔案或 Go 程式碼中讀取 SQL 腳本，並記錄每次 migration 的版本。
- **`pressly/goose`**: 另一個受歡迎的 migration 工具。

### `golang-migrate/migrate` 範例
1.  安裝 CLI 工具。
2.  建立 migration 檔案：
    `migrate create -ext sql -dir db/migrations -seq create_users_table`
    這會產生兩個檔案，一個 `up` (升級)，一個 `down` (降級)。
3.  編輯 `up.sql` 和 `down.sql` 檔案：
    ```sql
    -- 000001_create_users_table.up.sql
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL
    );

    -- 000001_create_users_table.down.sql
    DROP TABLE users;
    ```
4.  應用 migration：
    `migrate -database "postgres://..." -path db/migrations up`
5.  回滾 migration：
    `migrate -database "postgres://..." -path db/migrations down 1` (回滾一個版本)

---

# Chapter 6.4: Database Design

## 前言
良好的資料庫設計是建構穩健、高效能應用程式的基石。它不僅關乎資料的儲存，更影響到資料的完整性、查詢效率和未來的可擴展性。

## 核心原則
1.  **正規化 (Normalization):**
    - 這是減少資料冗餘和提高資料完整性的一系列準則。主要目標是確保資料庫中的每個「事實」只被儲存一次。
    - **第一正規化 (1NF):** 確保所有欄位都是不可分割的原子值。
    - **第二正規化 (2NF):** 滿足 1NF，且所有非主鍵欄位完全依賴於主鍵。
    - **第三正規化 (3NF):** 滿足 2NF，且所有非主鍵欄位不互相依賴。
    - **實務建議:** 在大多數應用中，達到 3NF 是一個很好的平衡點。過度的正規化有時會導致查詢需要 JOIN 過多的資料表，反而降低效能。

2.  **選擇正確的資料類型:**
    - 為每個欄位選擇最合適、最節省空間的資料類型。例如，如果一個整數欄位的範圍不會超過 32767，就使用 `SMALLINT` 而不是 `INTEGER`。
    - 對於價格等需要精確計算的數字，使用 `DECIMAL` 或 `NUMERIC`，而不是 `FLOAT`。

3.  **索引策略 (Indexing):**
    - **謹慎地新增索引:** 索引可以極大地提升 `SELECT` 查詢的速度，但會降低 `INSERT`, `UPDATE`, `DELETE` 的速度，因為每次資料變更時都需要更新索引。
    - **為 `WHERE` 子句、`JOIN` 操作和 `ORDER BY` 子句中頻繁使用的欄位建立索引。**
    - **使用複合索引 (Composite Indexes):** 當查詢經常同時涉及多個欄位時，可以為這些欄位建立一個複合索引。

4.  **關聯 (Relationships):**
    - **一對一 (One-to-One):** 例如 `users` 和 `user_profiles`。
    - **一對多 (One-to-Many):** 例如一個 `author` 可以有多篇 `posts`。通常在「多」的一方（`posts`）新增一個外鍵 (`author_id`)。
    - **多對多 (Many-to-Many):** 例如 `posts` 和 `tags`。一篇文章可以有多個標籤，一個標籤也可以對應多篇文章。這需要一個中間的「連接表」(pivot table)，例如 `post_tags`，它包含 `post_id` 和 `tag_id`。

## 結語
恭喜您完成了 ORM 與資料庫設計章節！您現在不僅知道如何使用 GORM 快速開發，也了解了原生 SQL 的強大之處，並掌握了透過資料庫遷移來管理資料庫結構的專業方法，以及良好資料庫設計的核心原則。

下一步，我們將進入 **API 開發** 的世界，學習如何設計和建構 RESTful API。