package store

import (
	"mini-etcd/internal/wal"
	"sync"
)

// Store 接口定义
type Store interface {
	Put(key string, value string) error
	Get(key string) (string, bool)
	Delete(key string) error        // [新增] 删除接口
	Watch(key string) <-chan string // 监听接口
}

// MemoryStore 内存存储实现
type MemoryStore struct {
	mu       sync.RWMutex
	data     map[string]string
	wal      *wal.WAL
	watchers map[string][]chan string
}

func NewMemoryStore(w *wal.WAL) *MemoryStore {
	return &MemoryStore{
		data:     make(map[string]string),
		wal:      w,
		watchers: make(map[string][]chan string),
	}
}

// Put 写入数据
func (s *MemoryStore) Put(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 写 WAL
	entry := wal.LogEntry{
		Command: "PUT",
		Key:     key,
		Value:   value,
	}
	if err := s.wal.Write(entry); err != nil {
		return err
	}

	// 2. 更新内存
	s.data[key] = value

	// 3. 通知 Watcher
	s.notifyWatchers(key, value)

	return nil
}

// Get 读取数据
func (s *MemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}

// Delete 删除数据 [本次新增实现]
func (s *MemoryStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 写 WAL
	entry := wal.LogEntry{
		Command: "DELETE",
		Key:     key,
		Value:   "", // 删除操作不需要 value
	}
	if err := s.wal.Write(entry); err != nil {
		return err
	}

	// 2. 更新内存
	delete(s.data, key)

	// 3. 通知 Watcher (发送空字符串代表删除)
	s.notifyWatchers(key, "")

	return nil
}

// Watch 注册监听
func (s *MemoryStore) Watch(key string) <-chan string {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan string, 1)
	s.watchers[key] = append(s.watchers[key], ch)
	return ch
}

// notifyWatchers 内部辅助方法，用于通知并清理 watcher
func (s *MemoryStore) notifyWatchers(key string, val string) {
	if channels, ok := s.watchers[key]; ok {
		for _, ch := range channels {
			ch <- val
		}
		// Long-Polling 模式：触发一次后清理
		delete(s.watchers, key)
	}
}

// Restore 启动恢复逻辑 [已更新以支持 DELETE]
func (s *MemoryStore) Restore(entries []wal.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range entries {
		switch entry.Command {
		case "PUT":
			s.data[entry.Key] = entry.Value
		case "DELETE":
			delete(s.data, entry.Key)
		}
	}
}
