// Package app 提供应用程序核心功能
package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/studyzy/aicli/pkg/i18n"
)

const (
	responseYes = "yes"
)

// confirmDangerousCommand 请求用户确认执行危险命令
// command: 要执行的命令
// description: 危险描述
// riskLevel: 风险等级
// 返回: true 表示用户确认，false 表示用户拒绝
func confirmDangerousCommand(command string, description string, riskLevel string) bool {
	// 显示警告信息
	fmt.Fprintf(os.Stderr, "\n⚠️  %s\n", i18n.T(i18n.WarnDangerousCommand))
	fmt.Fprintf(os.Stderr, "%s: %s\n", i18n.T(i18n.LabelCommand), command)
	fmt.Fprintf(os.Stderr, "%s: %s (%s: %s)\n\n", i18n.T(i18n.WarnRisk), description, i18n.T(i18n.WarnRiskLevel), riskLevel)

	// 请求确认
	msg := i18n.T(i18n.PromptConfirmRisky)
	fmt.Fprintf(os.Stderr, "%s", msg)

	// 读取用户输入
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	// 解析响应
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == responseYes
}
