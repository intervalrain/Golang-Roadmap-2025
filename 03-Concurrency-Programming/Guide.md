# Chapter 3.1: Goroutines

我們在上一章已經對 Goroutine 有了初步的認識。它是 Go 並發模型的核心，由 Go runtime 管理的輕量級執行緒。本節我們將更深入地探討 Goroutine 的使用，以及如何透過 `sync.WaitGroup` 來等待它們執行完成。

---

## 1. `sync.WaitGroup`

在上一章我們看到，如果 `main` 函式結束，所有的 Goroutine 都會被強制終止。為了可靠地等待 Goroutine 完成，我們可以使用 `sync` 套件中的 `WaitGroup`。

`WaitGroup` 是一個計數信號量，可以用來等待一組 Goroutine 的完成。它有三個主要方法：

- **`Add(delta int)`**: 將計數器增加 `delta`。通常 `delta` 就是您要等待的 Goroutine 的數量。
- **`Done()`**: 將計數器減一。每個 Goroutine 在完成任務後都應該呼叫此方法。
- **`Wait()`**: 阻塞當前的執行緒，直到計數器歸零。

### 範例：使用 WaitGroup

```go
import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    // 在函式退出時，通知 WaitGroup 任務已完成
    defer wg.Done()

    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second) // 模擬耗時操作
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    // 啟動 3 個 worker goroutine
    for i := 1; i <= 3; i++ {
        wg.Add(1) // 每啟動一個 goroutine，計數器加 1
        go worker(i, &wg)
    }

    // 等待所有 goroutine 完成
    fmt.Println("Waiting for workers to finish...")
    wg.Wait()
    fmt.Println("All workers done.")
}
```

在這個範例中，`main` 函式會一直阻塞在 `wg.Wait()`，直到三個 `worker` Goroutine 都呼叫了 `wg.Done()`，使計數器歸零為止。這確保了主程式會等待所有並發任務完成後才退出。

## 2. 匿名函式 Goroutine

您也可以直接使用匿名函式來快速啟動一個 Goroutine，這在處理一些簡單的並發任務時非常方便。

```go
func main() {
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        defer wg.Done()
        fmt.Println("I am an anonymous goroutine!")
    }()

    wg.Wait()
}
```

---

## Conclusion

Goroutine 是 Go 實現並發的基礎。單獨使用 Goroutine 已經可以執行非同步任務，但如果沒有一個可靠的同步機制（如 `WaitGroup`），我們將無法保證這些任務能夠執行完畢。`WaitGroup` 提供了一個簡單而強大的方法來等待一組 Goroutine 的完成。

然而，僅僅等待任務完成是不夠的。在不同的 Goroutine 之間安全地傳遞資料是並發程式設計的另一個關鍵挑戰。接下來，我們將學習 Go 的一個標誌性特性：`Channels`，它被設計用來解決這個問題。

---

# Chapter 3.2: Channels

**Channels** (通道) 是 Go 中專為 Goroutine 之間通訊而設計的管道。您可以把它想像成一個傳送帶，一個 Goroutine 可以將特定型別的值放入傳送帶的一端，而另一個 Goroutine 則可以從另一端取出。

Channel 的核心哲學是：「不要透過共享記憶體來通訊；而是透過通訊來共享記憶體。」 (Do not communicate by sharing memory; instead, share memory by communicating.)

---

## 1. 建立和使用 Channel

- **建立 Channel**: 使用 `make(chan val-type)` 來建立一個 Channel。Channel 是有型別的，只能傳遞特定型別的資料。
- **傳送 (Send)**: 使用 `channel <- value` 語法將一個值傳送到 Channel。
- **接收 (Receive)**: 使用 `<-channel` 語法從 Channel 接收一個值。

預設情況下，傳送和接收操作都是 **阻塞的**。這意味著：
- 當一個 Goroutine 向 Channel 傳送資料時，它會被阻塞，直到有另一個 Goroutine 從該 Channel 接收資料。
- 當一個 Goroutine 從 Channel 接收資料時，它會被阻塞，直到有另一個 Goroutine 向該 Channel 傳送資料。

```go
func main() {
    // 建立一個可以傳遞 string 型別的 channel
    messages := make(chan string)

    go func() {
        // 將 "ping" 傳入 channel
        // 這個操作會阻塞，直到 main goroutine 準備好接收
        messages <- "ping"
    }()

    // 從 channel 接收一個值
    // 這個操作會阻塞，直到有 goroutine 將值傳入
    msg := <-messages
    fmt.Println(msg)
}
```

## 2. 緩衝 Channel (Buffered Channels)

預設建立的是「無緩衝 Channel」。您也可以在 `make` 時提供第二個整數參數來建立一個「緩衝 Channel」：`make(chan int, 100)`。

- 向一個 **未滿** 的緩衝 Channel 傳送資料 **不會** 發生阻塞。
- 從一個 **非空** 的緩衝 Channel 接收資料 **不會** 發生阻塞。

只有在緩衝區滿時傳送，或在緩衝區空時接收，才會發生阻塞。

