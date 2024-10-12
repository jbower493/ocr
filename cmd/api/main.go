package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"ocr/internal/resizeImage"

	"github.com/chai2010/webp"
)

func handleOptimizeImagePath(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusNotFound)
		return
	}

	// Bitwise operator, shifts the bits of a number to the left, effectively multiplying it by 2. So this is like doing 10 * 2^20, which is aparently equal to 10MB (arg unit is bytes)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file_to_optimize")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	// param3 := r.FormValue("param3") // get extra form fields out of the request

	// Convert to webp
	var img image.Image
	var imageDecodeErr error

	// PNG
	// TODO: support other file types. Each one we decode, store the img and err in existing vars and just check them once after all if statements
	if mimeType == "image/png" {
		img, imageDecodeErr = png.Decode(file)
	} else if mimeType == "image/jpg" || mimeType == "image/jpeg" {
		img, imageDecodeErr = jpeg.Decode(file)
	} else {
		http.Error(w, "."+strings.Split(mimeType, "/")[1]+" extention not supported", http.StatusBadRequest)
		return
	}

	if imageDecodeErr != nil {
		http.Error(w, "Failed to decode png", http.StatusInternalServerError)
		return
	}

	// Resize image so that the width is 500px and the height is auto to maintain the aspect ratio
	resizedImg := resizeImage.ResizeImage(img)

	var buf bytes.Buffer

	// Setting "quality" to 80% reduces the file size from around 300kb to 50kb!
	webpEncodeErr := webp.Encode(&buf, resizedImg, &webp.Options{Quality: 80})
	if webpEncodeErr != nil {
		http.Error(w, "Failed to encode webp", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, "Unable to write filedata to response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/optimize-image", handleOptimizeImagePath)
	http.ListenAndServe(":8080", nil)
}
