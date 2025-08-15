package demo

import (
	"fmt"
	"testing"
)

// ✅ 这个会被调用 - 符合 TestXxx(t *testing.T) 规则
func TestWillRun(t *testing.T) {
	fmt.Println("✅ TestWillRun 被调用了!")
}

// ✅ 这个也会被调用
func TestAlsoWillRun(t *testing.T) {
	fmt.Println("✅ TestAlsoWillRun 被调用了!")
}

// ❌ 这个不会被调用 - 小写开头
func testWillNotRun(t *testing.T) {
	fmt.Println("❌ testWillNotRun 不会被调用")
}

// ❌ 这个不会被调用 - 不以 Test 开头
func MyTestFunction(t *testing.T) {
	fmt.Println("❌ MyTestFunction 不会被调用")
}

// ❌ 这个不会被调用 - 改名避免冲突
func NotTestNoParam() {
	fmt.Println("❌ NotTestNoParam 不会被调用")
}

// ✅ 这个是辅助函数，可以被其他测试函数调用
func helperFunction() string {
	return "我是辅助函数"
}

func TestUsingHelper(t *testing.T) {
	result := helperFunction()
	fmt.Printf("✅ TestUsingHelper 调用了: %s\n", result)
}
