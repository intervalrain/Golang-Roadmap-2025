# Chapter 1.1: Variables & Data Types

在任何程式語言中，變數 (Variables) 和資料型別 (Data Types) 都是最基本的構成要素。它們讓您能夠儲存、標記和操作資料。本節將帶您了解 Go 如何處理這些核心概念。

---

## 1. Variables

**Variable** 就像一個帶有標籤的容器，您可以在裡面存放資料。在 Go 中，您需要先「宣告」一個變數，才能使用它。

### 1.1. 宣告變數 (Declaring Variables)

Go 提供了幾種宣告變數的方式。最基本的是使用 `var` 關鍵字。

```go
// 宣告一個名為 'age' 的變數，其型別為 int
var age int

// 您可以賦予它一個值
age = 30

// 也可以在宣告時直接賦予初始值
var name string = "Alice"
```

**型別推斷 (Type Inference)**

如果宣告時就提供了初始值，Go 可以自動推斷出變數的型別，讓您省略型別宣告：

```go
// Go 會自動推斷 'score' 的型別為 float64
var score = 99.5
```

### 1.2. 短變數宣告 (Short Variable Declaration)

在函式內部，您可以使用更簡潔的 `:=` 運算子來宣告並初始化變數。這是 Go 中最常用也最受歡迎的方式。

```go
// 這等同於 var address string = "123 Main St"
address := "123 Main St"

// 這種方式只能在函式內部使用
func someFunction() {
    isValid := true
    // ...
}
```

### 1.3. 零值 (Zero Values)

在 Go 中，如果您宣告一個變數但沒有給予初始值，它會自動被賦予其型別的「零值」。這可以避免未定義行為，讓程式更安全。

- `int`, `float64`: `0`
- `bool`: `false`
- `string`: `""` (空字串)
- `pointers`, `functions`, `interfaces`, `slices`, `channels`, `maps`: `nil`

```go
var count int      // count 的值為 0
var enabled bool   // enabled 的值為 false
var message string // message 的值為 ""
```

---

## 2. Basic Data Types

Go 是一種靜態型別語言，這意味著每個變數在編譯時就必須有一個確定的型別。以下是 Go 的一些基本資料型別。

### 2.1. 整數 (Integers)

Go 提供了多種整數型別，包括有符號 (`int`, `int8`, `int16`, `int32`, `int64`) 和無符號 (`uint`, `uint8`, `uint16`, `uint32`, `uint64`)。

- `int`: 最常用的整數型別，其大小取決於作業系統（32位元或64位元）。
- `uint`: 無符號整數，只能表示非負數。

### 2.2. 浮點數 (Floating-Point Numbers)

用於表示小數。

- `float32`: 32位元浮點數。
- `float64`: 64位元浮點數，更常用，精確度更高。

### 2.3. 布林值 (Booleans)

`bool` 型別只有兩個可能的值：`true` 或 `false`。

### 2.4. 字串 (Strings)

`string` 型別用於表示文字。Go 中的字串是 **不可變的 (immutable)**，一旦建立就不能修改其內容。

```go
var greeting string = "Hello, World!"
```

---

## 3. Constants

如果您有一個值，並確定它在程式執行期間永遠不會改變，您應該使用 `const` 將其宣告為常數。

```go
const pi = 3.14159
const siteName string = "Go Roadmap"

// 常數的值必須在編譯時就確定
// const randomValue = rand.Intn(10) // 這會導致編譯錯誤
```

---

## Conclusion

恭喜！您已經學會了 Go 語言中最核心的變數和資料型別。理解如何宣告變數、各種基本型別的用途、零值的概念以及常數的重要性，是編寫清晰、高效 Go 程式的第一步。

在下一個主題中，我們將探索 `Control Flow`，學習如何讓您的程式根據不同條件執行不同的程式碼路徑。

---

# Chapter 1.2: Control Flow

**Control Flow** 是指程式碼執行的順序。透過控制流程的語法，您可以讓程式做出判斷、重複執行特定區塊的程式碼，讓程式不再只是由上到下地線性執行。

---

## 1. 條件判斷 (Conditional Statements)

### 1.1. `if-else`

`if` 是最基本的條件判斷。如果條件為 `true`，就執行其後的程式碼區塊。

