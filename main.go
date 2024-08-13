package main

import (
	"flag"
	_ "image_redactor/docs"
	"image_redactor/http"
	"image_redactor/storage"
	"log"
)

// @title My API Yoo
// @version 1.0
// @description This is a sample server.
// @host localhost:8080
// @BasePath /
func main() {
	addr := flag.String("addr", ":8080", "address for http server")

	s := storage.NewInMemoryTaskStorage()

	log.Printf("Starting server on %s", *addr)
	if err := http.CreateAndRunServer(s, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
