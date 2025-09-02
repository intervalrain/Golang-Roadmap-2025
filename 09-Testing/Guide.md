# Chapter 9: 測試

測試是軟體開發中確保程式碼品質和可靠性的關鍵實務。Go 語言內建了強大的測試框架，讓開發者能夠輕鬆撰寫各種類型的測試。本章將深入探討 Go 的測試生態系統，從基本的單元測試到進階的效能測試，幫助你建立全面的測試策略，確保程式碼的穩定性和可維護性。

## Unit Testing

**Unit Testing** (單元測試) 是測試個別程式單元（通常是函式或方法）的過程，確保每個單元都能正確執行預期的功能。

### Go 測試基礎

Go 的測試系統建立在 `testing` 套件之上，遵循簡單的命名約定：

- 測試檔案必須以 `_test.go` 結尾
- 測試函式必須以 `Test` 開頭
- 測試函式必須接受 `*testing.T` 參數

```go
// calculator.go
package calculator

func Add(a, b int) int {
    return a + b
}

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

```go
// calculator_test.go
package calculator

import (
    "testing"
)

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}

func TestDivide(t *testing.T) {
    // 測試正常情況
    result, err := Divide(10, 2)
    if err != nil {
        t.Errorf("Divide(10, 2) returned error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(10, 2) = %d; want 5", result)
    }
    
    // 測試錯誤情況
    _, err = Divide(10, 0)
    if err == nil {
        t.Error("Divide(10, 0) should return error")
    }
}
```

### 測試執行

```bash
# 執行當前套件的所有測試
go test

# 執行特定測試
go test -run TestAdd

# 顯示詳細輸出
go test -v

# 執行所有子套件的測試
go test ./...
```

## Table-driven Tests

**Table-driven Tests** 是 Go 社群廣泛採用的測試模式，透過資料表定義多組測試案例，提高測試的可讀性和維護性。

### 基本 Table-driven Test

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 5, 5},
        {"negative and positive", -3, 7, 4},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### 進階 Table-driven Test

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name        string
        a, b        int
        expected    int
        expectError bool
        errorMsg    string
    }{
        {"normal division", 10, 2, 5, false, ""},
        {"division by zero", 10, 0, 0, true, "division by zero"},
        {"negative result", -10, 2, -5, false, ""},
        {"both negative", -10, -2, 5, false, ""},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)
            
            if tt.expectError {
                if err == nil {
                    t.Errorf("Divide(%d, %d) expected error but got none", tt.a, tt.b)
                    return
                }
                if err.Error() != tt.errorMsg {
                    t.Errorf("Divide(%d, %d) error = %v; want %v", tt.a, tt.b, err.Error(), tt.errorMsg)
                }
            } else {
                if err != nil {
                    t.Errorf("Divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
                    return
                }
                if result != tt.expected {
                    t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
                }
            }
        })
    }
}
```

## Mocking

**Mocking** 是創建假物件來模擬真實依賴的技術，讓測試能夠獨立執行，不受外部系統影響。

### 介面和依賴注入

```go
// 定義介面
type UserRepository interface {
    GetUser(id int) (*User, error)
    SaveUser(user *User) error
}

type User struct {
    ID   int
    Name string
}

// 服務層
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUserName(id int) (string, error) {
    user, err := s.repo.GetUser(id)
    if err != nil {
        return "", err
    }
    return user.Name, nil
}
```

### 手動 Mock 實作

```go
// 測試用的 Mock 實作
type MockUserRepository struct {
    users map[int]*User
    err   error
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users: make(map[int]*User),
    }
}

func (m *MockUserRepository) GetUser(id int) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    
    if user, exists := m.users[id]; exists {
        return user, nil
    }
    
    return nil, errors.New("user not found")
}

func (m *MockUserRepository) SaveUser(user *User) error {
    if m.err != nil {
        return m.err
    }
    
    m.users[user.ID] = user
    return nil
}

// 設定 Mock 行為的輔助方法
func (m *MockUserRepository) SetUser(id int, name string) {
    m.users[id] = &User{ID: id, Name: name}
}

func (m *MockUserRepository) SetError(err error) {
    m.err = err
}
```

### 使用 Mock 進行測試

