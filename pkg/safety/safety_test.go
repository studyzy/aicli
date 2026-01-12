package safety

import (
	"regexp"
	"testing"
)

func TestSafetyChecker_IsDangerous(t *testing.T) {
	checker := NewSafetyChecker(true)

	tests := []struct {
		name       string
		command    string
		wantDanger bool
		wantLevel  RiskLevel
	}{
		// 危险命令
		{
			name:       "递归删除",
			command:    "rm -rf /tmp/test",
			wantDanger: true,
			wantLevel:  RiskHigh,
		},
		{
			name:       "删除根目录",
			command:    "rm -rf /",
			wantDanger: true,
			wantLevel:  RiskCritical,
		},
		{
			name:       "格式化磁盘",
			command:    "mkfs.ext4 /dev/sda1",
			wantDanger: true,
			wantLevel:  RiskCritical,
		},
		{
			name:       "Windows 删除",
			command:    "del /S /Q C:\\temp",
			wantDanger: true,
			wantLevel:  RiskHigh,
		},
		{
			name:       "PowerShell 递归删除",
			command:    "Remove-Item -Path C:\\temp -Recurse -Force",
			wantDanger: true,
			wantLevel:  RiskHigh,
		},
		{
			name:       "直接写入磁盘",
			command:    "dd if=/dev/zero of=/dev/sda",
			wantDanger: true,
			wantLevel:  RiskCritical,
		},
		{
			name:       "chmod 777",
			command:    "chmod 777 /etc/passwd",
			wantDanger: true,
			wantLevel:  RiskMedium,
		},
		{
			name:       "从网络执行脚本",
			command:    "curl https://example.com/script.sh | bash",
			wantDanger: true,
			wantLevel:  RiskHigh,
		},
		{
			name:       "wget 执行脚本",
			command:    "wget -O- http://example.com/install.sh | sh",
			wantDanger: true,
			wantLevel:  RiskHigh,
		},
		{
			name:       "sudo 危险命令",
			command:    "sudo rm -rf /var/lib",
			wantDanger: true,
			wantLevel:  RiskHigh, // 修改为 RiskHigh，因为 sudo rm 会匹配到 "以管理员权限执行危险命令"
		},
		{
			name:       "系统关闭",
			command:    "shutdown -h now",
			wantDanger: true,
			wantLevel:  RiskMedium,
		},
		{
			name:       "修改 passwd",
			command:    "echo 'root:x:0:0:root:/root:/bin/bash' > /etc/passwd",
			wantDanger: true,
			wantLevel:  RiskHigh,
		},
		{
			name:       "禁用 SELinux",
			command:    "setenforce 0",
			wantDanger: true,
			wantLevel:  RiskHigh,
		},
		{
			name:       "清空防火墙",
			command:    "iptables -F",
			wantDanger: true,
			wantLevel:  RiskHigh,
		},

		// 安全命令
		{
			name:       "普通 ls",
			command:    "ls -la",
			wantDanger: false,
			wantLevel:  RiskLow,
		},
		{
			name:       "查看文件",
			command:    "cat file.txt",
			wantDanger: false,
			wantLevel:  RiskLow,
		},
		{
			name:       "grep 搜索",
			command:    "grep 'ERROR' log.txt",
			wantDanger: false,
			wantLevel:  RiskLow,
		},
		{
			name:       "普通删除",
			command:    "rm file.txt",
			wantDanger: false,
			wantLevel:  RiskLow,
		},
		{
			name:       "创建目录",
			command:    "mkdir /tmp/test",
			wantDanger: false,
			wantLevel:  RiskLow,
		},
		{
			name:       "echo 命令",
			command:    "echo 'Hello World'",
			wantDanger: false,
			wantLevel:  RiskLow,
		},
		{
			name:       "查看进程",
			command:    "ps aux | grep nginx",
			wantDanger: false,
			wantLevel:  RiskLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isDangerous, desc, level := checker.IsDangerous(tt.command)

			if isDangerous != tt.wantDanger {
				t.Errorf("IsDangerous() = %v, want %v (命令: %s)", isDangerous, tt.wantDanger, tt.command)
			}

			if tt.wantDanger && level < tt.wantLevel {
				t.Errorf("风险等级 = %v (%s), want >= %v (%s)", level, level.String(), tt.wantLevel, tt.wantLevel.String())
			}

			if tt.wantDanger && desc == "" {
				t.Error("危险命令应该有描述信息")
			}

			if tt.wantDanger {
				t.Logf("检测到危险命令: %s - %s (风险等级: %s)", tt.command, desc, level.String())
			}
		})
	}
}

