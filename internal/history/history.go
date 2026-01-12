// Package history 提供命令历史记录管理功能
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// HistoryEntry 表示一条历史记录
type HistoryEntry struct {
	// ID 唯一标识符
	ID int `json:"id"`

	// Input 用户的自然语言输入
	Input string `json:"input"`

	// Command 转换后的命令
	Command string `json:"command"`

	// Timestamp 执行时间
	Timestamp time.Time `json:"timestamp"`

	// Success 命令是否执行成功
	Success bool `json:"success"`

	// ExitCode 命令退出码
	ExitCode int `json:"exit_code"`

	// Output 命令输出（截断）
	Output string `json:"output,omitempty"`

	// Error 错误信息
	Error string `json:"error,omitempty"`
}

// History 管理历史记录
type History struct {
	entries  []*HistoryEntry
	nextID   int
	maxSize  int
	mu       sync.RWMutex
	filePath string
}

// NewHistory 创建一个新的 History 实例
func NewHistory() *History {
	return &History{
		entries: make([]*HistoryEntry, 0),
		nextID:  1,
		maxSize: 1000, // 默认保留最近 1000 条
	}
}

// Add 添加一条历史记录
func (h *History) Add(entry *HistoryEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 自动设置 ID
	entry.ID = h.nextID
	h.nextID++

	// 添加到列表
	h.entries = append(h.entries, entry)

	// 限制数量
	if len(h.entries) > h.maxSize {
		// 移除最旧的记录
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}
}

// List 返回所有历史记录（最新的在前）
func (h *History) List() []*HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 创建副本并反转顺序
	result := make([]*HistoryEntry, len(h.entries))
	for i, entry := range h.entries {
		result[len(h.entries)-1-i] = entry
	}

	return result
}

// Get 根据 ID 获取历史记录
func (h *History) Get(id int) (*HistoryEntry, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, entry := range h.entries {
		if entry.ID == id {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("历史记录 ID %d 不存在", id)
}

// Save 保存历史记录到文件
func (h *History) Save(filePath string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 创建父目录
	dir := filePath[:strings.LastIndex(filePath, "/")]
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化历史记录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入历史文件失败: %w", err)
	}

	h.filePath = filePath
	return nil
}

// Load 从文件加载历史记录
func (h *History) Load(filePath string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在不算错误，返回空历史
		return nil
	}

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取历史文件失败: %w", err)
	}

	// 反序列化
	var entries []*HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("解析历史记录失败: %w", err)
	}

	h.entries = entries
	h.filePath = filePath

	// 更新 nextID
	if len(entries) > 0 {
		maxID := 0
		for _, entry := range entries {
			if entry.ID > maxID {
				maxID = entry.ID
			}
		}
		h.nextID = maxID + 1
	}

	return nil
}

// FilterBySuccess 筛选成功/失败的命令
func (h *History) FilterBySuccess(success bool) []*HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]*HistoryEntry, 0)
	for _, entry := range h.entries {
		if entry.Success == success {
			result = append(result, entry)
		}
	}

	return result
}

// Search 搜索历史记录（在 Input 和 Command 中搜索）
func (h *History) Search(query string) []*HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	query = strings.ToLower(query)
	result := make([]*HistoryEntry, 0)

	for _, entry := range h.entries {
		if strings.Contains(strings.ToLower(entry.Input), query) ||
			strings.Contains(strings.ToLower(entry.Command), query) {
			result = append(result, entry)
		}
	}

	return result
}

// Clear 清空所有历史记录
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.entries = make([]*HistoryEntry, 0)
	h.nextID = 1
}

// SetMaxSize 设置最大历史记录数量
func (h *History) SetMaxSize(size int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.maxSize = size

	// 如果当前记录超过新限制，截断
	if len(h.entries) > size {
		h.entries = h.entries[len(h.entries)-size:]
	}
}

// GetFilePath 获取历史文件路径
func (h *History) GetFilePath() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.filePath
}
