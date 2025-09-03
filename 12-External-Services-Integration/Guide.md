# Chapter 12: External Services Integration

在現代軟體架構中，應用程式很少單獨運作，通常需要與各種外部服務進行整合。Go 語言提供了豐富的工具和函式庫來簡化這些整合工作。本章將深入探討如何在 Go 應用程式中整合 HTTP 服務、gRPC 服務、訊息佇列系統、郵件服務以及第三方 API，幫助你建構強健且可擴展的分散式系統。

## HTTP Clients

**HTTP 客戶端** 是與基於 HTTP 的外部服務通信的基礎工具，Go 內建的 `net/http` 套件提供了強大的 HTTP 客戶端功能。

### 基本 HTTP 客戶端

```go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type HTTPClient struct {
    client  *http.Client
    baseURL string
    headers map[string]string
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
    return &HTTPClient{
        client: &http.Client{
            Timeout: timeout,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 100,
                IdleConnTimeout:     90 * time.Second,
            },
        },
        baseURL: baseURL,
        headers: make(map[string]string),
    }
}

func (c *HTTPClient) SetHeader(key, value string) {
    c.headers[key] = value
}

func (c *HTTPClient) Get(ctx context.Context, path string) (*http.Response, error) {
    return c.doRequest(ctx, "GET", path, nil)
}

func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
    return c.doRequest(ctx, "POST", path, body)
}

func (c *HTTPClient) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
    return c.doRequest(ctx, "PUT", path, body)
}

func (c *HTTPClient) Delete(ctx context.Context, path string) (*http.Response, error) {
    return c.doRequest(ctx, "DELETE", path, nil)
}

func (c *HTTPClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    var reqBody io.Reader
    
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
        reqBody = bytes.NewBuffer(jsonData)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // 設定預設 headers
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    
    // 設定自訂 headers
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }
    
    return c.client.Do(req)
}
```

### 重試機制與斷路器

```go
import (
    "math"
    "math/rand"
)

type RetryConfig struct {
    MaxRetries  int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
}

type CircuitBreakerConfig struct {
    MaxFailures     int
    ResetTimeout    time.Duration
    FailureTimeout  time.Duration
}

type ResilientHTTPClient struct {
    *HTTPClient
    retryConfig         RetryConfig
    circuitBreakerConfig CircuitBreakerConfig
    
    // 斷路器狀態
    failures      int
    lastFailTime  time.Time
    state         CircuitState
    mutex         sync.Mutex
}

type CircuitState int

const (
    Closed CircuitState = iota
    Open
    HalfOpen
)

func NewResilientHTTPClient(baseURL string, timeout time.Duration) *ResilientHTTPClient {
    return &ResilientHTTPClient{
        HTTPClient: NewHTTPClient(baseURL, timeout),
        retryConfig: RetryConfig{
            MaxRetries: 3,
            BaseDelay:  100 * time.Millisecond,
            MaxDelay:   5 * time.Second,
            Multiplier: 2.0,
        },
        circuitBreakerConfig: CircuitBreakerConfig{
            MaxFailures:    5,
            ResetTimeout:   60 * time.Second,
            FailureTimeout: 10 * time.Second,
        },
        state: Closed,
    }
}

func (c *ResilientHTTPClient) DoWithRetry(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    // 檢查斷路器狀態
    if err := c.checkCircuitBreaker(); err != nil {
        return nil, err
    }
    
    var lastErr error
    
    for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
        if attempt > 0 {
            // 計算延遲時間（指數退避 + 隨機抖動）
            delay := c.calculateDelay(attempt)
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(delay):
            }
        }
        
        resp, err := c.doRequest(ctx, method, path, body)
        
        if err == nil && c.isSuccessStatusCode(resp.StatusCode) {
            c.recordSuccess()
            return resp, nil
        }
        
        if err != nil {
            lastErr = err
        } else {
            lastErr = fmt.Errorf("HTTP error: %d", resp.StatusCode)
            resp.Body.Close()
        }
        
        // 檢查是否應該重試
        if !c.shouldRetry(resp, err) {
            break
        }
    }
    
    c.recordFailure()
    return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryConfig.MaxRetries+1, lastErr)
}

func (c *ResilientHTTPClient) calculateDelay(attempt int) time.Duration {
    delay := float64(c.retryConfig.BaseDelay) * math.Pow(c.retryConfig.Multiplier, float64(attempt-1))
    
    // 添加隨機抖動
    jitter := rand.Float64() * 0.1 * delay
    delay += jitter
    
    if delay > float64(c.retryConfig.MaxDelay) {
        delay = float64(c.retryConfig.MaxDelay)
    }
    
    return time.Duration(delay)
}

func (c *ResilientHTTPClient) shouldRetry(resp *http.Response, err error) bool {
    if err != nil {
        return true // 網路錯誤可重試
    }
    
    // HTTP 狀態碼判斷
    switch resp.StatusCode {
    case http.StatusTooManyRequests, http.StatusInternalServerError, 
         http.StatusBadGateway, http.StatusServiceUnavailable, 
         http.StatusGatewayTimeout:
        return true
    default:
        return false
    }
}

func (c *ResilientHTTPClient) checkCircuitBreaker() error {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    now := time.Now()
    
    switch c.state {
    case Open:
        if now.Sub(c.lastFailTime) > c.circuitBreakerConfig.ResetTimeout {
            c.state = HalfOpen
            return nil
        }
        return fmt.Errorf("circuit breaker is open")
        
    case HalfOpen:
        // 在半開狀態下允許一個請求通過
        return nil
        
    default: // Closed
        return nil
    }
}

func (c *ResilientHTTPClient) recordSuccess() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    c.failures = 0
    c.state = Closed
}

func (c *ResilientHTTPClient) recordFailure() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    c.failures++
    c.lastFailTime = time.Now()
    
    if c.failures >= c.circuitBreakerConfig.MaxFailures {
        c.state = Open
    }
}

func (c *ResilientHTTPClient) isSuccessStatusCode(code int) bool {
    return code >= 200 && code < 300
}
```

