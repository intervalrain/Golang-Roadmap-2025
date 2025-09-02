package main

import "fmt"

func declare1() {
	var age = 30
	var name = "Alice"
	var score = 99.5
	fmt.Println(age)
	fmt.Println(name)
	fmt.Println(score)
}

func declare2() {
	age := 33
	name := "Bob"
	score := 97
	fmt.Println(age)
	fmt.Println(name)
	fmt.Println(score)
}

func main() {
	declare1()
	declare2()
}