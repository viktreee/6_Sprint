package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpload(t *testing.T) {
	tests := []struct {
		name, field, filename, input, want string
		status                             int
	}{
		{"text", "myFile", "text.txt", "Привет", ".--. .-. .. .-- . -", http.StatusOK},
		{"Morse", "myFile", "test", ".--. .-. .. .-- . -", "ПРИВЕТ", http.StatusOK},
		{"words", "myFile", "words.txt", "Привет мир", ".--. .-. .. .-- . -   -- .. .-.", http.StatusOK},
		{"empty", "myFile", "empty.txt", "", "", http.StatusInternalServerError},
		{"invalid", "myFile", "bad.txt", ".......", "", http.StatusInternalServerError},
		{"unsupported", "myFile", "bad.txt", "Привет!", "", http.StatusInternalServerError},
		{"wrong field", "file", "text.txt", "Привет", "", http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			var body bytes.Buffer
			form := multipart.NewWriter(&body)
			part, err := form.CreateFormFile(tt.field, tt.filename)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := part.Write([]byte(tt.input)); err != nil {
				t.Fatal(err)
			}
			if err := form.Close(); err != nil {
				t.Fatal(err)
			}
			r := httptest.NewRequest(http.MethodPost, "/upload", &body)
			r.Header.Set("Content-Type", form.FormDataContentType())
			w := httptest.NewRecorder()
			Upload(w, r)
			if w.Code != tt.status {
				t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
			}
			files, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if tt.status != http.StatusOK {
				if len(files) != 0 {
					t.Fatal("failed upload created a result file")
				}
				return
			}
			if w.Body.String() != tt.want {
				t.Errorf("response = %q, want %q", w.Body.String(), tt.want)
			}
			if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
				t.Error("incorrect response content type")
			}
			if len(files) != 1 {
				t.Fatalf("created %d files, want 1", len(files))
			}
			if filepath.Ext(files[0].Name()) != filepath.Ext(tt.filename) {
				t.Error("file extension was not preserved")
			}
			data, err := os.ReadFile(filepath.Join(dir, files[0].Name()))
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tt.want {
				t.Errorf("saved result = %q, want %q", data, tt.want)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	t.Chdir(t.TempDir())
	w := httptest.NewRecorder()
	Index(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("missing HTML: status = %d", w.Code)
	}
	const html = "<html><body>Upload</body></html>"
	if err := os.WriteFile("index.html", []byte(html), 0600); err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	Index(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK || w.Body.String() != html {
		t.Fatalf("HTML response = %d %q", w.Code, w.Body.String())
	}
	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("incorrect HTML content type")
	}
}

func TestRejectedRequests(t *testing.T) {
	t.Chdir(t.TempDir())
	tests := []struct {
		name, method, path string
		handler            http.HandlerFunc
		status             int
	}{
		{"unknown path", http.MethodGet, "/missing", Index, http.StatusNotFound},
		{"index method", http.MethodPost, "/", Index, http.StatusMethodNotAllowed},
		{"upload method", http.MethodGet, "/upload", Upload, http.StatusMethodNotAllowed},
		{"malformed form", http.MethodPost, "/upload", Upload, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.handler(w, httptest.NewRequest(tt.method, tt.path, strings.NewReader("invalid form")))
			if w.Code != tt.status {
				t.Errorf("status = %d, want %d", w.Code, tt.status)
			}
		})
	}
}
