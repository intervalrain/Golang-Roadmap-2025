# Chapter 7: API 開發

在現代的軟體開發中，**API** (Application Programming Interface) 是系統間溝通的橋樑。Go 語言憑藉其高效能和簡潔的語法，成為開發 API 的絕佳選擇。本章將深入探討如何使用 Go 建構穩健、可擴展的 API，涵蓋從基本的 RESTful 設計到進階的 GraphQL 和 gRPC 實作。

## RESTful API Design

**REST** (Representational State Transfer) 是一種架構風格，定義了如何設計網路應用程式的介面。RESTful API 遵循以下核心原則：

### 核心原則

- **資源導向 (Resource-Oriented):** 每個 URL 代表一個資源
- **統一介面 (Uniform Interface):** 使用標準的 HTTP 方法
- **無狀態 (Stateless):** 每個請求都包含完整的資訊
- **可快取 (Cacheable):** 回應可以被快取以提升效能

### HTTP 方法對應

```go
// GET - 取得資源
GET /api/users          // 取得所有使用者
GET /api/users/123      // 取得特定使用者

// POST - 建立新資源
POST /api/users         // 建立新使用者

// PUT - 完全更新資源
PUT /api/users/123      // 完全更新使用者資訊

// PATCH - 部分更新資源
PATCH /api/users/123    // 部分更新使用者資訊

// DELETE - 刪除資源
DELETE /api/users/123   // 刪除使用者
```

### 狀態碼最佳實務

- **2xx 成功:** 200 (OK), 201 (Created), 204 (No Content)
- **4xx 客戶端錯誤:** 400 (Bad Request), 401 (Unauthorized), 404 (Not Found)
- **5xx 伺服器錯誤:** 500 (Internal Server Error), 503 (Service Unavailable)

## CRUD Operations

**CRUD** 代表 Create、Read、Update、Delete 四個基本資料操作。以下是使用 Go 實作的範例：

### 基本 CRUD 結構

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserHandler struct {
    users map[int]User
    nextID int
}

// Create - 建立使用者
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)
    
    user.ID = h.nextID
    h.nextID++
    h.users[user.ID] = user
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// Read - 取得使用者
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Path[len("/api/users/"):])
    
    if user, exists := h.users[id]; exists {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
}
```

## API Documentation

良好的 API 文件是成功 API 的關鍵。文件應該包含：

### 必要資訊
- **端點描述:** 每個 API 端點的功能說明
- **請求格式:** 參數類型、必填欄位、範例
- **回應格式:** 成功和錯誤回應的結構
- **認證方式:** 如何進行身份驗證
- **錯誤碼:** 可能的錯誤狀況和處理方式

### 文件工具推薦
- **Swagger/OpenAPI:** 業界標準的 API 規範
- **Postman Collections:** 互動式 API 測試
- **README.md:** 簡單的文字說明

## Swagger/OpenAPI

**OpenAPI** (原稱 Swagger) 是一個 API 規範標準，提供機器可讀的 API 描述。

### Go 中的 Swagger 整合

```go
// 使用 swaggo/swag 套件
//go:generate swag init

// @title User API
// @version 1.0
// @description 使用者管理 API

// @host localhost:8080
// @BasePath /api

// @Summary 取得所有使用者
// @Description 取得系統中所有使用者的清單
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
    // 實作邏輯
}
```

### 自動文件生成

```bash
# 安裝 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成 Swagger 文件
swag init
```

## GraphQL

**GraphQL** 是一種查詢語言和執行環境，讓客戶端可以精確指定需要的資料。

### GraphQL vs REST

- **單一端點:** GraphQL 只需要一個 URL
- **精確查詢:** 客戶端指定需要的欄位
- **類型系統:** 強型別的 Schema 定義
- **即時更新:** 支援 Subscription 功能

### Go GraphQL 實作

```go
// 使用 graphql-go/graphql 套件
type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var userType = graphql.NewObject(graphql.ObjectConfig{
    Name: "User",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.String,
        },
        "name": &graphql.Field{
            Type: graphql.String,
        },
    },
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Query",
    Fields: graphql.Fields{
        "user": &graphql.Field{
            Type: userType,
            Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                    Type: graphql.String,
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                // 查詢邏輯
                return User{ID: "1", Name: "John"}, nil
            },
        },
    },
})
```

## gRPC

**gRPC** 是由 Google 開發的高效能 RPC 框架，使用 Protocol Buffers 作為序列化格式。

### gRPC 優勢

- **高效能:** 使用 HTTP/2 和 Protocol Buffers
- **跨語言:** 支援多種程式語言
- **類型安全:** 透過 .proto 檔案定義介面
- **串流支援:** 支援單向和雙向串流

### Protocol Buffer 定義

```proto
syntax = "proto3";

package user;

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (stream User);
}

message User {
    int32 id = 1;
    string name = 2;
    string email = 3;
}

message GetUserRequest {
    int32 id = 1;
}
```

### Go gRPC 實作

```go
// 實作 gRPC 服務
type userServiceServer struct {
    users map[int32]*pb.User
}

func (s *userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    if user, exists := s.users[req.Id]; exists {
        return user, nil
    }
    return nil, status.Errorf(codes.NotFound, "使用者不存在")
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal(err)
    }
    
    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &userServiceServer{
        users: make(map[int32]*pb.User),
    })
    
    s.Serve(lis)
}
```

## 最佳實務總結

### API 設計原則
- **一致性:** 保持 URL 格式和命名規則一致
- **版本控制:** 透過 URL 或 Header 管理 API 版本
- **錯誤處理:** 提供清晰的錯誤訊息和狀態碼
- **安全性:** 實作適當的認證和授權機制

### 效能優化
- **快取策略:** 合理使用 HTTP 快取 Header
- **分頁處理:** 大量資料分頁回傳
- **壓縮傳輸:** 啟用 gzip 壓縮
- **連線池:** 複用資料庫連線

### 監控與測試
- **日誌記錄:** 記錄重要的 API 呼叫和錯誤
- **單元測試:** 為每個端點撰寫測試
- **整合測試:** 測試完整的 API 流程
- **效能測試:** 評估 API 在高負載下的表現

透過本章的學習，你現在具備了使用 Go 開發各種類型 API 的核心知識。無論是傳統的 RESTful API、現代的 GraphQL，還是高效能的 gRPC，Go 都提供了優秀的工具和框架來支援你的開發需求。接下來，建議你實際動手建構一個小型的 API 專案，將這些概念應用到實際開發中。