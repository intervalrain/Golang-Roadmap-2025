# Chapter 2.1: Interfaces

**Interfaces** (介面) 是 Go 語言中最強大、最核心的概念之一。它提供了一種定義「行為」的方式。一個介面型別定義了一組方法的集合，任何實作了這些方法的具體型別，我們就稱之為「滿足」了該介體。

Go 的介面是隱式實現的。這意味著您不需要像其他語言（如 Java）那樣明確宣告 `implements` 某個介面。只要您的型別擁有介面所定義的全部方法，它就自動地、隱式地滿足了該介面。

---

## 1. 定義介面 (Defining an Interface)

使用 `type` 和 `interface` 關鍵字來定義一個介面。

```go
// Shaper 是一個介面
// 任何擁有 Area() float64 方法的型別都滿足 Shaper 介面
type Shaper interface {
    Area() float64
}
```

## 2. 實作介面 (Implementing an Interface)

讓我們定義幾個 `struct`，並為它們實作 `Area()` 方法。

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// 為 Rectangle 實作 Area() 方法
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

type Circle struct {
    Radius float64
}

// 為 Circle 實作 Area() 方法
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}
```

因為 `Rectangle` 和 `Circle` 都定義了 `Area() float64` 方法，所以它們都隱式地滿足了 `Shaper` 介面。

## 3. 使用介面 (Using Interfaces)

介面的威力在於，您可以編寫一個只接受介面型別的函式。這個函式可以操作任何滿足該介面的具體型別，而不需要知道這些具體型別的內部細節。

```go
// 這個函式可以接收任何滿足 Shaper 介面的型別
func PrintShapeArea(s Shaper) {
    fmt.Printf("這個形狀的面積是: %0.2f\n", s.Area())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circ := Circle{Radius: 3}

    PrintShapeArea(rect) // 傳入 Rectangle
    PrintShapeArea(circ) // 傳入 Circle
}
```

## 4. 空介面 (The Empty Interface)

一個不包含任何方法的介面被稱為「空介面」，寫作 `interface{}`。因為任何型別都至少實作了零個方法，所以 **任何型別都滿足空介面**。

空介面可以用於儲存任意型別的值，這在您需要處理未知型別的資料時非常有用。

```go
// 接受任何型別的參數
func PrintAnything(v interface{}) {
    fmt.Printf("值: %v, 型別: %T\n", v, v)
}

PrintAnything(10)
PrintAnything("hello")
PrintAnything(Rectangle{Width: 1, Height: 2})
```

## 5. 型別斷言 (Type Assertions)

當您有一個介面型別的變數時，您可能需要知道它底層儲存的具體型別是什麼。這時就需要使用「型別斷言」。

語法是 `v.(T)`，其中 `v` 是介面型別的變數，`T` 是您要斷言的具體型別。

```go
var i interface{} = "hello"

// 斷言 i 儲存的是 string
s := i.(string)
fmt.Println(s) // "hello"

// 如果斷言失敗，會引發 panic
// f := i.(float64) // panic: interface conversion: interface {} is string, not float64

// 更安全的方式是使用「comma, ok」的語法
f, ok := i.(float64)
if ok {
    fmt.Printf("i 是 float64，值為 %f\n", f)
} else {
    fmt.Println("i 不是 float64")
}
```

---

## Conclusion

介面是 Go 實現多型 (polymorphism) 和編寫高彈性、可擴充套件程式碼的關鍵。透過定義行為而非具體資料，介面讓您的程式碼耦合度更低，更易於測試和維護。

接下來，我們將學習 Go 內建的測試工具，了解如何為您的程式碼編寫 `Testing Basics`。

---

# Chapter 2.2: Error Handling

在 Go 的設計哲學中，錯誤處理是至關重要的一環。Go 將 **錯誤視為一般的值**，而不是像其他語言那樣使用 `try-catch` 機制來處理「例外」。這種明確的、可預期的錯誤處理方式是 Go 語言健壯性的關鍵來源。

---

## 1. `error` 型別

Go 的 `error` 是一個內建的介面型別，其定義非常簡單：

```go
type error interface {
    Error() string
}
```

任何實作了 `Error() string` 方法的型別，都滿足 `error` 介面。這意味著您可以非常靈活地建立自訂的錯誤型別。

## 2. 慣用的錯誤處理模式

Go 最常見的錯誤處理模式是：讓函式在回傳結果的同時，也回傳一個 `error` 型別的值作為最後一個回傳值。如果函式執行成功，`error` 的值會是 `nil`；如果失敗，它會包含一個描述錯誤的非 `nil` 值。

呼叫端則必須 **立即檢查** 這個 `error`。

```go
import "strconv"

