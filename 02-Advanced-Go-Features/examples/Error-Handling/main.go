package main

import (
	"errors"
	"fmt"
)

// --- Custom Error Type ---
// OpError 是一個自訂的錯誤型別，可以包含更多上下文
type OpError struct {
	Op   string
	Code int
	Err  error
}

func (e *OpError) Error() string {
	return fmt.Sprintf("操作: %s, 代碼: %d, 錯誤: %s", e.Op, e.Code, e.Err)
}

// --- Error Wrapping ---
var ErrDataAccess = errors.New("資料存取失敗")

func loadData() error {
	// 模擬一個底層錯誤
	baseErr := ErrDataAccess

	// 將底層錯誤包裝起來，並添加更多上下文
	return &OpError{
		Op:   "loadData",
		Code: 500,
		Err:  baseErr,
	}
}

func main() {
	fmt.Println("--- Error Handling in Go ---")

	err := loadData()

	if err != nil {
		fmt.Println("發生錯誤:", err)

		// --- 使用 errors.Is() 檢查錯誤鏈 ---
		// 檢查錯誤鏈中是否包含 ErrDataAccess
		if errors.Is(err, ErrDataAccess) {
			fmt.Println("日誌: 這是一個資料存取錯誤，需要特別關注！")
		}

		// --- 使用 errors.As() 提取特定型別的錯誤 ---
		var opErr *OpError
		// 檢查錯誤鏈中是否有一個 *OpError 型別的錯誤
		// 如果有，將其賦值給 opErr
		if errors.As(err, &opErr) {
			fmt.Printf("日誌: 捕獲到操作錯誤 - 操作: %s, 代碼: %d\n", opErr.Op, opErr.Code)
		} else {
			fmt.Println("日誌: 這不是一個 *OpError 型別的錯誤")
		}
	}
}
