// Package app 提供应用程序核心功能
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