func main() {
    s := "123a"
    // strconv.Atoi 會嘗試將字串轉換為整數
    // 它會回傳 (int, error)
    n, err := strconv.Atoi(s)
    if err != nil {
        // 如果 err 不是 nil，表示發生了錯誤
        fmt.Printf("轉換失敗: %v\n", err)
        return // 提早返回，不再繼續執行
    }

    fmt.Printf("成功轉換為: %d\n", n)
}
```

這種 `if err != nil` 的模式在 Go 程式碼中隨處可見，它強迫開發者正視並處理每一個可能出錯的地方。

## 3. 建立錯誤

有兩種簡單的方式可以建立 `error`：

- **`errors.New()`**: 建立一個只包含簡單錯誤訊息的 `error`。
- **`fmt.Errorf()`**: 格式化一個錯誤訊息，功能更強大，可以包含變數內容。

```go
import (
    "errors"
    "fmt"
)

func checkAge(age int) error {
    if age < 0 {
        return errors.New("年齡不能是負數")
    }
    if age < 18 {
        return fmt.Errorf("年齡 %d 過小，必須年滿 18 歲", age)
    }
    return nil
}
```

## 4. 錯誤的包裝 (Wrapping Errors)

有時候，一個錯誤是由另一個更底層的錯誤引起的。Go 1.13 引入了「錯誤包裝」機制，讓您可以將錯誤層層包裹起來，形成一個錯誤鏈，同時保留完整的上下文資訊。

- 使用 `fmt.Errorf` 中的 `%w` 動詞來包裝錯誤。
- 使用 `errors.Is()` 來檢查錯誤鏈中是否 **包含** 特定的目標錯誤。
- 使用 `errors.As()` 來檢查錯誤鏈中是否有特定 **型別** 的錯誤，並將其取出。

```go
// 假設這是從資料庫層回傳的錯誤
var ErrNotFound = errors.New("not found")

func findUser(id int) error {
    // ... 模擬查詢失敗
    return fmt.Errorf("查詢使用者 %d 失敗: %w", id, ErrNotFound)
}

func main() {
    err := findUser(123)
    if err != nil {
        // 使用 errors.Is() 檢查錯誤鏈中是否包含 ErrNotFound
        if errors.Is(err, ErrNotFound) {
            fmt.Println("使用者不存在，請檢查您的輸入。")
        } else {
            fmt.Println("發生未知錯誤:", err)
        }
    }
}
```

---

## Conclusion

Go 的錯誤處理機制鼓勵開發者編寫清晰、明確且健壯的程式碼。將錯誤視為值，並透過 `if err != nil` 的模式進行處理，是每個 Go 開發者都必須掌握的核心技能。透過錯誤包裝，您可以為問題追蹤提供更豐富的上下文。

---

# Chapter 2.3: Package Management

**Packages** (套件) 是 Go 組織和重用程式碼的基本單位。一個套件就是一個目錄中的所有 Go 檔案的集合。良好的套件設計是建構大型、可維護 Go 應用程式的基礎。

我們在 `Chapter 1.6` 中已經介紹了 `Go Modules`，它是管理外部依賴的工具。本節將更深入地探討如何在您自己的專案中組織和使用套件。

---

## 1. 套件宣告與可見性 (Declaration and Visibility)

- **套件宣告**: 一個目錄中的所有 Go 檔案都必須屬於同一個套件，並且必須在檔案開頭使用 `package [package_name]` 來宣告。
- **可見性規則**: Go 透過 **名稱的大小寫** 來決定一個識別字（變數、常數、函式、`struct` 等）是否能被套件外部存取。
    - **大寫開頭 (Exported)**: 如果名稱以大寫字母開頭（例如 `MyVariable`, `DoSomething`），則它是「導出的」，可以被其他套件引用。
    - **小寫開頭 (Unexported)**: 如果名稱以小寫字母開頭（例如 `myVariable`, `doSomething`），則它是「未導出的」，只能在它自己所在的套件內部使用。

## 2. 專案結構範例 (Example Project Structure)

一個典型的 Go 專案結構可能如下：

```
my-project/
├── go.mod
├── main.go                 # main 套件
└── calculator/
    ├── calculator.go
    └── calculator_test.go
