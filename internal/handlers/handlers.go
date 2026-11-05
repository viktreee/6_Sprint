package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("read HTML file: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(data); err != nil {
		log.Printf("write HTML response: %v", err)
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if r.MultipartForm != nil {
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				log.Printf("remove temporary upload files: %v", err)
			}
		}()
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("parse upload form: %v", err), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, fmt.Sprintf("get uploaded file: %v", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close uploaded file: %v", err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("read uploaded file: %v", err), http.StatusInternalServerError)
		return
	}
	result, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("convert uploaded file: %v", err), http.StatusInternalServerError)
		return
	}
	if err := saveResult(result, header.Filename); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := io.WriteString(w, result); err != nil {
		log.Printf("write conversion response: %v", err)
	}
}

func saveResult(result, originalName string) error {
	// Unlike time.Time.String(), this format contains no Windows-invalid colons.
	stamp := time.Now().UTC().Format("20060102T150405.000000000Z")
	name := strings.ReplaceAll(stamp, ".", "_") + filepath.Ext(originalName)
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("create result file: %w", err)
	}
	_, writeErr := file.WriteString(result)
	closeErr := file.Close()
	if writeErr != nil {
		return fmt.Errorf("write result file: %w", writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close result file: %w", closeErr)
	}
	return nil
}
