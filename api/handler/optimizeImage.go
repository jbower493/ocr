package optimizeImage

import (
	"bytes"
	"net/http"
	"ocr/internal/decodeImage"
	"ocr/internal/imageOrientation"
	"ocr/internal/resizeImage"
	"ocr/internal/webpImage"
)

func Post(w http.ResponseWriter, r *http.Request) {
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

	// Get orientation
	orientationValue, _ := imageOrientation.GetImageOrientation(file)

	// Decode image
	decodedImage, imageDecodeError := decodeImage.Decode(file, mimeType)
	if imageDecodeError != nil {
		http.Error(w, imageDecodeError.Error(), http.StatusInternalServerError)
		return
	}

	// Rotate
	rotatedImg := imageOrientation.Rotate(decodedImage, orientationValue)

	// Resize
	resizedImg := resizeImage.Resize(rotatedImg)

	// Convert to webp
	var buf bytes.Buffer
	webpEncodeErr := webpImage.ConvertImageToWebp(&buf, resizedImg)
	if webpEncodeErr != nil {
		http.Error(w, "Failed to encode webp", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "image/webp")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, "Unable to write filedata to response", http.StatusInternalServerError)
	}
}
