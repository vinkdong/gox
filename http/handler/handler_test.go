package handler

import (
	"fmt"
	"net/http"
	"testing"
	"v2k.io/gox/http/mock"
)

func TestHandler_Do(t *testing.T) {
	handler := New()
	handler.Register("/abc", func(w http.ResponseWriter, r *http.Request) interface{} {
		fmt.Println("/abc")
		return nil
	})
	handler.Register("/abcd", func(w http.ResponseWriter, r *http.Request) interface{} {
		fmt.Println("/abcd")
		return nil
	})
	handler.Register("/abc/d", func(w http.ResponseWriter, r *http.Request) interface{} {
		fmt.Println("/abc/d")
		return nil
	})
	handler.Register("/abc/dc", func(w http.ResponseWriter, r *http.Request) interface{} {
		fmt.Println("/abc/dc")
		return nil
	})
	handler.Register("/abc/d/e/", func(w http.ResponseWriter, r *http.Request) interface{} {
		fmt.Println("/abc/d/e/")
		return nil
	})
	handler.Register("/abc/d/ef", func(w http.ResponseWriter, r *http.Request) interface{} {
		fmt.Println("/abc/d/ef")
		return nil
	})
	handler.Register("/xyz/d/ef", func(w http.ResponseWriter, r *http.Request) interface{} {
		fmt.Println("/xyz/d/ef")
		return nil
	})

	w := &mock.HttpWriter{}
	r := &http.Request{
		Method:     "GET",
		RequestURI: "/xyz/d/ef/s2",
	}
	handler.Do(w, r)

	r = &http.Request{
		Method:     "GET",
		RequestURI: "/abc/d/exx",
	}
	handler.Do(w, r)
}