```go
score := 85
if score >= 60 {
    fmt.Println("及格了！")
}
```

您可以加上 `else` 來處理條件為 `false` 的情況。

```go
if score >= 90 {
    fmt.Println("優等")
} else {
    fmt.Println("甲等")
}
```

也可以使用 `else if` 來處理多個條件。

```go
if score >= 90 {
    fmt.Println("優等")
} else if score >= 80 {
    fmt.Println("甲等")
} else {
    fmt.Println("乙等")
}
```

**`if` 的簡短陳述句**

`if` 語句可以在條件判斷前，先執行一個簡短的陳述句（例如，變數宣告）。這個變數的作用域僅限於 `if-else` 區塊內。

```go
if n := rand.Intn(10); n%2 == 0 {
    fmt.Printf("%d 是偶數\n", n)
} else {
    fmt.Printf("%d 是奇數\n", n)
}
// n 在這裡無法被存取
```

---

## 2. 迴圈 (Looping)

Go 語言只有一種迴圈結構，就是 `for` 迴圈，但它有多種形式。

### 2.1. 基本 `for` 迴圈

這是最常見的形式，包含三個部分：初始陳述句、條件運算式、以及結束陳述句。

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

### 2.2. `while` 形式

您可以省略初始和結束陳述句，使其功能類似其他語言的 `while` 迴圈。

```go
sum := 1
for sum < 100 {
    sum += sum
}
```

### 2.3. 無限迴圈

如果連條件運算式都省略，就成了一個無限迴圈。您通常會搭配 `break` 或 `return` 來跳出迴圈。

```go
for {
    // 無限執行，直到被中斷
}
```

### 2.4. `for-range`

`for-range` 用於遍歷 `slice`, `array`, `string`, `map` 或 `channel`。每次迭代，它會返回索引和對應的值。

```go
// 遍歷 slice
items := []string{"apple", "banana", "cherry"}
for index, item := range items {
    fmt.Printf("索引 %d: %s\n", index, item)
}

// 如果您不需要索引，可以使用底線 (_) 來忽略它
for _, item := range items {
    fmt.Println(item)
}
```

---

## 3. 分支 (Switching)

### 3.1. `switch`

`switch` 是一個更清晰、更強大的 `if-else if` 鏈。

```go
day := "Sunday"
switch day {
case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
    fmt.Println("工作日")
case "Saturday", "Sunday":
    fmt.Println("假日")
default:
    fmt.Println("無效的日期")
}
```

在 Go 中，`case` 預設就是獨立的，執行完畢後會自動跳出 `switch`，不需要像其他語言一樣手動 `break`。

### 3.2. 無運算式的 `switch`

`switch` 後面可以不接運算式，這樣它會將每個 `case` 後的運算式當作 `true` 來進行比對。這可以寫出更簡潔的 `if-else` 邏輯。

```go
score := 88
switch {
case score >= 90:
    fmt.Println("優等")
case score >= 80:
    fmt.Println("甲等")
default:
    fmt.Println("乙等")
}
```

---

## Conclusion

您現在已經掌握了 Go 的流程控制工具：`if` 用於條件判斷，`for` 用於各種迴圈，而 `switch` 用於多重分支。這些結構是編寫任何複雜程式的基礎。

接下來，我們將學習 `Functions`，了解如何將程式碼組織成可重複使用的區塊。

---

# Chapter 1.3: Functions

**Functions** (函式) 是 Go 程式的基本建構單位。它們是執行特定任務、可重複使用的程式碼區塊。透過將程式碼組織成函式，您可以讓您的程式更有結構、更易於閱讀和維護。

---

## 1. 函式定義 (Function Definition)

使用 `func` 關鍵字來定義一個函式。一個函式包含函式名稱、參數列表、回傳值型別以及函式主體。

```go
// 一個簡單的函式，接收一個 string 參數，沒有回傳值
func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// 接收兩個 int 參數，回傳一個 int
func add(a int, b int) int {
    return a + b
}

// 如果多個參數型別相同，可以簡寫
func subtract(a, b int) int {
    return a - b
}
```

---

## 2. 多重回傳值 (Multiple Return Values)