```

- **`main.go`**:

```go
package main

import (
    "fmt"
    "my-project/calculator" // 匯入我們自己的 calculator 套件
)

func main() {
    result := calculator.Add(10, 5) // 呼叫 calculator 套件中的 Add 函式
    fmt.Println("結果是:", result)

    // 下面這行會導致編譯錯誤，因為 subtract 是小寫開頭，無法被外部存取
    // res2 := calculator.subtract(10, 5)
}
```

- **`calculator/calculator.go`**:

```go
package calculator

// Add 函式因為是大寫開頭，所以可以被其他套件呼叫
func Add(a, b int) int {
    return a + b
}

// subtract 函式是小寫開頭，只能在 calculator 套件內部使用
func subtract(a, b int) int {
    return a - b
}
```

## 3. `internal` 套件

Go 有一個特殊的目錄名稱 `internal`。放在 `internal` 目錄下的套件只能被其 **直屬父目錄** 以及父目錄的子目錄中的程式碼所引用。

這提供了一種強制的保護機制，確保某些僅供內部使用的套件不會被專案外部的其他部分意外引用。

```
my-project/
├── internal/
│   └── auth/             # auth 套件只能被 my-project 內的程式碼引用
│       └── auth.go
├── main.go                 # 可以引用 auth
└── api/
    └── server.go         # 也可以引用 auth

another-project/
└── main.go                 # 無法引用 my-project/internal/auth
```

---

## Conclusion

透過將相關功能的程式碼組織到獨立的套件中，並利用大小寫規則來控制可見性，您可以建立出結構清晰、易於理解和維護的 Go 專案。`internal` 套件更進一步提供了程式碼隔離的強大工具。

---

# Chapter 2.4: Testing Basics

測試是確保軟體品質、進行重構以及驗證功能正確性的重要環節。Go 語言內建了強大且易用的測試工具鏈，讓開發者可以輕鬆地為他們的程式碼編寫測試。

---

## 1. 測試檔案的約定 (Test File Conventions)

Go 的測試工具依賴一套簡單的檔案命名和函式簽名約定：

- **檔案名稱**: 測試檔案必須以 `_test.go` 結尾（例如 `calculator_test.go`）。
- **函式簽名**: 測試函式的名稱必須以 `Test` 開頭，並接收一個 `*testing.T` 型別的參數（例如 `func TestAdd(t *testing.T) { ... }`）。
- **位置**: 測試檔案和被測試的程式碼通常放在同一個套件（同一個目錄）下。

## 2. 編寫一個基本測試 (Writing a Basic Test)

讓我們為 `Chapter 2.3` 中的 `calculator` 套件編寫一個測試。

- **`calculator/calculator.go`** (與上一節相同):

```go
package calculator

func Add(a, b int) int {
    return a + b
}
```

- **`calculator/calculator_test.go`**:

```go
package calculator

import "testing"

