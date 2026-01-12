package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// confirmDangerousCommand 请求用户确认执行危险命令
// command: 要执行的命令
// description: 危险描述
// riskLevel: 风险等级
// 返回: true 表示用户确认，false 表示用户拒绝
func confirmDangerousCommand(command string, description string, riskLevel string) bool {
	// 显示警告信息
	fmt.Fprintf(os.Stderr, "\n⚠️  检测到潜在危险命令！\n")
	fmt.Fprintf(os.Stderr, "命令: %s\n", command)
	fmt.Fprintf(os.Stderr, "风险: %s (等级: %s)\n\n", description, riskLevel)

	// 请求确认
	fmt.Fprintf(os.Stderr, "是否继续执行？(y/n): ")

	// 读取用户输入
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	// 解析响应
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// confirmWithMessage 通用确认函数
// message: 确认消息
// 返回: true 表示用户确认，false 表示用户拒绝
func confirmWithMessage(message string) bool {
	fmt.Fprintf(os.Stderr, "%s (y/n): ", message)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
