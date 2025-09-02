# Chapter 10: Performance Optimization

效能優化是軟體開發中至關重要的環節，特別是在高負載和大規模系統中。Go 語言雖然本身已經具備優異的效能特性，但透過適當的優化技術，我們可以進一步提升應用程式的執行效率。本章將深入探討 Go 應用程式的效能分析、記憶體管理、並發優化以及快取策略等關鍵技術，幫助你建構高效能的 Go 應用程式。

## Profiling (pprof)

**pprof** 是 Go 內建的效能分析工具，能夠幫助開發者識別程式碼中的效能瓶頸和資源使用情況。

### pprof 基本概念

pprof 可以分析多種效能指標：
- **CPU Profiling:** 分析 CPU 使用情況
- **Memory Profiling:** 分析記憶體分配和使用
- **Goroutine Profiling:** 分析 Goroutine 的建立和狀態
- **Block Profiling:** 分析阻塞操作
- **Mutex Profiling:** 分析互斥鎖競爭

### 在程式中啟用 pprof

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    _ "net/http/pprof" // 匯入 pprof
    "runtime"
    "time"
)

func main() {
    // 啟動 pprof HTTP 伺服器
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 模擬一些工作負載
    for i := 0; i < 1000000; i++ {
        go doWork(i)
    }
    
    // 保持程式運行
    select {}
}

func doWork(id int) {
    // 模擬 CPU 密集型工作
    sum := 0
    for i := 0; i < 1000000; i++ {
        sum += i
    }
    
    // 模擬記憶體分配
    data := make([]byte, 1024)
    _ = data
    
    time.Sleep(100 * time.Millisecond)
}
```

### 使用 pprof 分析效能

```bash
# 分析 CPU 效能（收集 30 秒資料）
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 分析記憶體使用
go tool pprof http://localhost:6060/debug/pprof/heap

# 分析 Goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 分析阻塞情況
go tool pprof http://localhost:6060/debug/pprof/block
```

### pprof 交互式分析

```bash
# 進入 pprof 交互模式
(pprof) top
Showing nodes accounting for 10.48s, 95.55% of 10.97s total
Dropped 63 nodes (cum < 0.05s)
      flat  flat%   sum%        cum   cum%
     3.64s 33.18% 33.18%      3.64s 33.18%  main.doWork
     2.11s 19.24% 52.42%      2.11s 19.24%  runtime.mallocgc

# 查看函式詳細資訊
(pprof) list main.doWork

# 查看呼叫圖
(pprof) web

# 查看火焰圖
(pprof) png
```

### 程式碼中的 CPU Profiling

```go
package main

import (
    "os"
    "runtime/pprof"
    "log"
)

func main() {
    // 創建 CPU profile 檔案
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    // 啟動 CPU profiling
    if err := pprof.StartCPUProfile(f); err != nil {
        log.Fatal(err)
    }
    defer pprof.StopCPUProfile()
    
    // 執行要分析的程式碼
    performanceCriticalFunction()
}

func performanceCriticalFunction() {
    // 模擬 CPU 密集型任務
    for i := 0; i < 100000000; i++ {
        _ = i * i
    }
}
```

## Memory Management

Go 的垃圾回收器 (GC) 自動管理記憶體，但理解記憶體分配和回收機制對效能優化至關重要。

### 記憶體分配模式

```go
// 避免不必要的記憶體分配
func inefficientString() string {
    var result string
    for i := 0; i < 1000; i++ {
        result += "hello" // 每次連接都會分配新的記憶體
    }
    return result
}

func efficientString() string {
    var builder strings.Builder
    builder.Grow(5000) // 預先分配足夠容量
    for i := 0; i < 1000; i++ {
        builder.WriteString("hello")
    }
    return builder.String()
}
```

### Slice 優化

```go
// 低效的 slice 使用
func inefficientSlice() []int {
    var slice []int
    for i := 0; i < 10000; i++ {
        slice = append(slice, i) // 可能觸發多次重新分配
    }
    return slice
}

