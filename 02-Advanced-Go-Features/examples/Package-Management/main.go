package main

import (
	"fmt"

	// 為了讓這個範例能夠獨立執行，我們假設 go.mod 的 module 路徑是 "advanced-go/examples"
	// 在實際專案中，您需要將其替換為您在 go.mod 中定義的真實 module 路徑。
	// 例如: "github.com/your-username/your-project/02-Advanced-Go-Features/examples/Package-Management/calculator"
	"golang-Roadmap-2025/02-Advanced-Go-Features/examples/Package-Management/calculator"
)

func main() {
	fmt.Println("--- Package Management ---")

	// 呼叫 calculator 套件中的 Add 函式
	result := calculator.Add(10, 5)
	fmt.Println("10 + 5 =", result)

	// 下面這行會導致編譯錯誤，因為 `subtract` 是小寫開頭，
	// 是 calculator 套件的未導出函式，無法從 main 套件中存取。
	// result = calculator.subtract(10, 5)
}