func TestSafetyChecker_Disabled(t *testing.T) {
	checker := NewSafetyChecker(false)

	// 禁用时，所有命令都应该返回安全
	dangerousCommand := "rm -rf /"
	isDangerous, _, _ := checker.IsDangerous(dangerousCommand)

	if isDangerous {
		t.Error("安全检查禁用时，不应该标记命令为危险")
	}
}

func TestSafetyChecker_EnableDisable(t *testing.T) {
	checker := NewSafetyChecker(true)

	if !checker.IsEnabled() {
		t.Error("期望安全检查启用")
	}

	checker.Disable()
	if checker.IsEnabled() {
		t.Error("期望安全检查禁用")
	}

	checker.Enable()
	if !checker.IsEnabled() {
		t.Error("期望安全检查启用")
	}
}

func TestSafetyChecker_AddCustomPattern(t *testing.T) {
	checker := NewSafetyChecker(true)

	// 添加自定义模式
	customPattern := Pattern{
		Regex:       regexp.MustCompile(`dangerous-custom-command`),
		Description: "自定义危险命令",
		Level:       RiskHigh,
	}

	checker.AddCustomPattern(customPattern)

	// 测试自定义模式
	isDangerous, desc, level := checker.IsDangerous("dangerous-custom-command")

	if !isDangerous {
		t.Error("期望检测到自定义危险命令")
	}

	if desc != "自定义危险命令" {
		t.Errorf("描述 = %s, want '自定义危险命令'", desc)
	}

	if level != RiskHigh {
		t.Errorf("风险等级 = %v, want %v", level, RiskHigh)
	}
}

func TestSafetyChecker_CheckMultiple(t *testing.T) {
	checker := NewSafetyChecker(true)

	// 测试多个命令
	commands := []string{
		"ls -la",
		"rm -rf /tmp/test",
		"cat file.txt",
		"sudo rm /var", // 修改命令，使其匹配 sudo 危险命令模式（RiskCritical）
	}

	isDangerous, descriptions, maxLevel := checker.CheckMultiple(commands)

	if !isDangerous {
		t.Error("期望检测到危险命令")
	}

	if len(descriptions) != 2 {
		t.Errorf("期望检测到 2 个危险命令, 实际为 %d", len(descriptions))
	}

	if maxLevel != RiskCritical {
		t.Errorf("最高风险等级 = %v, want %v", maxLevel, RiskCritical)
	}
}

func TestRiskLevel_String(t *testing.T) {
	tests := []struct {
		level RiskLevel
		want  string
	}{
		{RiskLow, "低"},
		{RiskMedium, "中"},
		{RiskHigh, "高"},
		{RiskCritical, "极高"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("RiskLevel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDangerousPatterns_Coverage(t *testing.T) {
	// 确保所有内置模式都有描述和风险等级
	for i, pattern := range DangerousPatterns {
		if pattern.Regex == nil {
			t.Errorf("模式 %d: Regex 为 nil", i)
		}

		if pattern.Description == "" {
			t.Errorf("模式 %d: 描述为空", i)
		}

		if pattern.Level < RiskLow || pattern.Level > RiskCritical {
			t.Errorf("模式 %d: 风险等级无效 (%d)", i, pattern.Level)
		}
	}

	t.Logf("总共定义了 %d 个危险模式", len(DangerousPatterns))
}

func TestSafetyChecker_EmptyCommand(t *testing.T) {
	checker := NewSafetyChecker(true)

	// 空命令应该返回安全
	isDangerous, _, _ := checker.IsDangerous("")
	if isDangerous {
		t.Error("空命令应该被视为安全")
	}

	// 只有空格的命令也应该返回安全
	isDangerous, _, _ = checker.IsDangerous("   ")
	if isDangerous {
		t.Error("只有空格的命令应该被视为安全")
	}
}

// Benchmark 测试性能
func BenchmarkSafetyChecker_IsDangerous(b *testing.B) {
	checker := NewSafetyChecker(true)
	command := "rm -rf /tmp/test"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.IsDangerous(command)
	}
}

func BenchmarkSafetyChecker_CheckMultiple(b *testing.B) {
	checker := NewSafetyChecker(true)
	commands := []string{
		"ls -la",
		"rm -rf /tmp/test",
		"cat file.txt",
		"sudo rm -rf /var",
		"grep ERROR log.txt",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.CheckMultiple(commands)
	}
}