```go
func TestUserService_GetUserName(t *testing.T) {
    tests := []struct {
        name        string
        userID      int
        setupMock   func(*MockUserRepository)
        expected    string
        expectError bool
    }{
        {
            name:   "existing user",
            userID: 1,
            setupMock: func(mock *MockUserRepository) {
                mock.SetUser(1, "John Doe")
            },
            expected:    "John Doe",
            expectError: false,
        },
        {
            name:   "user not found",
            userID: 999,
            setupMock: func(mock *MockUserRepository) {
                // 不設定任何使用者
            },
            expected:    "",
            expectError: true,
        },
        {
            name:   "repository error",
            userID: 1,
            setupMock: func(mock *MockUserRepository) {
                mock.SetError(errors.New("database connection failed"))
            },
            expected:    "",
            expectError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 設定 Mock
            mockRepo := NewMockUserRepository()
            tt.setupMock(mockRepo)
            
            // 建立服務
            service := NewUserService(mockRepo)
            
            // 執行測試
            result, err := service.GetUserName(tt.userID)
            
            // 驗證結果
            if tt.expectError {
                if err == nil {
                    t.Error("expected error but got none")
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if result != tt.expected {
                    t.Errorf("got %s; want %s", result, tt.expected)
                }
            }
        })
    }
}
```

## Integration Testing

**Integration Testing** (整合測試) 驗證多個元件或系統之間的交互作用是否正常運作。

### 資料庫整合測試

```go
func TestUserRepository_Integration(t *testing.T) {
    // 設定測試資料庫
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewUserRepository(db)
    
    // 測試儲存使用者
    user := &User{Name: "Test User"}
    err := repo.SaveUser(user)
    if err != nil {
        t.Fatalf("SaveUser failed: %v", err)
    }
    
    // 測試查詢使用者
    retrievedUser, err := repo.GetUser(user.ID)
    if err != nil {
        t.Fatalf("GetUser failed: %v", err)
    }
    
    if retrievedUser.Name != user.Name {
        t.Errorf("got %s; want %s", retrievedUser.Name, user.Name)
    }
}

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }
    
    // 建立測試表格
    _, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        )
    `)
    if err != nil {
        t.Fatalf("Failed to create test table: %v", err)
    }
    
    return db
}
```

### HTTP API 整合測試

```go
func TestUserAPI_Integration(t *testing.T) {
    // 設定測試伺服器
    server := setupTestServer(t)
    defer server.Close()
    
    tests := []struct {
        name           string
        method         string
        path           string
        body           string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "create user",
            method:         "POST",
            path:           "/users",
            body:           `{"name": "John Doe"}`,
            expectedStatus: http.StatusCreated,
            expectedBody:   `{"id":1,"name":"John Doe"}`,
        },
        {
            name:           "get user",
            method:         "GET",
            path:           "/users/1",
            expectedStatus: http.StatusOK,
            expectedBody:   `{"id":1,"name":"John Doe"}`,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var req *http.Request
            var err error
            
            if tt.body != "" {
                req, err = http.NewRequest(tt.method, server.URL+tt.path, strings.NewReader(tt.body))
                req.Header.Set("Content-Type", "application/json")
            } else {
                req, err = http.NewRequest(tt.method, server.URL+tt.path, nil)
            }
            
            if err != nil {
                t.Fatalf("Failed to create request: %v", err)
            }
            
            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                t.Fatalf("Request failed: %v", err)
            }
            defer resp.Body.Close()
            
            if resp.StatusCode != tt.expectedStatus {
                t.Errorf("got status %d; want %d", resp.StatusCode, tt.expectedStatus)
            }
            
            body, _ := io.ReadAll(resp.Body)
            if strings.TrimSpace(string(body)) != tt.expectedBody {
                t.Errorf("got body %s; want %s", string(body), tt.expectedBody)
            }
        })
    }
}
```

## Benchmark Testing

**Benchmark Testing** (效能測試) 用於測量程式碼的執行效能，幫助識別效能瓶頸。

### 基本 Benchmark

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(42, 24)
    }
}

func BenchmarkStringConcatenation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        result := ""
        for j := 0; j < 100; j++ {
            result += "hello"
        }
    }
}

func BenchmarkStringBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        for j := 0; j < 100; j++ {
            builder.WriteString("hello")
        }
        _ = builder.String()
    }
}
```

### 進階 Benchmark 技巧

