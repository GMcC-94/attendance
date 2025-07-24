package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

func UploadLogoHandler(imageStore db.ImageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Invalid multipart form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image uploade failed", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, file); err != nil {
			http.Error(w, "Failed to buffer file", http.StatusInternalServerError)
			return
		}

		url, err := helpers.UploadToS3(file, header)
		if err != nil {
			http.Error(w, "Failed to upload to S3", http.StatusInternalServerError)
			return
		}

		image := types.Image{
			FileName:   filename,
			FileURL:    url,
			Type:       "logo",
			UploadedAt: time.Now(),
		}

		if err := imageStore.SaveImage(image); err != nil {
			http.Error(w, "Failed to save image metadata", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"url": "%s"}`, url)))
	}
}

func GetLogoHandler(imageStore db.ImageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image, err := imageStore.GetImageByType("logo")
		if err != nil {
			http.Error(w, "Logo not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"fileURL": "%s"}`, image.FileURL)))
	}
}
