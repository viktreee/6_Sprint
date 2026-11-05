package server

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	logger := log.New(io.Discard, "", 0)
	srv := New(logger)
	if srv.Logger != logger || srv.HTTPServer.ErrorLog != logger {
		t.Error("logger was not passed to the server")
	}
	if srv.HTTPServer.Addr != ":8080" {
		t.Errorf("address = %q", srv.HTTPServer.Addr)
	}
	if srv.HTTPServer.ReadTimeout != 5*time.Second || srv.HTTPServer.WriteTimeout != 10*time.Second || srv.HTTPServer.IdleTimeout != 15*time.Second {
		t.Error("incorrect timeouts")
	}
	for _, path := range []string{"/", "/upload"} {
		w := httptest.NewRecorder()
		srv.HTTPServer.Handler.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, path, nil))
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("route %q: status = %d", path, w.Code)
		}
	}
}
