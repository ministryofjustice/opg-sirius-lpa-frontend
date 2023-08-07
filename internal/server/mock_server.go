package server

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

type MockServer struct {
	mux *chi.Mux
	err error
}

func newMockServer(route string, handler Handler) MockServer {
	mux := chi.NewRouter()
	server := MockServer{
		mux: mux,
		err: nil,
	}

	server.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		server.err = handler(w, r)
	})

	return server
}

func (s *MockServer) serve(req *http.Request) (*httptest.ResponseRecorder, error) {
	resp := httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)

	return resp, s.err
}
