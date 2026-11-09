package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	srv := server.New(logger)
	logger.Printf("server is listening on %s", srv.HTTPServer.Addr)
	if err := srv.HTTPServer.ListenAndServe(); err != nil {
		logger.Fatalf("start server: %v", err)
	}
}