## gRPC Clients

**gRPC 客戶端** 提供高效能的遠端程序呼叫，適合微服務間的內部通信。

### 基本 gRPC 客戶端

```protobuf
// user_service.proto
syntax = "proto3";

package user;
option go_package = "./user";

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc UpdateUser(UpdateUserRequest) returns (User);
    rpc DeleteUser(DeleteUserRequest) returns (Empty);
    rpc ListUsers(ListUsersRequest) returns (stream User);
}

message User {
    int32 id = 1;
    string name = 2;
    string email = 3;
    int64 created_at = 4;
}

message GetUserRequest {
    int32 id = 1;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
}

message UpdateUserRequest {
    int32 id = 1;
    string name = 2;
    string email = 3;
}

message DeleteUserRequest {
    int32 id = 1;
}

message ListUsersRequest {
    int32 page = 1;
    int32 page_size = 2;
}

message Empty {}
```

```go
// gRPC 客戶端實作
type GRPCUserClient struct {
    conn   *grpc.ClientConn
    client user.UserServiceClient
}

func NewGRPCUserClient(address string) (*GRPCUserClient, error) {
    conn, err := grpc.Dial(address, 
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:    30 * time.Second,
            Timeout: 5 * time.Second,
        }),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
    }
    
    client := user.NewUserServiceClient(conn)
    
    return &GRPCUserClient{
        conn:   conn,
        client: client,
    }, nil
}

func (c *GRPCUserClient) GetUser(ctx context.Context, id int32) (*user.User, error) {
    req := &user.GetUserRequest{Id: id}
    
    return c.client.GetUser(ctx, req)
}

func (c *GRPCUserClient) CreateUser(ctx context.Context, name, email string) (*user.User, error) {
    req := &user.CreateUserRequest{
        Name:  name,
        Email: email,
    }
    
    return c.client.CreateUser(ctx, req)
}

func (c *GRPCUserClient) ListUsers(ctx context.Context, page, pageSize int32) ([]*user.User, error) {
    req := &user.ListUsersRequest{
        Page:     page,
        PageSize: pageSize,
    }
    
    stream, err := c.client.ListUsers(ctx, req)
    if err != nil {
        return nil, err
    }
    
    var users []*user.User
    for {
        user, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, nil
}

func (c *GRPCUserClient) Close() error {
    return c.conn.Close()
}
```

