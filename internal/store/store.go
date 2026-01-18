package store

import (
	"mini-etcd/internal/wal" // 引入刚才写的 wal 包
	"sync"
)

type Store interface {
	Put(key string, value string) error // 修改返回值，增加 error
	Get(key string) (string, bool)
	// Delete(key string) error // 暂时先不改 Delete，留作练习
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
	wal  *wal.WAL // 持有 WAL 对象的指针
}

// NewMemoryStore 构造函数现在需要传入 WAL
func NewMemoryStore(w *wal.WAL) *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
		wal:  w,
	}
}

// Put 现在包含了两步：写日志 + 更新内存
func (s *MemoryStore) Put(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 先持久化 (Write Ahead)
	entry := wal.LogEntry{
		Command: "PUT",
		Key:     key,
		Value:   value,
	}

	if err := s.wal.Write(entry); err != nil {
		return err // 如果写磁盘失败，直接报错，不更新内存
	}

	// 2. 更新内存
	s.data[key] = value
	return nil
}

func (s *MemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}

// Restore 从日志列表中恢复内存状态
// 这个方法只在启动时调用，不需要加锁（因为那时服务还没开始监听）
func (s *MemoryStore) Restore(entries []wal.LogEntry) {
	for _, entry := range entries {
		if entry.Command == "PUT" {
			s.data[entry.Key] = entry.Value
		}
		// 如果支持 DELETE，这里也要处理
	}
}