```go
func main() {
    // 建立一個緩衝大小為 2 的 string channel
    messages := make(chan string, 2)

    // 因為有緩衝，這兩次傳送不會阻塞
    messages <- "buffered"
    messages <- "channel"

    // 取出值
    fmt.Println(<-messages)
    fmt.Println(<-messages)
}
```

## 3. `for-range` 和 `close`

- **`for-range`**: 您可以使用 `for-range` 迴圈來持續地從 Channel 中接收資料，直到該 Channel 被關閉。
- **`close`**: 傳送者可以透過 `close(channel)` 來關閉一個 Channel，表示不會再有新的值被傳送進來。

接收者可以透過接收操作的第二個回傳值來判斷 Channel 是否已被關閉：`v, ok := <-ch`。如果 `ok` 是 `false`，表示 Channel 已被關閉且其中已沒有值可接收。

```go
func producer(ch chan int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch) // 完成後關閉 channel
}

func main() {
    ch := make(chan int, 5)
    go producer(ch)

    // for-range 會自動處理 channel 的關閉
    for val := range ch {
        fmt.Println("接收到:", val)
    }
}
```

**重要**: 只有傳送者才應該關閉 Channel，接收者永遠不應該關閉。向一個已關閉的 Channel 傳送資料會引發 `panic`。

---

## Conclusion

Channel 是 Go 並發程式設計的基石，它提供了一種型別安全、同步的 Goroutine 通訊方式。透過阻塞式的傳送和接收，以及使用緩衝區，您可以精確地控制 Goroutine 之間的協作流程。

接下來，我們將學習 `Select Statement`，它讓您可以同時等待多個 Channel 操作。

---

# Chapter 3.3: Select Statement

Go 的 `select` 語句讓一個 Goroutine 可以同時等待多個通訊操作（在多個 Channel 上進行傳送或接收）。

`select` 會一直阻塞，直到其中一個 `case` 可以執行，然後它就會執行那個 `case`。如果有多個 `case` 同時就緒，`select` 會隨機選擇一個執行。這可以避免飢餓，並確保公平性。

---

## 1. `select` 的基本用法

`select` 的結構類似 `switch`，但它的 `case` 都是 Channel 操作。

```go
func main() {
    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        c1 <- "one"
    }()
    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
    }()

    // 我們需要等待兩次，所以迴圈兩次
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
            fmt.Println("received", msg2)
        }
    }
}
```

在這個範例中，`select` 會等待 `c1` 或 `c2` 之一變為可讀。第一次迴圈，一秒後 `c1` 就緒，第一個 `case` 被執行。第二次迴圈，`c1` 已經沒有值了，`select` 會繼續等待，直到兩秒後 `c2` 就緒，第二個 `case` 被執行。

## 2. `default` Case

`select` 可以有一個 `default` case。當 `select` 中的其他 `case` 都沒有就緒時，`default` case 會被執行。

這可以用來實現 **非阻塞** 的 Channel 操作。

```go
ch := make(chan string)

select {
case msg := <-ch:
    fmt.Println("received message", msg)
default:
    // 如果 ch 中沒有值可以接收，立即執行這裡
    fmt.Println("no message received")
}
```

## 3. 超時處理 (Timeouts)

`default` case 對於防止 Goroutine 在 Channel 操作上無限期阻塞非常有用。一個常見的應用場景是實現操作超時。

`time.After` 函式會回傳一個 Channel，它會在指定的時間過後接收到一個值。

```go
c1 := make(chan string, 1)
go func() {
    time.Sleep(2 * time.Second)
    c1 <- "result 1"
}()

select {
case res := <-c1:
    fmt.Println(res)
case <-time.After(1 * time.Second):
    // 如果 1 秒內沒有從 c1 接收到值，這個 case 就會被執行
    fmt.Println("timeout 1")
}
```

---

## Conclusion

`select` 是 Go 中一個強大的控制結構，它賦予了 Goroutine 監聽多個 Channel 的能力。透過結合 `default` case 或 `time.After`，您可以輕鬆實現非阻塞操作和超時控制，編寫出更健壯、反應更靈敏的並發程式。

接下來，我們將探索 `sync` 套件中除了 `WaitGroup` 之外的其他工具，例如 `Mutex`。

---

# Chapter 3.4: Sync Package

雖然 Go 推崇「透過通訊來共享記憶體」，但在某些情況下，我們仍然需要使用傳統的共享記憶體方式來處理並發，例如多個 Goroutine 需要存取同一個變數時。為了保證在這種情況下的資料安全，`sync` 套件提供了必要的同步原語 (synchronization primitives)。

我們已經學習了 `sync.WaitGroup`，本節將介紹另外兩個重要的工具：`sync.Mutex` 和 `sync.Once`。

---

## 1. `sync.Mutex`

`Mutex` 是「互斥鎖 (mutual exclusion lock)」的縮寫。它用來保護一段程式碼在同一時間只能被一個 Goroutine 執行，從而避免多個 Goroutine 同時修改共享資源而導致的「競爭條件 (race condition)」。