### gRPC 連線池

```go
type GRPCConnectionPool struct {
    address string
    pool    chan *grpc.ClientConn
    factory func() (*grpc.ClientConn, error)
    mu      sync.Mutex
    closed  bool
}

func NewGRPCConnectionPool(address string, poolSize int) *GRPCConnectionPool {
    pool := &GRPCConnectionPool{
        address: address,
        pool:    make(chan *grpc.ClientConn, poolSize),
        factory: func() (*grpc.ClientConn, error) {
            return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
        },
    }
    
    // 預先建立連線
    for i := 0; i < poolSize; i++ {
        conn, err := pool.factory()
        if err != nil {
            continue
        }
        pool.pool <- conn
    }
    
    return pool
}

func (p *GRPCConnectionPool) Get() (*grpc.ClientConn, error) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if p.closed {
        return nil, fmt.Errorf("connection pool is closed")
    }
    
    select {
    case conn := <-p.pool:
        if conn.GetState() == connectivity.Ready || conn.GetState() == connectivity.Idle {
            return conn, nil
        }
        // 連線狀態不佳，建立新連線
        fallthrough
    default:
        return p.factory()
    }
}

func (p *GRPCConnectionPool) Put(conn *grpc.ClientConn) {
    if conn == nil {
        return
    }
    
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if p.closed {
        conn.Close()
        return
    }
    
    select {
    case p.pool <- conn:
    default:
        // 池已滿，關閉連線
        conn.Close()
    }
}

func (p *GRPCConnectionPool) Close() {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if p.closed {
        return
    }
    
    p.closed = true
    close(p.pool)
    
    for conn := range p.pool {
        conn.Close()
    }
}
```

## Message Queues

**訊息佇列** 提供異步通信機制，實現系統間的解耦和可靠的訊息傳遞。

### RabbitMQ 整合

```go
import (
    "github.com/streadway/amqp"
)

type RabbitMQClient struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queues  map[string]amqp.Queue
    mu      sync.RWMutex
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }
    
    channel, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to open channel: %w", err)
    }
    
    client := &RabbitMQClient{
        conn:    conn,
        channel: channel,
        queues:  make(map[string]amqp.Queue),
    }
    
    return client, nil
}

func (c *RabbitMQClient) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    queue, err := c.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, nil)
    if err != nil {
        return fmt.Errorf("failed to declare queue: %w", err)
    }
    
    c.queues[name] = queue
    return nil
}

func (c *RabbitMQClient) Publish(queueName string, message []byte) error {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    return c.channel.Publish(
        "",        // exchange
        queueName, // routing key
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        message,
        },
    )
}

func (c *RabbitMQClient) Consume(queueName string, handler func([]byte) error) error {
    msgs, err := c.channel.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        return fmt.Errorf("failed to register consumer: %w", err)
    }
    
    go func() {
        for msg := range msgs {
            err := handler(msg.Body)
            if err != nil {
                msg.Nack(false, true) // 拒絕並重新排隊
            } else {
                msg.Ack(false) // 確認處理
            }
        }
    }()
    
    return nil
}

func (c *RabbitMQClient) Close() {
    if c.channel != nil {
        c.channel.Close()
    }
    if c.conn != nil {
        c.conn.Close()
    }
}
```

### Apache Kafka 整合

