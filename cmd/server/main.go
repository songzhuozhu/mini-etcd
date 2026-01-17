package main

import (
	"log"
	"mini-etcd/internal/server"
	"mini-etcd/internal/store"
	"net/http"
)

func main() {
	// 1. 初始化核心组件 (Dependency Injection)
	// 创建内存存储
	kv := store.NewMemoryStore()
	// 创建 HTTP 服务，并将存储注入进去
	srv := server.NewHTTPServer(kv)

	// 2. 注册路由 (Routing)
	// 将 URL 路径映射到具体的方法
	http.HandleFunc("/put", srv.HandlePut)
	http.HandleFunc("/get", srv.HandleGet)

	// 3. 启动服务
	addr := ":8080"
	log.Printf("Mini-Etcd Server starting on %s ...", addr)

	// ListenAndServe 会一直阻塞，直到出错
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
