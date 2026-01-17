package main

import (
	"fmt"
	"mini-etcd/internal/store" // 导入我们要刚才写的包
)

func main() {
	// 1. 初始化存储
	kv := store.NewMemoryStore()

	// 2. 写入数据
	fmt.Println("正在写入数据: key=name, value=etcd-learner")
	kv.Put("name", "etcd-learner")

	// 3. 读取数据
	// Go 支持多返回值，这里返回 value 和 是否存在的布尔值
	val, exists := kv.Get("name")
	if exists {
		fmt.Printf("读取成功: key=name, value=%s\n", val)
	} else {
		fmt.Println("key 不存在")
	}

	// 4. 读取不存在的数据
	_, exists = kv.Get("unknown")
	if !exists {
		fmt.Println("读取 unknown 失败，符合预期")
	}
}