```go
import (
    "github.com/segmentio/kafka-go"
)

type KafkaClient struct {
    brokers []string
    writer  *kafka.Writer
    readers map[string]*kafka.Reader
    mu      sync.RWMutex
}

func NewKafkaClient(brokers []string) *KafkaClient {
    return &KafkaClient{
        brokers: brokers,
        readers: make(map[string]*kafka.Reader),
    }
}

func (c *KafkaClient) CreateProducer(topic string) {
    c.writer = &kafka.Writer{
        Addr:     kafka.TCP(c.brokers...),
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    }
}

func (c *KafkaClient) Produce(ctx context.Context, key string, message []byte) error {
    return c.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(key),
        Value: message,
    })
}

func (c *KafkaClient) CreateConsumer(topic, groupID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  c.brokers,
        Topic:    topic,
        GroupID:  groupID,
        MinBytes: 10e3, // 10KB
        MaxBytes: 10e6, // 10MB
    })
    
    c.readers[topic] = reader
}

func (c *KafkaClient) Consume(ctx context.Context, topic string, handler func([]byte) error) error {
    c.mu.RLock()
    reader, exists := c.readers[topic]
    c.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("consumer for topic %s not found", topic)
    }
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                msg, err := reader.ReadMessage(ctx)
                if err != nil {
                    if err == context.Canceled {
                        return
                    }
                    continue
                }
                
                if err := handler(msg.Value); err != nil {
                    // 處理錯誤（可以記錄日誌或發送到死信佇列）
                    continue
                }
            }
        }
    }()
    
    return nil
}

func (c *KafkaClient) Close() {
    if c.writer != nil {
        c.writer.Close()
    }
    
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    for _, reader := range c.readers {
        reader.Close()
    }
}
```

### NATS 整合

```go
import (
    "github.com/nats-io/nats.go"
)

type NATSClient struct {
    conn *nats.Conn
    js   nats.JetStreamContext
}

func NewNATSClient(url string) (*NATSClient, error) {
    conn, err := nats.Connect(url)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to NATS: %w", err)
    }
    
    js, err := conn.JetStream()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to create JetStream context: %w", err)
    }
    
    return &NATSClient{
        conn: conn,
        js:   js,
    }, nil
}

func (c *NATSClient) CreateStream(streamName string, subjects []string) error {
    _, err := c.js.AddStream(&nats.StreamConfig{
        Name:     streamName,
        Subjects: subjects,
    })
    
    return err
}

func (c *NATSClient) Publish(subject string, data []byte) error {
    _, err := c.js.Publish(subject, data)
    return err
}

func (c *NATSClient) Subscribe(subject, durableName string, handler func([]byte) error) (*nats.Subscription, error) {
    return c.js.Subscribe(subject, func(msg *nats.Msg) {
        err := handler(msg.Data)
        if err != nil {
            msg.Nak() // Negative acknowledgment
        } else {
            msg.Ack() // Acknowledgment
        }
    }, nats.Durable(durableName))
}

func (c *NATSClient) Close() {
    if c.conn != nil {
        c.conn.Close()
    }
}
```

## Email Services

**郵件服務** 整合讓應用程式能夠發送通知、確認信件和行銷郵件。

### SMTP 郵件發送

