# Chapter 8: Authentication and Authorization

在現代 Web 應用程式中，**身份驗證** (Authentication) 和**授權** (Authorization) 是保護系統安全的關鍵機制。身份驗證確認使用者的身份，而授權則決定使用者能夠存取哪些資源。本章將探討如何在 Go 應用程式中實作各種身份驗證和授權機制，從基本的 Session 管理到現代的 JWT 和 OAuth 解決方案。

## Authentication vs Authorization

### Authentication（身份驗證）
**Authentication** 回答的問題是「你是誰？」，它是驗證使用者身份的過程：

- **目的:** 確認使用者聲稱的身份是否屬實
- **方法:** 密碼、生物特徵、多因素驗證等
- **結果:** 確認使用者的身份

### Authorization（授權）
**Authorization** 回答的問題是「你能做什麼？」，它是決定使用者權限的過程：

- **目的:** 確定已驗證使用者可以存取哪些資源
- **方法:** 角色基礎存取控制 (RBAC)、屬性基礎存取控制 (ABAC)
- **結果:** 允許或拒絕特定操作

### 實際應用流程

```go
// 1. Authentication - 驗證使用者身份
func authenticateUser(username, password string) (*User, error) {
    user, err := getUserByUsername(username)
    if err != nil {
        return nil, err
    }
    
    if !verifyPassword(password, user.HashedPassword) {
        return nil, errors.New("invalid credentials")
    }
    
    return user, nil
}

// 2. Authorization - 檢查使用者權限
func authorizeUser(user *User, resource string, action string) bool {
    for _, permission := range user.Permissions {
        if permission.Resource == resource && permission.Action == action {
            return true
        }
    }
    return false
}
```

## JWT (JSON Web Token)

**JWT** 是一種緊湊且安全的方式來在各方之間傳輸資訊。它特別適合用於分散式系統和 API 身份驗證。

### JWT 結構

JWT 由三個部分組成，以點號 (.) 分隔：

1. **Header:** 描述令牌類型和簽名演算法
2. **Payload:** 包含聲明 (claims) 的資料
3. **Signature:** 用於驗證令牌完整性的簽名

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

### Go 中的 JWT 實作

```go
package main

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Username string `json:"username"`
    UserID   int    `json:"user_id"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key")

