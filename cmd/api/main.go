package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"ocr/internal/httpHelpers"
	"ocr/internal/imageProcessing"
	"ocr/internal/textRecognition"
	"strings"

	"github.com/chai2010/webp"
)

func handleImageToTextPath(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		fmt.Fprintf(w, "Hello options")
		return
	}

	// Read request body
	decoder := json.NewDecoder(r.Body)

	type RequestBody struct {
		Image string `json:"image"`
	}

	var requestBody RequestBody
	decodeErr := decoder.Decode(&requestBody)
	if decodeErr != nil {
		httpHelpers.HandleErrorResponse(w, "Failed to decode request body", 500)
		return
	}

	// Parse out mime type and base64 data from request
	var imageUrlWithoutDataPart string = strings.Split(requestBody.Image, "data:")[1]
	var splitImageUrlWithoutDataPart []string = strings.Split(imageUrlWithoutDataPart, ";base64,")
	var mimeType string = splitImageUrlWithoutDataPart[0]
	var extension string = strings.Split(mimeType, "/")[1]
	var base64String string = splitImageUrlWithoutDataPart[1]

	// Grayscale
	preparedImg, prepareErr := imageProcessing.PrepareImageForOcr(base64String, extension, true)

	if prepareErr != nil {
		fmt.Println(prepareErr)
		httpHelpers.HandleErrorResponse(w, "Failed to prepare image for OCR", 500)
		return
	}

	// Extract text from image data
	extractedText, extractTextError := textRecognition.Base64BytesToText(preparedImg)
	if extractTextError != nil {
		fmt.Println(extractTextError)
		httpHelpers.HandleErrorResponse(w, "Failed to extract text from image", 500)
		return
	}

	// Respond
	type ResponseData struct {
		Text    string `json:"text"`
		Success bool   `json:"success"`
	}

	response := ResponseData{
		Text:    extractedText,
		Success: true,
	}

	httpHelpers.HandleSuccessResponse(w, response)
}

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

	// Decode image
	var img image.Image
	var imageDecodeErr error

	// PNG
	if mimeType == "image/png" {
		img, imageDecodeErr = png.Decode(file)
	} else {
		http.Error(w, "Extention not supported", http.StatusBadRequest)
		return
	}

	if imageDecodeErr != nil {
		http.Error(w, "Failed to decode png", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer

	webpEncodeErr := webp.Encode(&buf, img, &webp.Options{Lossless: true})
	if webpEncodeErr != nil {
		http.Error(w, "Failed to encode webp", http.StatusInternalServerError)
		return
	}

	// if err = ioutil.WriteFile("output.webp", buf.Bytes(), 0666); err != nil {
	// 	http.Error(w, "Failed to write webp file", http.StatusInternalServerError)
	// 	return
	// }

	// m, err := webp.Decode(bytes.NewReader(data))
	// if err != nil {
	// 	log.Println(err)
	// }

	// Read file binary data
	// fileData, err := io.ReadAll(file)
	// if err != nil {
	// 	http.Error(w, "Unable to read file data", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "image/webp")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, "Unable to write filedata to response", http.StatusInternalServerError)
	}

	// fmt.Fprintf(w, "Successfully failed lol")
}

func main() {
	http.HandleFunc("/image-to-text", handleImageToTextPath)
	http.HandleFunc("/optimize-image", handleOptimizeImagePath)
	http.ListenAndServe(":8080", nil)
}
