package main

import "fmt"

// 1. 定義一個 struct
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// 2. 定義一個方法 (值接收者)
// FullName 方法會回傳 Person 的全名
// (p Person) 是接收者，表示這個方法屬於 Person 型別
func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// 3. 定義一個方法 (指標接收者)
// SetAge 方法會修改 Person 的 Age
// (p *Person) 表示接收者是一個指向 Person 的指標
func (p *Person) SetAge(age int) {
	p.Age = age
}

func main() {
	// --- 建立和使用 Struct ---
	fmt.Println("--- Structs & Methods ---")

	// 建立一個 Person 實例
	p1 := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       40,
	}

	fmt.Printf("P1: %+v\n", p1)

	// --- 呼叫方法 ---
	fmt.Println("\n--- Calling Methods ---")

	// 呼叫值接收者方法
	fmt.Println("Full Name:", p1.FullName())

	// 呼叫指標接收者方法
	fmt.Println("Original Age:", p1.Age)
	p1.SetAge(41)
	fmt.Println("New Age:", p1.Age)

	// --- 指標與值接收者的差異 ---
	fmt.Println("\n--- Pointer vs Value Receiver ---")

	// Go 會自動處理值和指標之間的轉換
	// 即使 p2 是指標，也能直接呼叫值接收者的方法
	p2 := &Person{"Jane", "Doe", 28}
	fmt.Println("P2 Full Name:", p2.FullName())

	// 即使 p1 是值，也能直接呼叫指標接收者的方法
	// Go 會自動轉換為 (&p1).SetAge(42)
	p1.SetAge(42)
	fmt.Println("P1 New Age after implicit conversion:", p1.Age)
}
