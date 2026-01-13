package history

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestEntry_Creation 测试历史条目创建
func TestEntry_Creation(t *testing.T) {
	entry := &Entry{
		ID:        1,
		Input:     "list files",
		Command:   "ls -la",
		Timestamp: time.Now(),
		Success:   true,
		ExitCode:  0,
		Output:    "file1.txt\nfile2.txt",
	}

	if entry.ID != 1 {
		t.Errorf("ID = %d, want 1", entry.ID)
	}

	if entry.Input != "list files" {
		t.Errorf("Input = %s, want 'list files'", entry.Input)
	}

	if !entry.Success {
		t.Error("Success should be true")
	}
}

// TestHistory_Add 测试添加历史记录
func TestHistory_Add(t *testing.T) {
	history := NewHistory()

	entry1 := &Entry{
		Input:     "test1",
		Command:   "echo test1",
		Timestamp: time.Now(),
		Success:   true,
	}

	entry2 := &Entry{
		Input:     "test2",
		Command:   "echo test2",
		Timestamp: time.Now(),
		Success:   true,
	}

	history.Add(entry1)
	history.Add(entry2)

	if len(history.entries) != 2 {
		t.Errorf("entries count = %d, want 2", len(history.entries))
	}

	// 验证 ID 自动递增
	if history.entries[0].ID != 1 {
		t.Errorf("first entry ID = %d, want 1", history.entries[0].ID)
	}

	if history.entries[1].ID != 2 {
		t.Errorf("second entry ID = %d, want 2", history.entries[1].ID)
	}
}

// TestHistory_List 测试列出历史记录
func TestHistory_List(t *testing.T) {
	history := NewHistory()

	// 添加多条记录
	for i := 1; i <= 5; i++ {
		history.Add(&Entry{
			Input:   "test",
			Command: "echo test",
			Success: true,
		})
	}

	entries := history.List()
	if len(entries) != 5 {
		t.Errorf("List() returned %d entries, want 5", len(entries))
	}

	// 验证顺序（最新的在前）
	if entries[0].ID != 5 {
		t.Errorf("first entry ID = %d, want 5", entries[0].ID)
	}

	if entries[4].ID != 1 {
		t.Errorf("last entry ID = %d, want 1", entries[4].ID)
	}
}

// TestHistory_Get 测试获取单条历史记录
func TestHistory_Get(t *testing.T) {
	history := NewHistory()

	history.Add(&Entry{
		Input:   "test1",
		Command: "echo test1",
		Success: true,
	})

	history.Add(&Entry{
		Input:   "test2",
		Command: "echo test2",
		Success: true,
	})

	// 获取存在的记录
	entry, err := history.Get(1)
	if err != nil {
		t.Fatalf("Get(1) failed: %v", err)
	}

	if entry.Input != "test1" {
		t.Errorf("entry Input = %s, want 'test1'", entry.Input)
	}

	// 获取不存在的记录
	_, err = history.Get(999)
	if err == nil {
		t.Error("Get(999) should return error")
	}
}

// TestHistory_SaveAndLoad 测试保存和加载历史记录
func TestHistory_SaveAndLoad(t *testing.T) {
	// 创建临时文件
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "test_history.json")

	// 创建并保存历史记录
	history1 := NewHistory()
	history1.Add(&Entry{
		Input:     "test command",
		Command:   "echo test",
		Timestamp: time.Now(),
		Success:   true,
		ExitCode:  0,
		Output:    "test output",
	})

	err := history1.Save(historyFile)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		t.Fatal("History file was not created")
	}

	// 加载历史记录
	history2 := NewHistory()
	loadErr := history2.Load(historyFile)
	if loadErr != nil {
		t.Fatalf("Load() failed: %v", loadErr)
	}

	// 验证数据一致性
	if len(history2.entries) != 1 {
		t.Errorf("loaded entries count = %d, want 1", len(history2.entries))
	}

	entry := history2.entries[0]
	if entry.Input != "test command" {
		t.Errorf("loaded entry Input = %s, want 'test command'", entry.Input)
	}

	if entry.Command != "echo test" {
		t.Errorf("loaded entry Command = %s, want 'echo test'", entry.Command)
	}
}