```go
import (
    "net/smtp"
    "mime/quotedprintable"
    "strings"
)

type EmailClient struct {
    smtpHost string
    smtpPort string
    username string
    password string
    auth     smtp.Auth
}

type Email struct {
    From        string
    To          []string
    Cc          []string
    Bcc         []string
    Subject     string
    Body        string
    IsHTML      bool
    Attachments []Attachment
}

type Attachment struct {
    Name        string
    Content     []byte
    ContentType string
}

func NewEmailClient(smtpHost, smtpPort, username, password string) *EmailClient {
    auth := smtp.PlainAuth("", username, password, smtpHost)
    
    return &EmailClient{
        smtpHost: smtpHost,
        smtpPort: smtpPort,
        username: username,
        password: password,
        auth:     auth,
    }
}

func (c *EmailClient) SendEmail(email Email) error {
    msg := c.buildMessage(email)
    
    addr := c.smtpHost + ":" + c.smtpPort
    to := append(email.To, email.Cc...)
    to = append(to, email.Bcc...)
    
    return smtp.SendMail(addr, c.auth, email.From, to, []byte(msg))
}

func (c *EmailClient) buildMessage(email Email) string {
    var msg strings.Builder
    
    msg.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
    msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ",")))
    
    if len(email.Cc) > 0 {
        msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(email.Cc, ",")))
    }
    
    msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
    msg.WriteString("MIME-Version: 1.0\r\n")
    
    if len(email.Attachments) > 0 {
        boundary := "boundary123456"
        msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
        msg.WriteString("\r\n")
        
        // 郵件內容
        msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
        if email.IsHTML {
            msg.WriteString("Content-Type: text/html; charset=utf-8\r\n")
        } else {
            msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
        }
        msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
        msg.WriteString(c.encodeQuotedPrintable(email.Body))
        msg.WriteString("\r\n")
        
        // 附件
        for _, attachment := range email.Attachments {
            msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
            msg.WriteString(fmt.Sprintf("Content-Type: %s\r\n", attachment.ContentType))
            msg.WriteString("Content-Transfer-Encoding: base64\r\n")
            msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n\r\n", attachment.Name))
            msg.WriteString(c.encodeBase64(attachment.Content))
            msg.WriteString("\r\n")
        }
        
        msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
    } else {
        if email.IsHTML {
            msg.WriteString("Content-Type: text/html; charset=utf-8\r\n")
        } else {
            msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
        }
        msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
        msg.WriteString(c.encodeQuotedPrintable(email.Body))
    }
    
    return msg.String()
}

func (c *EmailClient) encodeQuotedPrintable(s string) string {
    var buf strings.Builder
    w := quotedprintable.NewWriter(&buf)
    w.Write([]byte(s))
    w.Close()
    return buf.String()
}

func (c *EmailClient) encodeBase64(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}
```

### 郵件模板系統

```go
import (
    "html/template"
    "bytes"
)

type EmailTemplate struct {
    templates map[string]*template.Template
}

func NewEmailTemplate() *EmailTemplate {
    return &EmailTemplate{
        templates: make(map[string]*template.Template),
    }
}

func (et *EmailTemplate) LoadTemplate(name, content string) error {
    tmpl, err := template.New(name).Parse(content)
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }
    
    et.templates[name] = tmpl
    return nil
}

func (et *EmailTemplate) RenderTemplate(name string, data interface{}) (string, error) {
    tmpl, exists := et.templates[name]
    if !exists {
        return "", fmt.Errorf("template %s not found", name)
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("failed to execute template: %w", err)
    }
    
    return buf.String(), nil
}

// 郵件服務包裝
type MailService struct {
    client   *EmailClient
    template *EmailTemplate
}

func NewMailService(client *EmailClient) *MailService {
    return &MailService{
        client:   client,
        template: NewEmailTemplate(),
    }
}

func (ms *MailService) SendWelcomeEmail(to, username string) error {
    // 載入歡迎郵件模板
    welcomeTemplate := `
    <html>
    <body>
        <h1>歡迎, {{.Username}}!</h1>
        <p>感謝您註冊我們的服務。</p>
        <p>如果您有任何問題，請隨時聯繫我們。</p>
    </body>
    </html>
    `
    
    if err := ms.template.LoadTemplate("welcome", welcomeTemplate); err != nil {
        return err
    }
    
    body, err := ms.template.RenderTemplate("welcome", map[string]interface{}{
        "Username": username,
    })
    if err != nil {
        return err
    }
    
    email := Email{
        From:    "noreply@example.com",
        To:      []string{to},
        Subject: "歡迎加入我們！",
        Body:    body,
        IsHTML:  true,
    }
    
    return ms.client.SendEmail(email)
}
```

## Third-party APIs

**第三方 API** 整合讓應用程式能夠利用外部服務的功能，如支付、地圖、社群媒體等。

### 通用 API 客戶端框架

