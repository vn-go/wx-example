package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed dist
var content embed.FS

func main() {
	// Tạo một FileServer từ biến `embed.FS`
	// Bằng cách sử dụng fs.Sub, chúng ta chỉ định rằng thư mục gốc là thư mục "dist"
	distFS, err := fs.Sub(content, "dist")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(distFS)))

	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