// TestHistory_LoadNonExistent 测试加载不存在的文件
func TestHistory_LoadNonExistent(t *testing.T) {
	history := NewHistory()
	err := history.Load("/nonexistent/path/history.json")

	// 加载不存在的文件不应该报错（返回空历史）
	if err != nil {
		t.Logf("Load() returned error: %v (expected for non-existent file)", err)
	}

	if len(history.entries) != 0 {
		t.Errorf("entries count = %d, want 0", len(history.entries))
	}
}

// TestHistory_MaxEntries 测试历史记录数量限制
func TestHistory_MaxEntries(t *testing.T) {
	history := NewHistory()

	// 添加超过最大限制的记录
	maxEntries := 1000
	for i := 0; i < maxEntries+100; i++ {
		history.Add(&Entry{
			Input:   "test",
			Command: "echo test",
			Success: true,
		})
	}

	// 验证不超过最大限制
	if len(history.entries) > maxEntries {
		t.Errorf("entries count = %d, should not exceed %d", len(history.entries), maxEntries)
	}
}

// TestHistory_FilterSuccess 测试筛选成功的命令
func TestHistory_FilterSuccess(t *testing.T) {
	history := NewHistory()

	// 添加成功和失败的命令
	history.Add(&Entry{
		Input:   "success1",
		Command: "echo ok",
		Success: true,
	})

	history.Add(&Entry{
		Input:   "failed1",
		Command: "invalid_cmd",
		Success: false,
	})

	history.Add(&Entry{
		Input:   "success2",
		Command: "echo ok2",
		Success: true,
	})

	// 筛选成功的命令
	successEntries := history.FilterBySuccess(true)
	if len(successEntries) != 2 {
		t.Errorf("FilterBySuccess(true) returned %d entries, want 2", len(successEntries))
	}

	// 筛选失败的命令
	failedEntries := history.FilterBySuccess(false)
	if len(failedEntries) != 1 {
		t.Errorf("FilterBySuccess(false) returned %d entries, want 1", len(failedEntries))
	}
}

// TestHistory_Search 测试搜索历史记录
func TestHistory_Search(t *testing.T) {
	history := NewHistory()

	history.Add(&Entry{
		Input:   "list files in directory",
		Command: "ls -la /home",
		Success: true,
	})

	history.Add(&Entry{
		Input:   "find text file",
		Command: "find . -name '*.txt'",
		Success: true,
	})

	history.Add(&Entry{
		Input:   "count lines",
		Command: "wc -l data.txt",
		Success: true,
	})

	// 搜索包含 "file" 的记录（应该匹配前两条）
	results := history.Search("file")
	if len(results) != 2 {
		t.Errorf("Search('file') returned %d entries, want 2", len(results))
	}

	// 验证搜索结果
	for _, entry := range results {
		if !strings.Contains(entry.Input, "file") && !strings.Contains(entry.Command, "file") {
			t.Errorf("Search result should contain 'file': %v", entry)
		}
	}
}

// TestHistory_Clear 测试清空历史记录
func TestHistory_Clear(t *testing.T) {
	history := NewHistory()

	// 添加一些记录
	for i := 0; i < 5; i++ {
		history.Add(&Entry{
			Input:   "test",
			Command: "echo test",
			Success: true,
		})
	}

	if len(history.entries) != 5 {
		t.Fatalf("entries count = %d, want 5", len(history.entries))
	}

	// 清空历史
	history.Clear()

	if len(history.entries) != 0 {
		t.Errorf("entries count after Clear() = %d, want 0", len(history.entries))
	}

	if history.nextID != 1 {
		t.Errorf("nextID after Clear() = %d, want 1", history.nextID)
	}
}

// TestHistory_ConcurrentAccess 测试并发访问
func TestHistory_ConcurrentAccess(t *testing.T) {
	history := NewHistory()

	// 并发添加记录
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			history.Add(&Entry{
				Input:   "concurrent test",
				Command: "echo test",
				Success: true,
			})
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有记录都被添加
	if len(history.entries) != 10 {
		t.Errorf("entries count = %d, want 10", len(history.entries))
	}
}