func TestAdd(t *testing.T) {
    // 準備測試資料
    a := 10
    b := 5
    expected := 15

    // 執行被測試的函式
    result := Add(a, b)

    // 斷言結果是否符合預期
    if result != expected {
        // t.Errorf 會將此測試標記為失敗，並輸出錯誤訊息
        t.Errorf("Add(%d, %d) = %d; 預期為 %d", a, b, result, expected)
    }
}
```

## 3. 執行測試 (Running Tests)

打開終端機，進入到您的專案目錄中，然後執行 `go test` 指令。

- **`go test`**: 執行當前目錄及其子目錄下的所有測試。
- **`go test ./...`**: 從專案根目錄執行所有套件的測試。
- **`go test -v`**: `-v` (verbose) 參數會顯示更詳細的測試過程，包括每個測試函式的名稱和執行結果。

當您在 `calculator` 目錄下執行 `go test -v`，您應該會看到類似以下的成功輸出：

```
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
PASS
ok      my-project/calculator 0.001s
```

## 4. 表格驅動測試 (Table-Driven Tests)

當您需要為同一個函式測試多組不同的輸入和輸出時，表格驅動測試是一種非常常見且推薦的模式。它讓您可以輕鬆地新增、刪除或修改測試案例。

```go
func TestAddTableDriven(t *testing.T) {
    // 定義測試案例的表格
    testCases := []struct {
        name     string // 測試案例名稱
        a, b     int    // 輸入
        expected int    // 預期輸出
    }{
        {"正數相加", 2, 3, 5},
        {"負數相加", -2, -3, -5},
        {"與零相加", 7, 0, 7},
    }

    // 遍歷所有測試案例
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := Add(tc.a, tc.b)
            if result != tc.expected {
                t.Errorf("Add(%d, %d) = %d; 預期為 %d", tc.a, tc.b, result, tc.expected)
            }
        })
    }
}
```

`t.Run` 可以讓您建立子測試，這樣在輸出結果時會更清晰，並且可以獨立執行某個子測試。

---

## Conclusion

Go 內建的測試工具提供了一個簡單而強大的方式來確保您的程式碼品質。透過遵循 `_test.go` 的命名約定和編寫測試函式，特別是採用表格驅動測試的模式，您可以建立一個健壯的測試套件，為您的專案保駕護航。

---

# Chapter 2.5: Goroutines 入門

**Goroutine** 是 Go 語言並發 (concurrency) 模型的核心。您可以將它想像成一個由 Go runtime 所管理的、非常輕量的執行緒。與作業系統的執行緒相比，Goroutine 的創建成本極低，您可以在一個程式中輕鬆地啟動成千上萬個 Goroutine。

並發 (Concurrency) 不等於並行 (Parallelism)。並發是關於「處理」多件事情，而並行是關於「執行」多件事情。Go 透過 Goroutine 和 Channels 讓編寫並發程式變得簡單。

---

## 1. 啟動一個 Goroutine (Starting a Goroutine)

啟動一個 Goroutine 非常簡單，只需要在函式呼叫前加上 `go` 關鍵字即可。

```go
import (
    "fmt"
    "time"
)

func say(s string) {
    for i := 0; i < 3; i++ {
        fmt.Println(s)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    go say("World") // 啟動一個新的 Goroutine
    say("Hello")    // 在主 Goroutine 中執行
}
```

當您執行這段程式碼，您會發現 "Hello" 和 "World" 的輸出是交錯的。這是因為 `say("World")` 在一個新的 Goroutine 中執行，而 `say("Hello")` 在主 Goroutine（也就是 `main` 函式本身）中執行。它們是並發執行的。

## 2. Goroutine 與主函式 (Goroutines and the Main Function)

需要特別注意的是：當 `main` 函式結束時，整個程式就會退出，**所有正在執行的 Goroutine 也會被立即終止**，不論它們是否執行完畢。

觀察以下範例：

```go
func main() {
    go say("I am a goroutine")
    // main 函式在這裡沒有做任何等待，它會立刻結束
}
```

如果您執行這段程式碼，您很可能什麼也看不到。因為在 `say` 這個 Goroutine 有機會印出任何東西之前，`main` 函式就已經執行完畢並退出了。

為了讓 Goroutine 有時間執行，我們需要一種方法來等待它完成。在上面的第一個範例中，`say("Hello")` 的執行恰好給了 `say("World")` 足夠的執行時間。但這不是一個可靠的方法。

在後續的章節中，我們將學習 `Channels` 和 `sync` 套件，它們提供了更可靠、更優雅的方式來協調和同步 Goroutine。

---

## Conclusion

恭喜！您已經完成了第二章的學習。

您已經了解了 Go 中一些最强大的特性：介面提供了靈活性，明確的錯誤處理確保了健壯性，套件系統幫助組織程式碼，內建的測試工具保障了品質，而 Goroutine 則為高效的並發程式設計打開了大門。

在下一章，我們將深入 `並發程式設計 (Concurrency)`，學習如何使用 `Channels` 和 `sync` 套件來駕馭 Goroutine 的強大威力。

```