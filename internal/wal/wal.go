package wal

import (
	"encoding/json"
	"os"
	"sync"
)

// LogEntry 代表日志文件中的一行记录
type LogEntry struct {
	Command string `json:"command"` // "PUT" 或 "DELETE"
	Key     string `json:"key"`
	Value   string `json:"value"`
}

// WAL 结构体管理日志文件
type WAL struct {
	file *os.File
	mu   sync.Mutex
}

// NewWAL 创建或打开一个 WAL 文件
func NewWAL(filename string) (*WAL, error) {
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
	return json.NewEncoder(w.file).Encode(entry)
}

// ReadAll 读取所有日志
func (w *WAL) ReadAll() ([]LogEntry, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var entries []LogEntry
	decoder := json.NewDecoder(w.file)

	for decoder.More() {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (w *WAL) Close() error {
	return w.file.Close()
}
