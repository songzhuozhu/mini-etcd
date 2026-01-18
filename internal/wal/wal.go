package wal

import (
	"encoding/json"
	"os"
	"sync"
)

// LogEntry 代表日志文件中的一行记录
type LogEntry struct {
	Command string `json:"command"` // 操作类型: "PUT" 或 "DELETE"
	Key     string `json:"key"`
	Value   string `json:"value"`
}

// WAL 结构体管理日志文件
type WAL struct {
	file *os.File
	mu   sync.Mutex // 文件写入锁，防止多个 Goroutine 同时写乱文件
}

// NewWAL 创建或打开一个 WAL 文件
func NewWAL(filename string) (*WAL, error) {
	// os.O_APPEND: 追加模式
	// os.O_CREATE: 文件不存在则创建
	// os.O_RDWR: 读写模式
	// 0644: 文件权限 (rw-r--r--)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{file: f}, nil
}

// Write 写入一条日志
func (w *WAL) Write(entry LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 使用 JSON 格式写入，每条记录一行
	// NewEncoder 会自动处理缓冲
	return json.NewEncoder(w.file).Encode(entry)
}

// ReadAll 读取所有日志（用于启动时恢复）
func (w *WAL) ReadAll() ([]LogEntry, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 将文件指针重置到开头
	_, err := w.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var entries []LogEntry
	decoder := json.NewDecoder(w.file)

	// 循环读取直到文件结束 (EOF)
	for decoder.More() {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// Close 关闭文件
func (w *WAL) Close() error {
	return w.file.Close()
}
