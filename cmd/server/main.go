package main

import (
	"fmt"
	"log"
	"mini-etcd/internal/server"
	"mini-etcd/internal/store"
	"mini-etcd/internal/wal"
	"net/http"
)

func main() {
	// 1. 初始化 WAL
	// 数据会保存在当前目录下的 server.wal 文件中
	w, err := wal.NewWAL("server.wal")
	if err != nil {
		log.Fatal("无法打开 WAL 文件:", err)
	}
	defer w.Close()

	// 2. 初始化存储并注入 WAL
	kv := store.NewMemoryStore(w)

	// 3. 恢复数据 (Replay)
	fmt.Println("正在从磁盘恢复数据...")
	entries, err := w.ReadAll()
	if err != nil {
		log.Fatal("无法读取 WAL 文件:", err)
	}
	kv.Restore(entries)
	fmt.Printf("成功恢复了 %d 条记录\n", len(entries))

	// 4. 启动 HTTP 服务
	srv := server.NewHTTPServer(kv)
	http.HandleFunc("/put", srv.HandlePut)
	http.HandleFunc("/get", srv.HandleGet)
	http.HandleFunc("/watch", srv.HandleWatch)

	addr := ":8080"
	log.Printf("Mini-Etcd Server starting on %s ...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