```go
func BenchmarkMap(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            m := make(map[int]int, size)
            
            // 準備測試資料
            for i := 0; i < size; i++ {
                m[i] = i * 2
            }
            
            b.ResetTimer() // 重設計時器，排除準備時間
            
            for i := 0; i < b.N; i++ {
                key := i % size
                _ = m[key]
            }
        })
    }
}

// 測記憶體分配
func BenchmarkSliceAppend(b *testing.B) {
    b.ReportAllocs() // 報告記憶體分配統計
    
    for i := 0; i < b.N; i++ {
        var slice []int
        for j := 0; j < 1000; j++ {
            slice = append(slice, j)
        }
    }
}
```

### 執行 Benchmark

```bash
# 執行所有 Benchmark
go test -bench=.

# 執行特定 Benchmark
go test -bench=BenchmarkAdd

# 執行多次取平均值
go test -bench=. -count=5

# 顯示記憶體分配統計
go test -bench=. -benchmem

# 設定執行時間
go test -bench=. -benchtime=10s
```

## Test Coverage

**Test Coverage** (測試覆蓋率) 衡量測試執行時程式碼被執行的比例，幫助識別未測試的程式碼區域。

### 生成覆蓋率報告

```bash
# 執行測試並生成覆蓋率報告
go test -cover

# 生成詳細的覆蓋率檔案
go test -coverprofile=coverage.out

# 查看覆蓋率詳情
go tool cover -func=coverage.out

# 生成 HTML 覆蓋率報告
go tool cover -html=coverage.out -o coverage.html
```

### 設定覆蓋率目標

```go
// 使用建構標籤區分測試和生產程式碼
//go:build !test

package main

// 生產程式碼
```

```go
//go:build test

package main

// 測試相關程式碼
```

## 測試最佳實務

### 測試命名規範

```go
// 好的測試命名
func TestUserService_GetUser_ValidID_ReturnsUser(t *testing.T) {}
func TestUserService_GetUser_InvalidID_ReturnsError(t *testing.T) {}
func TestUserService_GetUser_RepositoryError_ReturnsError(t *testing.T) {}

// 使用子測試提高可讀性
func TestUserService_GetUser(t *testing.T) {
    t.Run("valid ID returns user", func(t *testing.T) {})
    t.Run("invalid ID returns error", func(t *testing.T) {})
    t.Run("repository error returns error", func(t *testing.T) {})
}
```

### 測試隔離和清理

```go
func TestUserService(t *testing.T) {
    t.Run("create user", func(t *testing.T) {
        // 設定
        service := setupUserService(t)
        defer cleanup(t, service)
        
        // 執行和驗證測試
    })
}

func setupUserService(t *testing.T) *UserService {
    // 建立測試依賴
    return &UserService{}
}

func cleanup(t *testing.T, service *UserService) {
    // 清理資源
}
```

### 輔助函式

```go
// 測試輔助函式
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper() // 標記為輔助函式
    if got != want {
        t.Errorf("got %v; want %v", got, want)
    }
}

func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

func assertError(t *testing.T, err error, msg string) {
    t.Helper()
    if err == nil {
        t.Error("expected error but got none")
        return
    }
    if err.Error() != msg {
        t.Errorf("got error %v; want %v", err.Error(), msg)
    }
}
```

## 進階測試工具

### Testify 框架

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestWithTestify(t *testing.T) {
    // 使用 assert 簡化測試
    result := Add(2, 3)
    assert.Equal(t, 5, result, "Add should return correct result")
    
    // 測試錯誤
    _, err := Divide(10, 0)
    assert.Error(t, err, "Division by zero should return error")
    assert.Contains(t, err.Error(), "division by zero")
}
```

### 測試環境設定

```go
func TestMain(m *testing.M) {
    // 設定測試環境
    setupTestEnvironment()
    
    // 執行所有測試
    code := m.Run()
    
    // 清理測試環境
    cleanupTestEnvironment()
    
    // 結束程式
    os.Exit(code)
}
```

透過本章的學習，你已經掌握了 Go 語言測試的完整知識體系。從基本的單元測試到進階的效能測試，從簡單的測試案例到複雜的整合測試，這些技能將幫助你建構更加穩定可靠的應用程式。記住，良好的測試不僅是程式碼品質的保證，更是重構和維護的信心來源。接下來，建議你將這些測試技巧應用到實際專案中，建立完整的測試套件。