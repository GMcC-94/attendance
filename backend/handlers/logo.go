package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/types"
)

func UploadLogoHandler(imageStore db.ImageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // max 10MB
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

		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err := os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				http.Error(w, "Could not create uploads directory", http.StatusInternalServerError)
				return
			}
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
		filePath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(filePath)
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
			FileName:   filename,
			FilePath:   filePath,
			Type:       "logo", // save this as a logo
			UploadedAt: time.Now(),
		}

		err = imageStore.SaveImage(image)
		if err != nil {
			http.Error(w, "Could not store image metadata", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"message": "Image uploaded and saved to database",
			"path":    "/uploads/" + filename,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func GetLatestLogoHandler(imageStore db.ImageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image, err := imageStore.GetLatestLogo()
		if err != nil {
			http.Error(w, "No logo found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"fileName": "%s"}`, image.FileName)
	}
}