```go
type APIClient struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
    rateLimiter *RateLimiter
}

type RateLimiter struct {
    limiter *rate.Limiter
    burst   int
}

func NewRateLimiter(requestsPerSecond int, burst int) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
        burst:   burst,
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    return rl.limiter.Wait(ctx)
}

func NewAPIClient(baseURL, apiKey string, timeout time.Duration) *APIClient {
    return &APIClient{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: timeout,
        },
        rateLimiter: NewRateLimiter(10, 20), // 每秒10次請求，突發20次
    }
}

func (c *APIClient) CallAPI(ctx context.Context, method, endpoint string, payload interface{}, response interface{}) error {
    // 速率限制
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return fmt.Errorf("rate limit error: %w", err)
    }
    
    var reqBody io.Reader
    if payload != nil {
        jsonData, err := json.Marshal(payload)
        if err != nil {
            return fmt.Errorf("failed to marshal payload: %w", err)
        }
        reqBody = bytes.NewBuffer(jsonData)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, reqBody)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    
    // 設定標頭
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("User-Agent", "MyApp/1.0")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
    }
    
    if response != nil {
        if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
            return fmt.Errorf("failed to decode response: %w", err)
        }
    }
    
    return nil
}
```

### 支付服務整合範例

```go
type PaymentService struct {
    client *APIClient
}

type PaymentRequest struct {
    Amount   int    `json:"amount"`
    Currency string `json:"currency"`
    Source   string `json:"source"`
    Description string `json:"description,omitempty"`
}

type PaymentResponse struct {
    ID     string `json:"id"`
    Status string `json:"status"`
    Amount int    `json:"amount"`
}

func NewPaymentService(apiKey string) *PaymentService {
    client := NewAPIClient("https://api.stripe.com/v1", apiKey, 30*time.Second)
    return &PaymentService{client: client}
}

func (ps *PaymentService) CreateCharge(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
    var response PaymentResponse
    
    err := ps.client.CallAPI(ctx, "POST", "/charges", req, &response)
    if err != nil {
        return nil, fmt.Errorf("failed to create charge: %w", err)
    }
    
    return &response, nil
}

func (ps *PaymentService) GetCharge(ctx context.Context, chargeID string) (*PaymentResponse, error) {
    var response PaymentResponse
    
    endpoint := fmt.Sprintf("/charges/%s", chargeID)
    err := ps.client.CallAPI(ctx, "GET", endpoint, nil, &response)
    if err != nil {
        return nil, fmt.Errorf("failed to get charge: %w", err)
    }
    
    return &response, nil
}
```

### API 回應快取

```go
type CachedAPIClient struct {
    *APIClient
    cache *Cache
    ttl   time.Duration
}

func NewCachedAPIClient(client *APIClient, ttl time.Duration) *CachedAPIClient {
    return &CachedAPIClient{
        APIClient: client,
        cache:     NewCache(),
        ttl:       ttl,
    }
}

func (c *CachedAPIClient) CallAPIWithCache(ctx context.Context, method, endpoint string, payload interface{}, response interface{}) error {
    // 只快取 GET 請求
    if method != "GET" {
        return c.CallAPI(ctx, method, endpoint, payload, response)
    }
    
    cacheKey := fmt.Sprintf("%s:%s", method, endpoint)
    
    // 嘗試從快取取得
    if cached, found := c.cache.Get(cacheKey); found {
        cachedData := cached.([]byte)
        return json.Unmarshal(cachedData, response)
    }
    
    // 快取未命中，呼叫 API
    if err := c.CallAPI(ctx, method, endpoint, payload, response); err != nil {
        return err
    }
    
    // 儲存到快取
    if responseData, err := json.Marshal(response); err == nil {
        c.cache.Set(cacheKey, responseData, c.ttl)
    }
    
    return nil
}
```

透過本章的學習，你已經掌握了在 Go 應用程式中整合各種外部服務的完整技術棧。從 HTTP 客戶端的重試機制到訊息佇列的異步通信，從郵件服務的模板系統到第三方 API 的快取策略，這些技能將幫助你建構強健且可擴展的分散式系統。良好的外部服務整合不僅提升了系統的功能性，更是現代微服務架構的基礎。接下來，建議你將這些整合技術應用到實際專案中，體驗分散式系統開發的完整流程。