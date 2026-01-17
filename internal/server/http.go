package server

import (
	"encoding/json"
	"mini-etcd/internal/store"
	"net/http"
)

// HTTPServer 结构体，持有 Store 接口
// 这里体现了 Go 的组合特性，server 依赖 store
type HTTPServer struct {
	Store store.Store
}

// NewHTTPServer 构造函数
func NewHTTPServer(s store.Store) *HTTPServer {
	return &HTTPServer{Store: s}
}

// --- 数据传输对象 (DTO) ---

// PutRequest 用于解析用户发来的 JSON Body
// 注意：字段首字母必须大写，否则 json 包无法通过反射访问！
// `json:"key"` 是 Struct Tag，告诉 Go 解析器对应 JSON 中的哪个字段
type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetResponse 用于返回 JSON 响应
type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Found bool   `json:"found"`
}

// --- Handlers (处理逻辑) ---

// HandlePut 处理写入请求
func (s *HTTPServer) HandlePut(w http.ResponseWriter, r *http.Request) {
	// 1. 校验请求方法
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. 解析 JSON Body
	var req PutRequest
	// 使用 NewDecoder 流式解析，比 ioutil.ReadAll 性能更好
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 3. 调用核心存储逻辑
	s.Store.Put(req.Key, req.Value)

	// 4. 返回结果
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// HandleGet 处理读取请求
func (s *HTTPServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. 获取 URL 参数
	key := r.URL.Query().Get("key")

	// 2. 查询存储
	value, found := s.Store.Get(key)

	// 3. 构建响应对象
	resp := GetResponse{
		Key:   key,
		Value: value,
		Found: found,
	}

	// 4. 设置 Header 并编码返回
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