// 生成 JWT Token
func generateJWT(username string, userID int, role string) (string, error) {
    claims := Claims{
        Username: username,
        UserID:   userID,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   username,
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// 驗證 JWT Token
func validateJWT(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrTokenInvalidClaims
}
```

### JWT 中介軟體

```go
func jwtMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := validateJWT(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // 將使用者資訊添加到請求內容中
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## OAuth 2.0 & OpenID Connect

**OAuth 2.0** 是一個授權框架，讓應用程式可以獲得對使用者帳戶的有限存取權限。**OpenID Connect** 則是建立在 OAuth 2.0 之上的身份驗證層。

### OAuth 2.0 流程

1. **Authorization Request:** 客戶端重導向使用者到授權伺服器
2. **User Authorization:** 使用者同意授權
3. **Authorization Grant:** 授權伺服器回傳授權碼
4. **Access Token Request:** 客戶端交換授權碼獲取存取權杖
5. **Access Token Response:** 授權伺服器回傳存取權杖

### Go OAuth 2.0 實作

```go
package main

import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var oauth2Config = &oauth2.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    RedirectURL:  "http://localhost:8080/callback",
    Scopes:       []string{"openid", "profile", "email"},
    Endpoint:     google.Endpoint,
}

// 初始化 OAuth 登入
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
    state := generateRandomState() // 產生隨機狀態值
    url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// 處理 OAuth 回呼
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")
    
    // 驗證狀態值
    if !validateState(state) {
        http.Error(w, "Invalid state", http.StatusBadRequest)
        return
    }
    
    // 交換授權碼獲取令牌
    token, err := oauth2Config.Exchange(r.Context(), code)
    if err != nil {
        http.Error(w, "Token exchange failed", http.StatusInternalServerError)
        return
    }
    
    // 使用存取權杖獲取使用者資訊
    client := oauth2Config.Client(r.Context(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    // ... 處理使用者資訊
}
```

## Session Management

**Session Management** 是傳統的狀態管理方式，在伺服器端儲存使用者的會話資訊。

### Session 的優缺點

**優點:**
- 伺服器完全控制會話狀態
- 可以立即撤銷會話
- 敏感資訊儲存在伺服器端

**缺點:**
- 需要伺服器端儲存空間
- 在分散式系統中需要共享 Session
- 不適合無狀態的微服務架構

### Go Session 實作

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "net/http"
    "sync"
    "time"
)

type Session struct {
    ID       string
    UserID   int
    Username string
    Created  time.Time
    LastUsed time.Time
}

type SessionManager struct {
    sessions map[string]*Session
    mutex    sync.RWMutex
}

func NewSessionManager() *SessionManager {
    sm := &SessionManager{
        sessions: make(map[string]*Session),
    }
    
    // 定期清理過期 Session
    go sm.cleanupExpiredSessions()
    return sm
}

// 創建新 Session
func (sm *SessionManager) CreateSession(userID int, username string) *Session {
    sessionID := generateSessionID()
    session := &Session{
        ID:       sessionID,
        UserID:   userID,
        Username: username,
        Created:  time.Now(),
        LastUsed: time.Now(),
    }
    
    sm.mutex.Lock()
    sm.sessions[sessionID] = session
    sm.mutex.Unlock()
    
    return session
}

// 獲取 Session
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
    sm.mutex.RLock()
    session, exists := sm.sessions[sessionID]
    sm.mutex.RUnlock()
    
    if exists {
        // 更新最後使用時間
        sm.mutex.Lock()
        session.LastUsed = time.Now()
        sm.mutex.Unlock()
    }
    
    return session, exists
}

// 刪除 Session
func (sm *SessionManager) DeleteSession(sessionID string) {
    sm.mutex.Lock()
    delete(sm.sessions, sessionID)
    sm.mutex.Unlock()
}

func generateSessionID() string {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    return base64.URLEncoding.EncodeToString(bytes)
}
```

### Session 中介軟體

```go
func sessionMiddleware(sm *SessionManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            cookie, err := r.Cookie("session_id")
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            session, exists := sm.GetSession(cookie.Value)
            if !exists {
                http.Error(w, "Invalid session", http.StatusUnauthorized)
                return
            }
            
            // 檢查 Session 是否過期
            if time.Since(session.LastUsed) > 30*time.Minute {
                sm.DeleteSession(session.ID)
                http.Error(w, "Session expired", http.StatusUnauthorized)
                return
            }
            
            ctx := context.WithValue(r.Context(), "session", session)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## 角色基礎存取控制 (RBAC)

**RBAC** 是一種根據使用者角色來控制存取權限的機制，廣泛應用於企業級系統中。

### RBAC 模型實作

```go
type Permission struct {
    Resource string
    Action   string
}

type Role struct {
    Name        string
    Permissions []Permission
}

type User struct {
    ID       int
    Username string
    Roles    []Role
}

// 檢查使用者是否具有特定權限
func (u *User) HasPermission(resource, action string) bool {
    for _, role := range u.Roles {
        for _, permission := range role.Permissions {
            if permission.Resource == resource && permission.Action == action {
                return true
            }
        }
    }
    return false
}

// 權限檢查中介軟體
func requirePermission(resource, action string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := r.Context().Value("user").(*User)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            if !user.HasPermission(resource, action) {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## 安全最佳實務

### 密碼安全

```go
import "golang.org/x/crypto/bcrypt"

// 雜湊密碼
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// 驗證密碼
func verifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 防止常見攻擊

- **SQL Injection:** 使用參數化查詢
- **XSS:** 適當的輸出編碼和 CSP Header
- **CSRF:** 使用 CSRF Token
- **Brute Force:** 實作登入嘗試限制

```go
// CSRF 保護中介軟體
func csrfProtection(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            token := r.Header.Get("X-CSRF-Token")
            if !validateCSRFToken(token, r) {
                http.Error(w, "CSRF token invalid", http.StatusForbidden)
                return
            }
        }
        next.ServeHTTP(w, r)
    })
}
```

## 選擇適當的驗證方式

### Session vs JWT 比較

| 特性 | Session | JWT |
|------|---------|-----|
| 儲存位置 | 伺服器端 | 客戶端 |
| 可撤銷性 | 立即 | 困難 |
| 擴展性 | 需要共享儲存 | 無狀態 |
| 安全性 | 較高 | 需要小心處理 |
| 適用場景 | 傳統 Web 應用 | API 和微服務 |

### 建議選擇策略

- **內部系統:** 使用 Session + Cookie
- **公開 API:** 使用 JWT
- **第三方整合:** 使用 OAuth 2.0
- **企業應用:** 結合 RBAC 和 LDAP/AD

透過本章的學習，你已經掌握了在 Go 中實作身份驗證和授權的核心技術。無論是傳統的 Session 管理、現代的 JWT 方案，還是標準的 OAuth 2.0 流程，每種方式都有其適用的場景。選擇合適的驗證機制，並遵循安全最佳實務，是建構安全可靠應用程式的基礎。接下來，建議你實作一個包含多種驗證方式的範例專案，深入體驗這些概念在實際開發中的應用。