package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/types"
)

func UploadImageHandler(imageStore db.ImageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Could not parse uploaded file", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Could not read uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
		FilePath := filepath.Join("uploads", filename)

		dst, err := os.Create(FilePath)
		if err != nil {
			http.Error(w, "Could not save file", http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		image := types.Image{
			FileName: filename,
			FilePath: FilePath,
			Type:     r.FormValue("type"),
		}

		err = imageStore.SaveImage(image)
		if err != nil {
			http.Error(w, "Could not store image metadata", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`
		{
		"message": "image uploaded successfully,
		"path":"` + FilePath + `"}`))
	}
}