Go 函式可以回傳多個值。這個特性經常用於同時回傳結果和錯誤狀態。

```go
// 回傳一個字串和一個錯誤
func divide(a, b float64) (float64, error) {
    if b == 0 {
        // Go 的慣例是當發生錯誤時，將錯誤作為第二個回傳值
        return 0, fmt.Errorf("除數不能為零")
    }
    return a / b, nil // nil 表示沒有錯誤發生
}
```

---

## 3. 具名回傳值 (Named Return Values)

您可以為函式的回傳值命名。這樣做可以讓程式碼更清晰，並且在函式內部，這些命名的回傳值就像變數一樣可以直接被賦值。一個沒有參數的 `return` 語句會自動回傳這些變數的當前值。

```go
// result 和 err 是具名回傳值
func divideWithNamedReturn(a, b float64) (result float64, err error) {
    if b == 0 {
        err = fmt.Errorf("除數不能為零")
        return // 會自動回傳 result 的零值 (0) 和 err
    }
    result = a / b
    return // 會自動回傳 result 和 err 的 nil 值
}
```

---

## 4. 可變參數函式 (Variadic Functions)

如果函式的最後一個參數型別前有 `...`，表示這是一個可變參數函式，它可以接收任意數量的該型別參數。

在函式內部，這個可變參數會被當作一個該型別的 `slice` 來處理。

```go
// 接收任意數量的 int
func sumAll(numbers ...int) int {
    total := 0
    for _, number := range numbers {
        total += number
    }
    return total
}

// 使用方式：
// sumAll(1, 2, 3)       // 結果為 6
// sumAll(10, 20)        // 結果為 30
// nums := []int{4, 5, 6}
// sumAll(nums...)     // 如果您已經有一個 slice，可以用 ... 將其展開傳入
```

---

## Conclusion

函式是組織和重用程式碼的強大工具。透過學習定義函式、處理多重回傳值以及使用可變參數，您現在可以開始編寫更有組織、更模組化的 Go 程式了。

接下來，我們將深入探討一個更進階但非常重要的主題：`Pointers`。

---

# Chapter 1.4: Pointers

**Pointers** (指標) 是 Go 語言中一個強大但初學者可能會感到困惑的概念。簡單來說，一個指標是一個儲存了另一個變數 **記憶體位址** 的變數。

想像一下，變數是儲存資料的房子，而指標就是記錄這些房子地址的筆記本。透過地址，您可以找到並修改房子的內部（即變數的值）。

---

## 1. 什麼是指標？ (What are Pointers?)

在 Go 中，對一個變數取其記憶體位址，是使用 `&` 運算子。

```go
x := 10
p := &x // p 是一個指標，它儲存了 x 的記憶體位址

fmt.Printf("x 的值: %d\n", x)
fmt.Printf("x 的記憶體位址: %p\n", p)
```

`p` 的型別是 `*int` (讀作 "pointer to int")，表示它是一個指向 `int` 型別變數的指標。

## 2. 解參考 (Dereferencing)

當您擁有一個指標時，您可以使用 `*` 運算子來「解參考」，也就是取得該指標指向的記憶體位址上所儲存的 **值**。

```go
fmt.Printf("p 指向的值: %d\n", *p) // *p 會取得 x 的值，所以結果是 10
```

您也可以透過解參考來修改原始變數的值。

```go
*p = 20 // 透過指標 p，將 x 的值修改為 20
fmt.Printf("x 現在的值: %d\n", x) // 結果會是 20
```

## 3. 為何使用指標？ (Why Use Pointers?)

使用指標主要有兩個原因：

1.  **效率**: 當您將大型的資料結構（例如一個很大的 `struct`）傳遞給函式時，如果直接傳遞值，Go 會複製整個資料結構，這會消耗時間和記憶體。如果傳遞指標，函式只會複製一個很小的記憶體位址，效率更高。

2.  **修改原始資料**: Go 的函式參數預設是「傳值」(pass-by-value)，意味著函式內部對參數的修改不會影響到原始變數。如果您希望函式能夠修改傳入的原始變數，您就必須傳遞該變數的指標。