// 高效的 slice 使用
func efficientSlice() []int {
    slice := make([]int, 0, 10000) // 預先分配容量
    for i := 0; i < 10000; i++ {
        slice = append(slice, i)
    }
    return slice
}

// 正確重設 slice
func resetSlice(s []int) []int {
    return s[:0] // 重複使用底層陣列
}
```

### 物件池模式

```go
import "sync"

type Buffer struct {
    data []byte
}

var bufferPool = sync.Pool{
    New: func() interface{} {
        return &Buffer{
            data: make([]byte, 0, 1024),
        }
    },
}

func processData(data []byte) []byte {
    // 從池中獲取 buffer
    buf := bufferPool.Get().(*Buffer)
    defer bufferPool.Put(buf)
    
    // 重設 buffer
    buf.data = buf.data[:0]
    
    // 處理資料
    buf.data = append(buf.data, data...)
    // ... 其他處理邏輯
    
    // 複製結果（因為 buffer 會被回收）
    result := make([]byte, len(buf.data))
    copy(result, buf.data)
    
    return result
}
```

### 記憶體洩漏偵測

```go
func detectMemoryLeak() {
    var m runtime.MemStats
    
    // 記錄初始記憶體使用量
    runtime.GC()
    runtime.ReadMemStats(&m)
    initialAlloc := m.Alloc
    
    // 執行可能造成記憶體洩漏的操作
    for i := 0; i < 100000; i++ {
        data := make([]byte, 1024)
        _ = data
    }
    
    // 強制垃圾回收
    runtime.GC()
    runtime.ReadMemStats(&m)
    finalAlloc := m.Alloc
    
    fmt.Printf("Memory leak: %d bytes\n", finalAlloc-initialAlloc)
}
```

## Goroutine Optimization

Goroutine 是 Go 並發程式設計的核心，但不當使用可能導致效能問題。

### Goroutine Pool

```go
type WorkerPool struct {
    workerChan chan func()
    quit       chan bool
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
    pool := &WorkerPool{
        workerChan: make(chan func(), maxWorkers),
        quit:       make(chan bool),
    }
    
    // 啟動固定數量的 worker
    for i := 0; i < maxWorkers; i++ {
        go pool.worker()
    }
    
    return pool
}

func (p *WorkerPool) worker() {
    for {
        select {
        case work := <-p.workerChan:
            work()
        case <-p.quit:
            return
        }
    }
}

func (p *WorkerPool) Submit(work func()) {
    p.workerChan <- work
}

func (p *WorkerPool) Stop() {
    close(p.quit)
}

// 使用範例
func useWorkerPool() {
    pool := NewWorkerPool(10)
    defer pool.Stop()
    
    for i := 0; i < 1000; i++ {
        taskID := i
        pool.Submit(func() {
            fmt.Printf("Processing task %d\n", taskID)
            time.Sleep(100 * time.Millisecond)
        })
    }
}
```

### 控制 Goroutine 數量

```go
// 使用 channel 控制並發數量
func controlledConcurrency(tasks []Task, maxWorkers int) {
    semaphore := make(chan struct{}, maxWorkers)
    var wg sync.WaitGroup
    
    for _, task := range tasks {
        wg.Add(1)
        go func(t Task) {
            defer wg.Done()
            
            // 獲取信號量
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            // 執行任務
            processTask(t)
        }(task)
    }
    
    wg.Wait()
}
```

### Context 的效能影響

```go
// 避免過度使用 context.WithTimeout
func efficientContext(ctx context.Context) {
    // 共用 timeout context
    timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            doWorkWithContext(timeoutCtx, id)
        }(i)
    }
    
    wg.Wait()
}
```

## Caching Strategies

快取是提升應用程式效能的重要策略，可以顯著減少資料存取延遲。

### In-Memory Cache

```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]CacheItem
}

