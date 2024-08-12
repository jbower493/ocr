package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ocr/internal/httpHelpers"
	"ocr/internal/imageProcessing"
	"ocr/internal/textRecognition"
	"strings"
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
	grayscaleImg, grayscaleError := imageProcessing.ConvertToGrayscale(base64String, extension, false)

	if grayscaleError != nil {
		httpHelpers.HandleErrorResponse(w, "Failed to convert image to grayscale", 500)
		return
	}

	// Extract text from image data
	extractedText, extractTextError := textRecognition.Base64BytesToText(grayscaleImg)
	if extractTextError != nil {
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

func main() {
	http.HandleFunc("/image-to-text", handleImageToTextPath)
	http.ListenAndServe(":8080", nil)
}
