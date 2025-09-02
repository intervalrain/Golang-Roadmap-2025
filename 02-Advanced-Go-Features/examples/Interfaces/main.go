package main

import (
	"fmt"
	"math"
)

// 1. 定義介面
type Shaper interface {
	Area() float64
}

// 2. 實作介面
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// 3. 使用介面作為函式參數
func PrintShapeArea(s Shaper) {
	fmt.Printf("這個形狀的面積是: %0.2f\n", s.Area())
}

// 4. 空介面與型別斷言
func PrintAnything(v interface{}) {
	fmt.Printf("--- 正在處理: %v ---\n", v)
	fmt.Printf("值的型別是: %T\n", v)

	// 型別斷言
	s, ok := v.(string)
	if ok {
		fmt.Println("這是一個 string，內容是:", s)
		return
	}

	i, ok := v.(int)
	if ok {
		fmt.Println("這是一個 int，值是:", i)
		return
	}

	fmt.Println("這是一個未知的型別")
}

func main() {
	fmt.Println("--- Interfaces ---")
	rect := Rectangle{Width: 10, Height: 5}
	circ := Circle{Radius: 3}

	// 因為 Rectangle 和 Circle 都滿足 Shaper 介面，
	// 所以它們都可以被傳遞給 PrintShapeArea 函式。
	PrintShapeArea(rect)
	PrintShapeArea(circ)

	fmt.Println("\n--- Empty Interface & Type Assertion ---")
	PrintAnything(100)
	PrintAnything("Hello, Go!")
	PrintAnything(rect)
}
