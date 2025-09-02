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
  - `DROP TABLE users;`

## 關鍵概念

- **`Primary Key` (主鍵)**: 表中每一筆資料的**唯一識別符**。例如 `users` 表中的 `id`。
- **`Foreign Key` (外鍵)**: 一個表中用於連結到另一個表主鍵的欄位，用於建立表與表之間的關聯。
- **`Index` (索引)**: 用於**加速查詢**操作的資料結構。可以把它想像成書本的目錄，讓資料庫能更快地找到所需的資料。

## 結語與提示

我們剛剛快速瀏覽了 SQL 的基礎。雖然看起來很簡單，但 SQL 的世界非常廣闊，組合查詢、JOIN 操作、交易處理等都是更進階的主題。

在下一個章節，我們將學習如何使用 Go 內建的 `database/sql` 套件來執行這些 SQL 指令。
