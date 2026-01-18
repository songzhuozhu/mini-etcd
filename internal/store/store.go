package store

import (
	"mini-etcd/internal/wal"
	"sync"
)

type Store interface {
	Put(key string, value string) error
	Get(key string) (string, bool)
	// 新增 Watch 接口，返回一个只读通道 (<-chan)
	Watch(key string) <-chan string
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
	wal  *wal.WAL

	// 新增 watchers
	// key -> 一组正在等待的通道
	// 比如 key="foo" 可能有 3 个客户端在同时监听
	watchers map[string][]chan string
}

func NewMemoryStore(w *wal.WAL) *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
		wal:  w,
		// 初始化 map
		watchers: make(map[string][]chan string),
	}
}

// Watch 注册一个监听器
func (s *MemoryStore) Watch(key string) <-chan string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 创建一个带缓冲的通道
	// 缓冲区为 1 是为了防止 Put 操作被阻塞（非常重要的并发细节！）
	ch := make(chan string, 1)

	// 2. 将通道加入到 watchers 列表中
	s.watchers[key] = append(s.watchers[key], ch)

	// 3. 返回通道给调用者（server 层）
	return ch
}

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

	// 3. [新增] 通知监听者 (Notify Watchers)
	// 检查有没有人在监听这个 key
	if channels, ok := s.watchers[key]; ok {
		for _, ch := range channels {
			// 向通道发送最新的 value
			ch <- value
		}
		// 简化版实现：通知完一次后，清空列表（Long Polling 模式）
		// 真实的 Etcd 是流式推送，这里我们做成“一次性触发”
		delete(s.watchers, key)
	}

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