type CacheItem struct {
    value      interface{}
    expiration time.Time
}

func NewCache() *Cache {
    cache := &Cache{
        data: make(map[string]CacheItem),
    }
    
    // 定期清理過期項目
    go cache.cleanup()
    return cache
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.data[key] = CacheItem{
        value:      value,
        expiration: time.Now().Add(ttl),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, exists := c.data[key]
    if !exists {
        return nil, false
    }
    
    if time.Now().After(item.expiration) {
        return nil, false
    }
    
    return item.value, true
}

func (c *Cache) cleanup() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, item := range c.data {
            if now.After(item.expiration) {
                delete(c.data, key)
            }
        }
        c.mu.Unlock()
    }
}
```

### LRU Cache

```go
type LRUCache struct {
    mu       sync.Mutex
    capacity int
    cache    map[string]*list.Element
    list     *list.List
}

type entry struct {
    key   string
    value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[string]*list.Element),
        list:     list.New(),
    }
}

func (lru *LRUCache) Get(key string) (interface{}, bool) {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if elem, exists := lru.cache[key]; exists {
        lru.list.MoveToFront(elem)
        return elem.Value.(*entry).value, true
    }
    
    return nil, false
}

func (lru *LRUCache) Put(key string, value interface{}) {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if elem, exists := lru.cache[key]; exists {
        elem.Value.(*entry).value = value
        lru.list.MoveToFront(elem)
        return
    }
    
    elem := lru.list.PushFront(&entry{key: key, value: value})
    lru.cache[key] = elem
    
    if lru.list.Len() > lru.capacity {
        oldest := lru.list.Back()
        lru.list.Remove(oldest)
        delete(lru.cache, oldest.Value.(*entry).key)
    }
}
```

### HTTP Response Cache

```go
type HTTPCache struct {
    cache *Cache
}

func NewHTTPCache() *HTTPCache {
    return &HTTPCache{
        cache: NewCache(),
    }
}

func (hc *HTTPCache) CacheMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 只快取 GET 請求
        if r.Method != http.MethodGet {
            next.ServeHTTP(w, r)
            return
        }
        
        cacheKey := r.URL.String()
        
        // 嘗試從快取取得回應
        if cachedResponse, found := hc.cache.Get(cacheKey); found {
            response := cachedResponse.([]byte)
            w.Header().Set("X-Cache", "HIT")
            w.Write(response)
            return
        }
        
        // 快取未命中，執行原始處理器
        recorder := &responseRecorder{
            ResponseWriter: w,
            body:          &bytes.Buffer{},
        }
        
        next.ServeHTTP(recorder, r)
        
        // 快取回應
        hc.cache.Set(cacheKey, recorder.body.Bytes(), 10*time.Minute)
        w.Header().Set("X-Cache", "MISS")
    })
}

type responseRecorder struct {
    http.ResponseWriter
    body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
    r.body.Write(b)
    return r.ResponseWriter.Write(b)
}
```

## Load Testing

負載測試幫助驗證系統在高負載下的效能表現。

### 基本負載測試

```go
func loadTest() {
    const (
        numRequests    = 10000
        concurrency    = 100
        targetURL      = "http://localhost:8080/api/test"
    )
    
    var wg sync.WaitGroup
    requestChan := make(chan struct{}, concurrency)
    results := make(chan time.Duration, numRequests)
    
    // 統計結果
    go func() {
        var total time.Duration
        var count int
        var min, max time.Duration = time.Hour, 0
        
        for duration := range results {
            total += duration
            count++
            
            if duration < min {
                min = duration
            }
            if duration > max {
                max = duration
            }
            
            if count%1000 == 0 {
                avg := total / time.Duration(count)
                fmt.Printf("Completed: %d, Avg: %v, Min: %v, Max: %v\n", 
                    count, avg, min, max)
            }
        }
    }()
    
    start := time.Now()
    
    for i := 0; i < numRequests; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            requestChan <- struct{}{}
            defer func() { <-requestChan }()
            
            requestStart := time.Now()
            resp, err := http.Get(targetURL)
            if err != nil {
                fmt.Printf("Request failed: %v\n", err)
                return
            }
            resp.Body.Close()
            
            results <- time.Since(requestStart)
        }()
    }
    
    wg.Wait()
    close(results)
    
    fmt.Printf("Total time: %v\n", time.Since(start))
    fmt.Printf("Requests per second: %.2f\n", 
        float64(numRequests)/time.Since(start).Seconds())
}
```

### 壓力測試與監控

```go
type LoadTester struct {
    client     *http.Client
    results    chan TestResult
    monitoring chan MonitorData
}

