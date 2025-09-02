# Chapter 11: Logging and Monitoring

日誌記錄和監控是現代應用程式運維的核心組成部分，對於系統的可觀測性、問題診斷和效能調優至關重要。Go 語言提供了豐富的日誌和監控工具，幫助開發者建構可靠、易於維護的系統。本章將深入探討日誌記錄最佳實務、結構化日誌、監控指標收集以及錯誤追蹤等關鍵技術，讓你的 Go 應用程式具備完善的可觀測性。

## Logging Best Practices

**日誌記錄** 是應用程式運行狀態的記錄和追蹤機制，良好的日誌實務能夠大幅提升系統的可維護性。

### 日誌層級管理

```go
package main

import (
    "log/slog"
    "os"
)

// 定義日誌層級
const (
    LevelTrace = slog.Level(-8)
    LevelDebug = slog.LevelDebug
    LevelInfo  = slog.LevelInfo
    LevelWarn  = slog.LevelWarn
    LevelError = slog.LevelError
    LevelFatal = slog.Level(12)
)

type Logger struct {
    *slog.Logger
}

func NewLogger(level slog.Level) *Logger {
    opts := &slog.HandlerOptions{
        Level: level,
        AddSource: true,
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    logger := slog.New(handler)
    
    return &Logger{Logger: logger}
}

func (l *Logger) Trace(msg string, args ...any) {
    l.Log(nil, LevelTrace, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...any) {
    l.Log(nil, LevelFatal, msg, args...)
    os.Exit(1)
}

// 使用範例
func main() {
    logger := NewLogger(slog.LevelInfo)
    
    logger.Info("Application starting", "version", "1.0.0", "port", 8080)
    logger.Warn("High memory usage detected", "usage", "85%")
    logger.Error("Database connection failed", "error", "connection timeout")
}
```

### 內容相關的日誌記錄

```go
import (
    "context"
    "log/slog"
)

// 為 context 添加日誌欄位
type contextKey string

const (
    requestIDKey contextKey = "request_id"
    userIDKey    contextKey = "user_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, requestIDKey, requestID)
}

func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

// 從 context 提取日誌欄位
func LoggerFromContext(ctx context.Context, logger *slog.Logger) *slog.Logger {
    if requestID, ok := ctx.Value(requestIDKey).(string); ok {
        logger = logger.With("request_id", requestID)
    }
    
    if userID, ok := ctx.Value(userIDKey).(string); ok {
        logger = logger.With("user_id", userID)
    }
    
    return logger
}

// HTTP 中介軟體範例
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestID := generateRequestID()
            ctx := WithRequestID(r.Context(), requestID)
            r = r.WithContext(ctx)
            
            contextLogger := LoggerFromContext(ctx, logger)
            contextLogger.Info("Request started",
                "method", r.Method,
                "path", r.URL.Path,
                "remote_addr", r.RemoteAddr,
            )
            
            start := time.Now()
            next.ServeHTTP(w, r)
            duration := time.Since(start)
            
            contextLogger.Info("Request completed",
                "duration", duration,
                "status", getStatusCode(w),
            )
        })
    }
}
```

### 日誌輪轉和歸檔

```go
import (
    "io"
    "log/slog"
    "gopkg.in/natefinch/lumberjack.v2"
)

func setupRotatingLogger() *slog.Logger {
    // 設定日誌輪轉
    rotator := &lumberjack.Logger{
        Filename:   "/var/log/myapp/app.log",
        MaxSize:    100, // MB
        MaxBackups: 3,
        MaxAge:     28, // days
        Compress:   true,
    }
    
    // 同時輸出到檔案和標準輸出
    multiWriter := io.MultiWriter(os.Stdout, rotator)
    
    opts := &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: true,
    }
    
    handler := slog.NewJSONHandler(multiWriter, opts)
    return slog.New(handler)
}
```

## Structured Logging

**結構化日誌** 使用結構化格式（如 JSON）記錄日誌，便於自動化解析和分析。

### JSON 結構化日誌

