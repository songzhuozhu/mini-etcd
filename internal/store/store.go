package store

import (
	"sync"
)

// Store 定义了我们存储引擎的接口
// 使用接口(Interface)是 Go 语言解耦的核心
type Store interface {
	Put(key string, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

// MemoryStore 是 Store 接口的一个内存实现
type MemoryStore struct {
	// mu 是读写互斥锁。
	// Etcd 这种读多写少的场景，RWMutex 比 Mutex 性能更好
	mu sync.RWMutex

	// data 是实际存储数据的 Map
	data map[string]string
}

// NewMemoryStore 是一个构造函数工厂
// Go 没有 class 的构造函数，通常用 New... 函数返回指针
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

// Put 写入数据
func (s *MemoryStore) Put(key string, value string) {
	// 加写锁，阻止其他读写
	s.mu.Lock()
	// defer 关键字非常重要，它确保函数退出前（无论是否报错）一定会执行 Unlock
	// 类似 Python 的 try...finally，但更简洁
	defer s.mu.Unlock()

	s.data[key] = value
}

// Get 读取数据
func (s *MemoryStore) Get(key string) (string, bool) {
	// 加读锁，允许其他 Goroutine 同时读取，但不允许写入
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}

// Delete 删除数据
func (s *MemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
