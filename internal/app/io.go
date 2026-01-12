package app

import (
	"io"
	"os"
)

// hasStdin 检测是否有标准输入（管道或重定向）
// 返回 true 表示有输入数据，false 表示是终端输入
func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	// 检查是否是字符设备（终端）
	// 如果不是字符设备，则是管道或文件重定向
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// readStdin 读取所有标准输入数据
// 返回读取的字符串和可能的错误
func readStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// readStdinIfAvailable 检测并读取标准输入（如果有）
// 返回读取的数据（可能为空字符串）和错误
func readStdinIfAvailable() (string, error) {
	if !hasStdin() {
		return "", nil
	}
	return readStdin()
}

// truncateStdin 截断过长的 stdin 数据
// maxLen: 最大长度（字节）
// 返回截断后的字符串
func truncateStdin(data string, maxLen int) string {
	if len(data) <= maxLen {
		return data
	}

	truncateMsg := "... (truncated)"
	keepLen := maxLen - len(truncateMsg)
	if keepLen < 0 {
		keepLen = 0
	}

	return data[:keepLen] + truncateMsg
}