```go
type StructuredLogger struct {
    *slog.Logger
    serviceName string
    version     string
}

func NewStructuredLogger(serviceName, version string) *StructuredLogger {
    opts := &slog.HandlerOptions{
        Level: slog.LevelInfo,
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            // 自訂時間戳格式
            if a.Key == slog.TimeKey {
                return slog.String("timestamp", a.Value.Time().Format(time.RFC3339))
            }
            return a
        },
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    baseLogger := slog.New(handler).With(
        "service", serviceName,
        "version", version,
        "hostname", getHostname(),
    )
    
    return &StructuredLogger{
        Logger:      baseLogger,
        serviceName: serviceName,
        version:     version,
    }
}

// 業務日誌方法
func (l *StructuredLogger) LogUserAction(ctx context.Context, action, userID string, metadata map[string]interface{}) {
    attrs := []slog.Attr{
        slog.String("event_type", "user_action"),
        slog.String("action", action),
        slog.String("user_id", userID),
    }
    
    // 添加 metadata
    for key, value := range metadata {
        attrs = append(attrs, slog.Any(key, value))
    }
    
    l.LogAttrs(ctx, slog.LevelInfo, "User action performed", attrs...)
}

func (l *StructuredLogger) LogAPICall(ctx context.Context, method, endpoint string, duration time.Duration, statusCode int) {
    l.Info("API call completed",
        "event_type", "api_call",
        "method", method,
        "endpoint", endpoint,
        "duration_ms", duration.Milliseconds(),
        "status_code", statusCode,
    )
}

func (l *StructuredLogger) LogError(ctx context.Context, err error, component string, metadata map[string]interface{}) {
    attrs := []slog.Attr{
        slog.String("event_type", "error"),
        slog.String("error", err.Error()),
        slog.String("component", component),
    }
    
    // 添加錯誤堆疊
    if stackTracer, ok := err.(interface{ StackTrace() string }); ok {
        attrs = append(attrs, slog.String("stack_trace", stackTracer.StackTrace()))
    }
    
    for key, value := range metadata {
        attrs = append(attrs, slog.Any(key, value))
    }
    
    l.LogAttrs(ctx, slog.LevelError, "Error occurred", attrs...)
}
```

### 自訂日誌格式化器

```go
type CustomHandler struct {
    slog.Handler
    serviceName string
}

func NewCustomHandler(serviceName string) *CustomHandler {
    opts := &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }
    
    return &CustomHandler{
        Handler:     slog.NewJSONHandler(os.Stdout, opts),
        serviceName: serviceName,
    }
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
    // 添加自訂欄位
    r.AddAttrs(
        slog.String("service", h.serviceName),
        slog.String("trace_id", getTraceIDFromContext(ctx)),
        slog.String("environment", os.Getenv("ENV")),
    )
    
    return h.Handler.Handle(ctx, r)
}

// 使用自訂處理器
func setupCustomLogger() *slog.Logger {
    handler := NewCustomHandler("my-service")
    return slog.New(handler)
}
```

## Log Aggregation

**日誌聚合** 將分散在不同服務和節點的日誌統一收集、儲存和分析。

### 日誌發送器

```go
type LogShipper struct {
    endpoint string
    client   *http.Client
    buffer   chan LogEntry
    quit     chan struct{}
}

type LogEntry struct {
    Timestamp time.Time              `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields"`
    Service   string                 `json:"service"`
}

func NewLogShipper(endpoint string, bufferSize int) *LogShipper {
    shipper := &LogShipper{
        endpoint: endpoint,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        buffer: make(chan LogEntry, bufferSize),
        quit:   make(chan struct{}),
    }
    
    go shipper.ship()
    return shipper
}

func (ls *LogShipper) ship() {
    batch := make([]LogEntry, 0, 100)
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case entry := <-ls.buffer:
            batch = append(batch, entry)
            
            // 當批次達到上限時立即發送
            if len(batch) >= 100 {
                ls.sendBatch(batch)
                batch = batch[:0]
            }
            
        case <-ticker.C:
            // 定期發送批次（即使未滿）
            if len(batch) > 0 {
                ls.sendBatch(batch)
                batch = batch[:0]
            }
            
        case <-ls.quit:
            // 關閉前發送剩餘的日誌
            if len(batch) > 0 {
                ls.sendBatch(batch)
            }
            return
        }
    }
}

func (ls *LogShipper) sendBatch(entries []LogEntry) {
    data, err := json.Marshal(entries)
    if err != nil {
        fmt.Printf("Failed to marshal log entries: %v\n", err)
        return
    }
    
    resp, err := ls.client.Post(ls.endpoint, "application/json", bytes.NewReader(data))
    if err != nil {
        fmt.Printf("Failed to ship logs: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Log shipping failed with status: %d\n", resp.StatusCode)
    }
}

func (ls *LogShipper) Ship(entry LogEntry) {
    select {
    case ls.buffer <- entry:
    default:
        // 緩衝區滿時丟棄日誌
        fmt.Println("Log buffer full, dropping entry")
    }
}

