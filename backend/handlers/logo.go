package handlers

import (
	"io"
	"net/http"
	"os"
)

func UploadLogoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Falied to parse form", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("logo")
		if err != nil {
			http.Error(w, "No file uploaded", http.StatusBadRequest)
			return
		}

		defer file.Close()

		dst, err := os.Create("uploads/logo.png")
		if err != nil {
			http.Error(w, "Failed to save logo", http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to write logo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		{"message": "Logo uploaded successfully}`))
	}
}
