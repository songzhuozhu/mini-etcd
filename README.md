# Mini-Etcd

[![Go Version](https://img.shields.io/badge/go-1.20+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[English](#english) | [中文](#chinese)

---

<a id="english"></a>
## English

### Project Overview
Mini-Etcd is a learning project built from scratch to practice Go. It mirrors the core architecture of etcd and implements in-memory storage, write-ahead logging (WAL) persistence, and a Watch mechanism based on Go concurrency.

### Evolution Roadmap (via Git tags)

#### Step 1: Core Data Structure
Tag: `step-1-basic-store`
- Goal: implement a thread-safe in-memory key-value store
- Concepts: interfaces, structs and pointers, `sync.RWMutex`

#### Step 2: HTTP API Layer
Tag: `step-2-http-api`
- Goal: expose the store via REST APIs
- Concepts: `net/http`, `encoding/json`, streaming decoder, dependency injection

#### Step 3: Persistence (WAL)
Tag: `step-3-wal-persistence`
- Goal: survive restarts using WAL
- Concepts: file I/O, log-first pattern, crash recovery, error handling

#### Step 4: Watch Mechanism
Tag: `step-4-watch`
- Goal: allow clients to subscribe to key changes (long-polling)
- Concepts: channels, buffered channels, event-driven notifications

#### Step 5: Timeout Control
Tag: `step-5-timeout-control`
- Goal: avoid infinite Watch blocking
- Concepts: `context.WithTimeout`, `select`, cleanup with `defer cancel()`

### Architecture Mapping: Mini-Etcd vs Real etcd

| Component | Mini-Etcd (Our Code) | Real etcd Source Path | Conceptual Difference |
| :--- | :--- | :--- | :--- |
| Persistence | `internal/wal/wal.go` | `server/wal/` | JSON Lines vs Protobuf + checksums + compression |
| Storage Engine | `internal/store/store.go` | `server/mvcc/` | In-memory map vs BoltDB + MVCC |
| Watcher | `internal/store` | `server/mvcc/watcher.go` | simple list vs watcher groups |
| API Layer | `internal/server/http.go` | `server/etcdserver/api/v3rpc` | REST JSON vs gRPC Protobuf |
| Bootstrap | `cmd/server/main.go` | `server/etcdserver/server.go` | single-node init vs Raft init |
| Consensus | (not implemented) | `raft/` | standalone vs Raft protocol |

### Usage

#### 1) Start the Server
```bash
go run cmd/server/main.go
```

#### 2) API Examples

Put key:
```bash
curl -X POST -d '{"key": "cloud", "value": "native"}' http://localhost:8080/put
```

Get key:
```bash
curl "http://localhost:8080/get?key=cloud"
```

Delete key:
```bash
# Standard DELETE
curl -X DELETE "http://localhost:8080/delete?key=cloud"

# Compatibility mode (e.g., PowerShell)
curl -X POST "http://localhost:8080/delete?key=cloud"
```

Watch key:
```bash
# Terminal A
curl -v "http://localhost:8080/watch?key=status"

# Terminal B
curl -X POST -d '{"key": "status", "value": "ready"}' http://localhost:8080/put
```

---

<a id="chinese"></a>
## 中文

### 项目简介
Mini-Etcd 是一个从零构建的 Go 学习项目，模仿 etcd 的核心架构，包含内存存储、WAL 持久化和基于 Go 并发的 Watch 机制。

### 演进路线（通过 Git Tag 查看）

#### Step 1: 核心数据结构
Tag: `step-1-basic-store`
- 目标：实现线程安全的内存键值存储
- 要点：接口设计、结构体与指针、`sync.RWMutex`

#### Step 2: HTTP 接口层
Tag: `step-2-http-api`
- 目标：通过 REST API 对外提供能力
- 要点：`net/http`、`encoding/json`、流式解码、依赖注入

#### Step 3: 持久化（WAL）
Tag: `step-3-wal-persistence`
- 目标：通过 WAL 保证重启不丢数据
- 要点：文件 IO、先写日志再写内存、启动回放、错误处理

#### Step 4: Watch 机制
Tag: `step-4-watch`
- 目标：支持订阅 Key 变更（长轮询）
- 要点：通道、缓冲通道、事件驱动

#### Step 5: 超时控制
Tag: `step-5-timeout-control`
- 目标：防止 Watch 请求无限阻塞
- 要点：`context.WithTimeout`、`select`、`defer cancel()`

### 架构映射：Mini-Etcd 与真实 etcd

| 组件 | Mini-Etcd（本项目） | etcd 源码路径 | 核心差异 |
| :--- | :--- | :--- | :--- |
| 持久化 | `internal/wal/wal.go` | `server/wal/` | JSON Lines vs Protobuf + 校验 + 压缩 |
| 存储引擎 | `internal/store/store.go` | `server/mvcc/` | 内存 map vs BoltDB + MVCC |
| 监听 | `internal/store` | `server/mvcc/watcher.go` | 简单列表 vs watcher 组 |
| API 层 | `internal/server/http.go` | `server/etcdserver/api/v3rpc` | REST JSON vs gRPC Protobuf |
| 启动流程 | `cmd/server/main.go` | `server/etcdserver/server.go` | 单节点初始化 vs Raft 初始化 |
| 共识 | （未实现） | `raft/` | 单机 vs Raft 协议 |

### 使用指南

#### 1) 启动服务
```bash
go run cmd/server/main.go
```

#### 2) API 示例

写入：
```bash
curl -X POST -d '{"key": "cloud", "value": "native"}' http://localhost:8080/put
```

读取：
```bash
curl "http://localhost:8080/get?key=cloud"
```

删除：
```bash
# 标准 DELETE
curl -X DELETE "http://localhost:8080/delete?key=cloud"

# 兼容模式（如 PowerShell）
curl -X POST "http://localhost:8080/delete?key=cloud"
```

监听：
```bash
# 终端 A
curl -v "http://localhost:8080/watch?key=status"

# 终端 B
curl -X POST -d '{"key": "status", "value": "ready"}' http://localhost:8080/put
```