func (ls *LogShipper) Close() {
    close(ls.quit)
}
```

### ELK Stack 整合

```go
import (
    "github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchLogger struct {
    client *elasticsearch.Client
    index  string
}

func NewElasticsearchLogger(addresses []string, index string) (*ElasticsearchLogger, error) {
    cfg := elasticsearch.Config{
        Addresses: addresses,
    }
    
    client, err := elasticsearch.NewClient(cfg)
    if err != nil {
        return nil, err
    }
    
    return &ElasticsearchLogger{
        client: client,
        index:  index,
    }, nil
}

func (el *ElasticsearchLogger) IndexLog(ctx context.Context, logEntry LogEntry) error {
    data, err := json.Marshal(logEntry)
    if err != nil {
        return err
    }
    
    req := esapi.IndexRequest{
        Index:   el.index,
        Body:    bytes.NewReader(data),
        Refresh: "true",
    }
    
    res, err := req.Do(ctx, el.client)
    if err != nil {
        return err
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return fmt.Errorf("indexing failed: %s", res.Status())
    }
    
    return nil
}
```

## Monitoring & Metrics

**監控和指標** 收集系統運行時的各種指標數據，用於效能分析和問題預警。

### Prometheus 指標收集

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
    requestsTotal    *prometheus.CounterVec
    requestDuration  *prometheus.HistogramVec
    activeConnections prometheus.Gauge
    memoryUsage      prometheus.GaugeFunc
}

func NewMetrics() *Metrics {
    metrics := &Metrics{
        requestsTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "endpoint", "status"},
        ),
        
        requestDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "http_request_duration_seconds",
                Help:    "Duration of HTTP requests",
                Buckets: prometheus.DefBuckets,
            },
            []string{"method", "endpoint"},
        ),
        
        activeConnections: promauto.NewGauge(
            prometheus.GaugeOpts{
                Name: "active_connections",
                Help: "Number of active connections",
            },
        ),
    }
    
    // 記憶體使用指標
    metrics.memoryUsage = promauto.NewGaugeFunc(
        prometheus.GaugeOpts{
            Name: "memory_usage_bytes",
            Help: "Current memory usage in bytes",
        },
        func() float64 {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            return float64(m.Alloc)
        },
    )
    
    return metrics
}

// HTTP 監控中介軟體
func (m *Metrics) HTTPMetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // 包裝 ResponseWriter 以捕獲狀態碼
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start)
        method := r.Method
        endpoint := r.URL.Path
        status := fmt.Sprintf("%d", wrapped.statusCode)
        
        // 記錄指標
        m.requestsTotal.WithLabelValues(method, endpoint, status).Inc()
        m.requestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### 自訂業務指標

```go
type BusinessMetrics struct {
    userRegistrations *prometheus.CounterVec
    orderValues       *prometheus.HistogramVec
    cacheHitRate      *prometheus.GaugeVec
}

func NewBusinessMetrics() *BusinessMetrics {
    return &BusinessMetrics{
        userRegistrations: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "user_registrations_total",
                Help: "Total number of user registrations",
            },
            []string{"source", "country"},
        ),
        
        orderValues: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "order_value_dollars",
                Help:    "Distribution of order values",
                Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
            },
            []string{"category"},
        ),
        
        cacheHitRate: promauto.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "cache_hit_rate",
                Help: "Cache hit rate by cache name",
            },
            []string{"cache_name"},
        ),
    }
}

func (bm *BusinessMetrics) RecordUserRegistration(source, country string) {
    bm.userRegistrations.WithLabelValues(source, country).Inc()
}

func (bm *BusinessMetrics) RecordOrderValue(category string, value float64) {
    bm.orderValues.WithLabelValues(category).Observe(value)
}

func (bm *BusinessMetrics) UpdateCacheHitRate(cacheName string, hitRate float64) {
    bm.cacheHitRate.WithLabelValues(cacheName).Set(hitRate)
}
```

## Health Checks

**健康檢查** 提供系統健康狀態的端點，用於負載平衡器和監控系統。

### 基本健康檢查

```go
type HealthChecker struct {
    checks map[string]HealthCheck
    mu     sync.RWMutex
}

type HealthCheck interface {
    Name() string
    Check(ctx context.Context) error
}

type HealthStatus struct {
    Status string                    `json:"status"`
    Checks map[string]CheckResult    `json:"checks"`
    Uptime string                    `json:"uptime"`
    Version string                   `json:"version"`
}

