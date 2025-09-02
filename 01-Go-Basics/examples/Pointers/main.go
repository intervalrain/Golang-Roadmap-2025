package main

import "fmt"

// 函式接收一般的值 (pass-by-value)
func addOne(val int) {
	val = val + 1
	fmt.Println("Inside addOne:", val)
}

// 函式接收一個指標 (pass-by-reference)
func addOneWithPointer(val *int) {
	*val = *val + 1
	fmt.Println("Inside addOneWithPointer:", *val)
}

func main() {
	// --- Pointer Basics ---
	fmt.Println("--- Pointer Basics ---")
	x := 10
	p := &x // p 儲存了 x 的記憶體位址

	fmt.Printf("x 的值: %d\n", x)
	fmt.Printf("x 的記憶體位址: %p\n", &x)
	fmt.Printf("p 儲存的位址: %p\n", p)
	fmt.Printf("p 指向的值 (解參考): %d\n", *p)

	// --- Modifying through Pointer ---
	fmt.Println("\n--- Modifying through Pointer ---")
	*p = 20 // 透過指標 p 修改 x 的值
	fmt.Printf("x 現在的值: %d\n", x)

	// --- Pointers in Functions ---
	fmt.Println("\n--- Pointers in Functions ---")
	i := 100
	fmt.Println("Original i:", i)

	addOne(i) // 傳遞 i 的副本
	fmt.Println("After addOne:", i) // i 的值沒有改變

	addOneWithPointer(&i) // 傳遞 i 的記憶體位址
	fmt.Println("After addOneWithPointer:", i) // i 的值被函式修改了

	// --- Nil Pointer ---
	fmt.Println("\n--- Nil Pointer ---")
	var z *int // z 是一個 nil pointer
	fmt.Printf("z 的值: %v\n", z)
	if z == nil {
		fmt.Println("z is a nil pointer.")
	}
	// 對 nil pointer 解參考會導致 panic
	// *z = 1 // 取消這行註解會導致程式崩潰
}
