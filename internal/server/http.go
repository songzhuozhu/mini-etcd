package server

import (
	"context"
	"encoding/json"
	"mini-etcd/internal/store"
	"net/http"
	"time"
)

type HTTPServer struct {
	Store store.Store
}

func NewHTTPServer(s store.Store) *HTTPServer {
	return &HTTPServer{Store: s}
}

// --- DTO ---

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Found bool   `json:"found"`
}

type WatchResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// --- Handlers ---

func (s *HTTPServer) HandlePut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.Store.Put(req.Key, req.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *HTTPServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	val, found := s.Store.Get(key)

	resp := GetResponse{
		Key:   key,
		Value: val,
		Found: found,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleDelete 处理删除请求 [本次新增]
func (s *HTTPServer) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// 按照 RESTful 规范，这里应该用 DELETE 方法，
	// 但为了方便用浏览器或简单工具测试，POST 也可以
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	if err := s.Store.Delete(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted"))
}

func (s *HTTPServer) HandleWatch(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	ch := s.Store.Watch(key)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	select {
	case newVal := <-ch:
		resp := WatchResponse{Key: key, Value: newVal}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	case <-ctx.Done():
		http.Error(w, "Watch timeout", http.StatusRequestTimeout)
	}
}