type TestResult struct {
    Duration   time.Duration
    StatusCode int
    Error      error
}

type MonitorData struct {
    Timestamp    time.Time
    ActiveConns  int
    MemoryUsage  uint64
    GoroutineNum int
}

func NewLoadTester() *LoadTester {
    return &LoadTester{
        client: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 100,
            },
        },
        results:    make(chan TestResult, 1000),
        monitoring: make(chan MonitorData, 1000),
    }
}

func (lt *LoadTester) RunTest(targetURL string, duration time.Duration, concurrency int) {
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()
    
    // 啟動監控
    go lt.monitor(ctx)
    
    // 啟動結果收集
    go lt.collectResults()
    
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, concurrency)
    
    for {
        select {
        case <-ctx.Done():
            wg.Wait()
            return
        default:
            wg.Add(1)
            go func() {
                defer wg.Done()
                
                semaphore <- struct{}{}
                defer func() { <-semaphore }()
                
                start := time.Now()
                resp, err := lt.client.Get(targetURL)
                duration := time.Since(start)
                
                result := TestResult{
                    Duration: duration,
                    Error:    err,
                }
                
                if resp != nil {
                    result.StatusCode = resp.StatusCode
                    resp.Body.Close()
                }
                
                lt.results <- result
            }()
        }
    }
}

func (lt *LoadTester) monitor(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            lt.monitoring <- MonitorData{
                Timestamp:    time.Now(),
                MemoryUsage:  m.Alloc,
                GoroutineNum: runtime.NumGoroutine(),
            }
        }
    }
}
```

## 效能最佳實務

### 程式碼層級優化

```go
// 避免不必要的反射
func avoidReflection(data []interface{}) {
    // 慢：使用反射
    for _, item := range data {
        v := reflect.ValueOf(item)
        if v.Kind() == reflect.String {
            // 處理字串
        }
    }
}

// 使用型別斷言
func useTypeAssertion(data []interface{}) {
    for _, item := range data {
        if str, ok := item.(string); ok {
            // 處理字串
            _ = str
        }
    }
}

// 避免在迴圈中進行字串連接
func efficientStringBuilding(items []string) string {
    if len(items) == 0 {
        return ""
    }
    
    // 估算最終字串長度
    totalLen := 0
    for _, item := range items {
        totalLen += len(item)
    }
    
    var builder strings.Builder
    builder.Grow(totalLen)
    
    for _, item := range items {
        builder.WriteString(item)
    }
    
    return builder.String()
}
```

### 編譯器優化

```go
// 使用編譯時常數
const (
    BufferSize = 4096
    MaxRetries = 3
)

// 利用內聯函式
//go:noinline
func expensiveFunction() int {
    // 複雜計算
    return 42
}

// 小函式會被自動內聯
func simpleFunction(x int) int {
    return x * 2
}
```

透過本章的學習，你已經掌握了 Go 應用程式效能優化的全面技術。從基本的效能分析工具到進階的快取策略，從記憶體管理到並發優化，這些技能將幫助你建構高效能、可擴展的 Go 應用程式。記住，效能優化是一個持續的過程，需要根據實際使用情況和效能測試結果來指導優化決策。接下來，建議你將這些優化技術應用到實際專案中，並建立完整的效能監控體系。