type CheckResult struct {
    Status  string        `json:"status"`
    Message string        `json:"message,omitempty"`
    Latency time.Duration `json:"latency"`
}

func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        checks: make(map[string]HealthCheck),
    }
}

func (hc *HealthChecker) AddCheck(check HealthCheck) {
    hc.mu.Lock()
    defer hc.mu.Unlock()
    hc.checks[check.Name()] = check
}

func (hc *HealthChecker) CheckHealth(ctx context.Context) HealthStatus {
    hc.mu.RLock()
    defer hc.mu.RUnlock()
    
    results := make(map[string]CheckResult)
    overallStatus := "healthy"
    
    for name, check := range hc.checks {
        start := time.Now()
        err := check.Check(ctx)
        latency := time.Since(start)
        
        result := CheckResult{
            Status:  "healthy",
            Latency: latency,
        }
        
        if err != nil {
            result.Status = "unhealthy"
            result.Message = err.Error()
            overallStatus = "unhealthy"
        }
        
        results[name] = result
    }
    
    return HealthStatus{
        Status:  overallStatus,
        Checks:  results,
        Uptime:  time.Since(startTime).String(),
        Version: version,
    }
}

// HTTP 處理器
func (hc *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    
    status := hc.CheckHealth(ctx)
    
    w.Header().Set("Content-Type", "application/json")
    if status.Status == "unhealthy" {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    
    json.NewEncoder(w).Encode(status)
}
```

### 具體健康檢查實作

```go
// 資料庫健康檢查
type DatabaseHealthCheck struct {
    db *sql.DB
}

func NewDatabaseHealthCheck(db *sql.DB) *DatabaseHealthCheck {
    return &DatabaseHealthCheck{db: db}
}

func (dhc *DatabaseHealthCheck) Name() string {
    return "database"
}

func (dhc *DatabaseHealthCheck) Check(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    return dhc.db.PingContext(ctx)
}

// Redis 健康檢查
type RedisHealthCheck struct {
    client *redis.Client
}

func NewRedisHealthCheck(client *redis.Client) *RedisHealthCheck {
    return &RedisHealthCheck{client: client}
}

func (rhc *RedisHealthCheck) Name() string {
    return "redis"
}

func (rhc *RedisHealthCheck) Check(ctx context.Context) error {
    return rhc.client.Ping(ctx).Err()
}

// 外部服務健康檢查
type HTTPServiceHealthCheck struct {
    name string
    url  string
    client *http.Client
}

func NewHTTPServiceHealthCheck(name, url string) *HTTPServiceHealthCheck {
    return &HTTPServiceHealthCheck{
        name: name,
        url:  url,
        client: &http.Client{Timeout: 5 * time.Second},
    }
}

func (hshc *HTTPServiceHealthCheck) Name() string {
    return hshc.name
}

func (hshc *HTTPServiceHealthCheck) Check(ctx context.Context) error {
    req, err := http.NewRequestWithContext(ctx, "GET", hshc.url, nil)
    if err != nil {
        return err
    }
    
    resp, err := hshc.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        return fmt.Errorf("service returned status %d", resp.StatusCode)
    }
    
    return nil
}
```

## Error Tracking

**錯誤追蹤** 自動收集、分析和報告應用程式中的錯誤和異常。

### 錯誤收集器

```go
type ErrorTracker struct {
    projectID string
    client    *http.Client
    buffer    chan ErrorEvent
    quit      chan struct{}
}

