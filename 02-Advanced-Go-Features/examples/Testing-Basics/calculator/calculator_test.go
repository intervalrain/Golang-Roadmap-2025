package calculator

import "testing"

// TestAdd 是一個基本的單元測試
func TestAdd(t *testing.T) {
	a := 10
	b := 5
	expected := 15

	result := Add(a, b)

	if result != expected {
		t.Errorf("Add(%d, %d) = %d; 預期為 %d", a, b, result, expected)
	}
}

// TestAddTableDriven 是一個表格驅動測試
func TestAddTableDriven(t *testing.T) {
	// 定義測試案例的表格
	testCases := []struct {
		name     string // 測試案例名稱
		a, b     int    // 輸入
		expected int    // 預期輸出
	}{
		{"正數相加", 2, 3, 5},
		{"負數相加", -2, -3, -5},
		{"與零相加", 7, 0, 7},
		{"結果為負數", 5, -10, -5},
	}

	// 遍歷所有測試案例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Add(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("Add(%d, %d) = %d; 預期為 %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}
