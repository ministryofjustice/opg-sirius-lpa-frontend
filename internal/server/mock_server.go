package server

import (
	"net/http"
	"net/http/httptest"
)

type MockServer struct {
	mux *http.ServeMux
	err error
}

func newMockServer(route string, handler Handler) *MockServer {
	mux := http.NewServeMux()
	server := MockServer{
		mux: mux,
		err: nil,
	}

	server.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		server.err = handler(w, r)
	})

	return &server
}

func (s *MockServer) serve(req *http.Request) (*httptest.ResponseRecorder, error) {
	resp := httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)

	return resp, s.err
}
