package mock

import "net/http"

type HttpWriter struct {
}

func (hw *HttpWriter) Header() http.Header {
	return nil
}

func (hw *HttpWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (hw *HttpWriter) WriteHeader(code int) {
}
