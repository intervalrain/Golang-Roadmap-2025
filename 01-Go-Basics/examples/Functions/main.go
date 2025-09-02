package main

import (
	"fmt"
)

// 1. 基本函式定義
// 接收兩個 int 參數，回傳一個 int
func add(a int, b int) int {
	return a + b
}

// 如果多個參數型別相同，可以簡寫
func subtract(a, b int) int {
	return a - b
}

// 2. 多重回傳值
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("除數不能為零")
	}
	return a / b, nil
}

// 3. 具名回傳值
func divideWithNamedReturn(a, b float64) (result float64, err error) {
	if b == 0 {
		err = fmt.Errorf("除數不能為零")
		return
	}
	result = a / b
	return
}

// 4. 可變參數函式
func sumAll(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

func main() {
	fmt.Println("--- Basic Functions ---")
	fmt.Println("Add:", add(10, 5))       // 15
	fmt.Println("Subtract:", subtract(10, 5)) // 5

	fmt.Println("\n--- Multiple Return Values ---")
	if result, err := divide(10, 2); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Division result:", result)
	}

	if _, err := divide(10, 0); err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\n--- Named Return Values ---")
	if result, err := divideWithNamedReturn(20, 4); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Division result:", result)
	}

	fmt.Println("\n--- Variadic Functions ---")
	fmt.Println("Sum(1, 2, 3):", sumAll(1, 2, 3))             // 6
	fmt.Println("Sum(10, 20, 30, 40):", sumAll(10, 20, 30, 40)) // 100

	// 將 slice 展開傳入
	nums := []int{5, 6, 7}
	fmt.Println("Sum(slice):", sumAll(nums...)) // 18
}
