package main

import (
	"encoding/json"
	"fmt"
	"media/controllers/media"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form, max 10MB file
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// Get file from form field "file"
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	md := &media.Media{}
	md.New()
	newFileName, err := md.SaveFile(*fileHeader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("File uploaded successfully: %s", newFileName)))
	}

}
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Handler GET /media/list
	r.Get("/api/media/list", func(w http.ResponseWriter, r *http.Request) {
		md := &media.Media{}
		md.New()
		files, err := md.ListAllFolderAndFiles()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(files); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	r.Post("/api/media/upload", uploadHandler)
	//http.HandleFunc("/api/media/upload", uploadHandler)
	r.Get("/api/media/chart.png", func(w http.ResponseWriter, r *http.Request) {
		md := &media.Media{}
		md.New()
		md.LatencyChartHandler(w)
	})

	http.ListenAndServe(":8081", r)
}
