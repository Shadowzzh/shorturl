package unit

import (
	"testing"

	"short-url/utils"

	"github.com/stretchr/testify/assert"
)

// 单元测试：只测试 GenerateShortID 函数本身
// 不依赖数据库、网络、文件系统等外部资源
func TestGenerateShortID(t *testing.T) {

	t.Run("生成的ID长度应该是6", func(t *testing.T) {
		id := utils.GenerateShortID()
		assert.Len(t, id, 6)
	})

	t.Run("生成的ID只包含合法字符", func(t *testing.T) {
		id := utils.GenerateShortID()
		validChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

		for _, char := range id {
			assert.Contains(t, validChars, string(char))
		}
	})

	t.Run("连续生成的ID应该不同", func(t *testing.T) {
		ids := make(map[string]bool)

		// 生成100个ID，检查是否有重复
		for i := 0; i < 100; i++ {
			id := utils.GenerateShortID()
			assert.False(t, ids[id], "发现重复ID: %s", id)
			ids[id] = true
		}
	})
}
