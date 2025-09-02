package calculator

// Add 函式因為是大寫開頭，所以可以被其他套件呼叫 (Exported)。
func Add(a, b int) int {
	return a + b
}

// subtract 函式是小寫開頭，只能在 calculator 套件內部使用 (Unexported)。
func subtract(a, b int) int {
	return a - b
}