type ErrorEvent struct {
    ID        string                 `json:"id"`
    Timestamp time.Time              `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Exception ExceptionInfo          `json:"exception"`
    Context   map[string]interface{} `json:"context"`
    User      UserInfo               `json:"user"`
    Tags      map[string]string      `json:"tags"`
}

type ExceptionInfo struct {
    Type       string `json:"type"`
    Value      string `json:"value"`
    Stacktrace string `json:"stacktrace"`
}

type UserInfo struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func NewErrorTracker(projectID string) *ErrorTracker {
    tracker := &ErrorTracker{
        projectID: projectID,
        client:    &http.Client{Timeout: 10 * time.Second},
        buffer:    make(chan ErrorEvent, 1000),
        quit:      make(chan struct{}),
    }
    
    go tracker.processEvents()
    return tracker
}

func (et *ErrorTracker) CaptureError(err error, ctx context.Context, tags map[string]string) {
    event := ErrorEvent{
        ID:        generateUUID(),
        Timestamp: time.Now(),
        Level:     "error",
        Message:   err.Error(),
        Exception: ExceptionInfo{
            Type:       fmt.Sprintf("%T", err),
            Value:      err.Error(),
            Stacktrace: getStackTrace(),
        },
        Context: extractContextInfo(ctx),
        User:    extractUserInfo(ctx),
        Tags:    tags,
    }
    
    select {
    case et.buffer <- event:
    default:
        // 緩衝區滿時記錄到本地日誌
        slog.Warn("Error tracker buffer full, dropping event")
    }
}

func (et *ErrorTracker) processEvents() {
    batch := make([]ErrorEvent, 0, 50)
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case event := <-et.buffer:
            batch = append(batch, event)
            if len(batch) >= 50 {
                et.sendBatch(batch)
                batch = batch[:0]
            }
            
        case <-ticker.C:
            if len(batch) > 0 {
                et.sendBatch(batch)
                batch = batch[:0]
            }
            
        case <-et.quit:
            if len(batch) > 0 {
                et.sendBatch(batch)
            }
            return
        }
    }
}

func (et *ErrorTracker) sendBatch(events []ErrorEvent) {
    data, err := json.Marshal(events)
    if err != nil {
        slog.Error("Failed to marshal error events", "error", err)
        return
    }
    
    endpoint := fmt.Sprintf("https://api.errortracker.com/projects/%s/events", et.projectID)
    resp, err := et.client.Post(endpoint, "application/json", bytes.NewReader(data))
    if err != nil {
        slog.Error("Failed to send error events", "error", err)
        return
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        slog.Error("Error tracking service returned error", "status", resp.StatusCode)
    }
}
```

### 錯誤處理中介軟體

```go
func ErrorTrackingMiddleware(tracker *ErrorTracker) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if recovered := recover(); recovered != nil {
                    err := fmt.Errorf("panic recovered: %v", recovered)
                    
                    tags := map[string]string{
                        "method": r.Method,
                        "path":   r.URL.Path,
                        "type":   "panic",
                    }
                    
                    tracker.CaptureError(err, r.Context(), tags)
                    
                    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                }
            }()
            
            next.ServeHTTP(w, r)
        })
    }
}

// 錯誤處理包裝器
func HandleError(tracker *ErrorTracker, handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := handler(w, r); err != nil {
            tags := map[string]string{
                "method": r.Method,
                "path":   r.URL.Path,
                "type":   "handler_error",
            }
            
            tracker.CaptureError(err, r.Context(), tags)
            
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    }
}
```

### 整合監控儀表板

```go
type MonitoringDashboard struct {
    logger       *slog.Logger
    metrics      *Metrics
    healthChecker *HealthChecker
    errorTracker *ErrorTracker
}

func NewMonitoringDashboard(logger *slog.Logger) *MonitoringDashboard {
    return &MonitoringDashboard{
        logger:       logger,
        metrics:      NewMetrics(),
        healthChecker: NewHealthChecker(),
        errorTracker: NewErrorTracker("your-project-id"),
    }
}

func (md *MonitoringDashboard) SetupRoutes(mux *http.ServeMux) {
    // Prometheus 指標端點
    mux.Handle("/metrics", promhttp.Handler())
    
    // 健康檢查端點
    mux.HandleFunc("/health", md.healthChecker.HealthHandler)
    
    // 詳細的系統資訊端點
    mux.HandleFunc("/debug/info", md.systemInfoHandler)
}

func (md *MonitoringDashboard) systemInfoHandler(w http.ResponseWriter, r *http.Request) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    info := map[string]interface{}{
        "go_version":    runtime.Version(),
        "goroutines":    runtime.NumGoroutine(),
        "memory": map[string]interface{}{
            "alloc":     m.Alloc,
            "sys":       m.Sys,
            "heap_objects": m.HeapObjects,
        },
        "gc": map[string]interface{}{
            "num_gc":     m.NumGC,
            "last_gc":    time.Unix(0, int64(m.LastGC)),
            "pause_total": time.Duration(m.PauseTotalNs),
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(info)
}
```

透過本章的學習，你已經掌握了建構完整可觀測性系統的核心技術。從結構化日誌記錄到指標收集，從健康檢查到錯誤追蹤，這些工具和技術將幫助你建構易於監控、診斷和維護的 Go 應用程式。良好的日誌和監控體系不僅能幫助快速定位問題，更是系統可靠性和效能優化的基礎。接下來，建議你將這些監控技術整合到實際專案中，建立全面的可觀測性平台。