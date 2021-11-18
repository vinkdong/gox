package handler

import (
	"net/http"
	"sort"
	"strings"
)

// 定义处理handler执行器

type Handler struct {
	path     []string
	sorted   bool
	handlers map[string]func(w http.ResponseWriter, r *http.Request) interface{}
}

func New() *Handler {
	return &Handler{
		path:     nil,
		sorted:   false,
		handlers: make(map[string]func(w http.ResponseWriter, r *http.Request) interface{}, 0),
	}
}

func (h *Handler) Register(path string, fn func(w http.ResponseWriter, r *http.Request) interface{}) {
	h.handlers[path] = fn
	h.path = append(h.path, path)
}

func (h *Handler) Do(w http.ResponseWriter, r *http.Request) interface{} {
	if !h.sorted {
		sort.Strings(h.path)
		h.sorted = true
	}
	for x := len(h.path) - 1; x >= 0; x-- {
		path := h.path[x]
		if strings.HasPrefix(r.RequestURI, path) {
			return h.handlers[path](w, r)
		}
	}
	return nil
}