```go
func addOne(val int) {
    val = val + 1 // 這裡修改的是 val 的副本
}

func addOneWithPointer(val *int) {
    *val = *val + 1 // 這裡修改的是指標指向的原始值
}

func main() {
    i := 10
    addOne(i)
    fmt.Println("addOne 之後:", i) // i 仍然是 10

    addOneWithPointer(&i)
    fmt.Println("addOneWithPointer 之後:", i) // i 變成了 11
}
```

## 4. `nil` 指標

一個指標的零值是 `nil`。一個 `nil` 指標不指向任何記憶體位址。對一個 `nil` 指標進行解參考會導致執行時錯誤 (runtime panic)。

```go
var p *int // p 是一個 nil 指標
// fmt.Println(*p) // 這會引發 panic
```

---

## Conclusion

指標是 Go 中一個不可或缺的工具，它讓您能夠更有效率地管理記憶體，並在函式間共享和修改資料。雖然一開始可能不易掌握，但理解指標的工作原理對於編寫高效能的 Go 程式至關重要。

接下來，我們將結合指標和自訂型別，學習 `Structs & Methods`。

---

# Chapter 1.5: Structs & Methods

**Structs** (結構) 是 Go 中用來建立自訂資料型別的工具。它是一個將多個不同型別的欄位 (fields) 集合在一起的複合型別。如果您有物件導向程式設計的經驗，可以將 `struct` 看作是沒有方法的 class。

**Methods** (方法) 則是附加到特定型別上的函式。當一個函式附加到一個 `struct` 上時，它就成為了該 `struct` 的方法。

---

## 1. 定義 Struct (Defining a Struct)

使用 `type` 和 `struct` 關鍵字來定義一個 `struct`。

```go
// 定義一個名為 Person 的 struct
type Person struct {
    FirstName string
    LastName  string
    Age       int
}
```

## 2. 建立和初始化 Struct

```go
// 建立一個 Person 型別的變數
var p1 Person
p1.FirstName = "Alice"
p1.LastName = "Smith"
p1.Age = 30

// 使用 struct literal 來初始化
p2 := Person{
    FirstName: "Bob",
    LastName:  "Johnson",
    Age:       25,
}

// 如果您知道欄位的順序，可以省略欄位名稱
p3 := Person{"Charlie", "Brown", 35}

// 也可以建立一個指向 struct 的指標
p4 := &Person{FirstName: "Diana"}
```

## 3. 定義方法 (Defining a Method)

方法是一個帶有「接收者 (receiver)」的函式。接收者出現在 `func` 關鍵字和方法名稱之間。

```go
// 為 Person struct 定義一個名為 FullName 的方法
// (p Person) 是接收者
func (p Person) FullName() string {
    return p.FirstName + " " + p.LastName
}
```

`p` 在這裡被稱為接收者，您可以任意命名它（例如 `person` 或 `self`），但慣例是使用該型別名稱的第一個小寫字母。

## 4. 指標接收者 vs. 值接收者 (Pointer vs. Value Receivers)

方法的接收者可以是值型別，也可以是指標型別。

- **值接收者 (`func (p Person) ...`)**: 方法會得到該 `struct` 的一個副本。在方法內部對 `struct` 的任何修改，都不會影響到原始的 `struct`。
- **指標接收者 (`func (p *Person) ...`)**: 方法會得到一個指向該 `struct` 的指標。在方法內部對 `struct` 的修改，會影響到原始的 `struct`。這是非常常見且高效的做法。

```go
// 使用指標接收者來修改 Age
func (p *Person) SetAge(age int) {
    p.Age = age
}

func main() {
    p := Person{"John", "Doe", 40}
    fmt.Println(p.FullName()) // "John Doe"

    p.SetAge(41)
    fmt.Println(p.Age) // 41
}
```

**Go 的自動解參考**: 當您呼叫一個需要指標接收者的方法時（如 `p.SetAge(42)`），Go 會自動為您轉換成 `(&p).SetAge(42)`，讓程式碼更簡潔。

---

## Conclusion

`Structs` 和 `Methods` 是 Go 實現「封裝」的基礎，讓您可以建立具有行為的自訂資料型別。這是 Go 程式設計中建構複雜、有組織應用程式的核心。

接下來，我們將學習 Go 的依賴管理系統：`Go Modules`。

```