- **`Lock()`**: 獲取鎖。如果鎖已經被其他 Goroutine 持有，則當前的 Goroutine 會被阻塞，直到鎖被釋放。
- **`Unlock()`**: 釋放鎖。一個 Goroutine 在完成對共享資源的存取後，必須呼叫 `Unlock()`，以便其他等待的 Goroutine 可以獲取鎖。

```go
// SafeCounter 是一個線程安全的計數器
type SafeCounter struct {
    mu      sync.Mutex
    counter map[string]int
}

// Inc 方法安全地對計數器進行遞增
func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()   // 在存取 map 前鎖定
    defer c.mu.Unlock() // 使用 defer 確保在函式退出時解鎖
    c.counter[key]++
}

// Value 方法安全地讀取計數器的值
func (c *SafeCounter) Value(key string) int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.counter[key]
}
```

在上面的範例中，任何對 `counter` map 的讀寫操作都被 `mu.Lock()` 和 `mu.Unlock()` 保護起來，確保了併發安全。

## 2. `sync.Once`

`sync.Once` 是一個非常有用的工具，它能保證某個動作在整個程式的生命週期中 **只執行一次**。

- **`Do(f func())`**: `Do` 方法接收一個函式作為參數。只有第一次呼叫 `Do` 時，這個函式 `f` 才會被執行。無論多少個 Goroutine 同時呼叫 `Do`，`f` 也只會被執行一次。

這在初始化單例物件或設定共享資源時非常有用。

```go
var once sync.Once

func initialize() {
    fmt.Println("Initializing...")
    // ... 執行一些耗時的初始化操作
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // 儘管有 10 個 goroutine，initialize() 只會被執行一次
            once.Do(initialize)
        }()
    }
    wg.Wait()
    fmt.Println("Done.")
}
```

---

## Conclusion

`sync` 套件提供了強大的工具來處理傳統的共享記憶體並發模型。`Mutex` 是保護共享資源、防止競爭條件的關鍵，而 `Once` 則為「只執行一次」的初始化場景提供了簡潔、安全的解決方案。

在了解了這些基本的並發原語後，我們將在下一節探討一些基於它們建構的、更高級的 `Concurrency Patterns`。

---

# Chapter 3.5: Concurrency Patterns

在掌握了 Goroutine、Channel 和 `sync` 套件這些基礎工具後，我們可以開始探索一些更高級的、可重用的並發模式。這些模式是 Go 社群在長期實踐中總結出的解決常見並發問題的優雅方案。

---

## 1. Fan-Out, Fan-In

這是一個非常強大的模式，用於將一個耗時的任務分散給多個 Goroutine 並行處理，然後再將它們的結果匯總起來。

- **Fan-Out**: 一個 Goroutine（通常是 `producer`）產生任務，並將這些任務發送到一個 Channel。多個 `worker` Goroutine 從這個 Channel 中取出任務並行處理。
- **Fan-In**: 多個 `worker` Goroutine 將它們的處理結果發送到 **同一個** 結果 Channel。一個 `consumer` Goroutine 從結果 Channel 中讀取所有結果。

```go
// worker 函式模擬一個耗時的計算
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("worker %d started job %d\n", id, j)
        time.Sleep(time.Second)
        fmt.Printf("worker %d finished job %d\n", id, j)
        results <- j * 2
    }
}

func main() {
    numJobs := 5
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    // Fan-Out: 啟動 3 個 worker goroutines
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // 將 5 個任務發送到 jobs channel
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    // Fan-In: 收集所有任務的結果
    for a := 1; a <= numJobs; a++ {
        <-results
    }
}
```

## 2. Rate Limiting (速率限制)

速率限制是控制資源使用的一個重要機制，例如限制對某个 API 的請求頻率。Go 的 Channel 和 `time.Ticker` 可以非常優雅地實現速率限制。

`time.NewTicker` 會回傳一個 `Ticker` 物件，它包含一個 Channel，會以固定的時間間隔向該 Channel 發送事件。

```go
func main() {
    // 建立一個每 200 毫秒觸發一次的 Ticker
    ticker := time.NewTicker(200 * time.Millisecond)
    defer ticker.Stop()

    // 模擬有 5 個請求需要處理
    requests := make(chan int, 5)
    for i := 1; i <= 5; i++ {
        requests <- i
    }
    close(requests)

    // 處理請求
    for req := range requests {
        <-ticker.C // 等待 ticker 觸發
        fmt.Println("processing request", req, time.Now())
    }
}
```

在這個範例中，`<-ticker.C` 操作會阻塞，直到 Ticker 觸發。這確保了我們的請求處理速率不會超過每 200 毫秒一次。

---

## Conclusion

恭喜！您已經完成了第三章，深入探索了 Go 的並發世界。

您不僅掌握了 Goroutine、Channel、Select 和 Sync 等核心工具，還學習了如何將它們組合起來，形成如 Fan-Out/Fan-In 和速率限制等強大的並發模式。這些知識是您編寫高效、可擴展的現代化 Go 應用程式的關鍵。

在下一章，我們將把目光轉向一個非常實用的領域：`HTTP & Web 開發`。
