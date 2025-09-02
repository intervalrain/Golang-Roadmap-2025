package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// --- if-else ---
	fmt.Println("---")
	score := 85
	if score >= 90 {
		fmt.Println("優等")
	} else if score >= 80 {
		fmt.Println("甲等")
	} else {
		fmt.Println("乙等")
	}

	// if with a short statement
	if n := rand.Intn(10); n%2 == 0 {
		fmt.Printf("%d 是偶數\n", n)
	} else {
		fmt.Printf("%d 是奇數\n", n)
	}

	// --- for loop ---
	fmt.Println("\n---")
	// Basic for loop
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}

	// "while" style
	sum := 1
	sfor sum < 20 {
		sum += sum
	}
	fmt.Println("Sum is", sum)

	// for-range over a slice
	items := []string{"apple", "banana", "cherry"}
	for index, item := range items {
		fmt.Printf("索引 %d: %s\n", index, item)
	}

	// --- switch ---
	fmt.Println("\n---")
	day := "Sunday"
	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Println("工作日")
	case "Saturday", "Sunday":
		fmt.Println("假日")
	default:
		fmt.Println("無效的日期")
	}

	// switch without an expression
	grade := 88
	switch {
	case grade >= 90:
		fmt.Println("優等")
	case grade >= 80:
		fmt.Println("甲等")
	default:
		fmt.Println("乙等")
	}
}
