package handlers

import (
	"fmt"
	"net/http"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
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

		url, err := helpers.UploadToS3(file, header)
		if err != nil {
			http.Error(w, "Failed to upload to S3: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"url": "%s"}`, url)))
	}
}
