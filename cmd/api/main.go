package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
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
		http.Error(w, "Failed to decode request body", http.StatusInternalServerError)
	}

	// Parse out mime type and base64 data from request
	var imageUrlWithoutDataPart string = strings.Split(requestBody.Image, "data:")[1]
	var splitImageUrlWithoutDataPart []string = strings.Split(imageUrlWithoutDataPart, ";base64,")
	var mimeType string = splitImageUrlWithoutDataPart[0]
	var base64String string = splitImageUrlWithoutDataPart[1]

	// Print mime type
	fmt.Printf("Mime type: %q\n", mimeType)

	/***** TEMP GRAYSCALE START *****/
	// Grayscale
	grayscaleImg, grayscaleError := imageProcessing.ConvertToGrayscale(base64String)
	if grayscaleError != nil {
		http.Error(w, "Failed to convert to grayscale", http.StatusInternalServerError)
	}
	fmt.Println(grayscaleImg)
	/***** TEMP END *****/

	// Extract text from image data
	base64ByteSlice, decodeBase64StringError := base64.StdEncoding.DecodeString(base64String)
	if decodeBase64StringError != nil {
		http.Error(w, "Failed to decode base64 string", http.StatusInternalServerError)
	}

	// TODO: once I have the grayscale byte slice, feed it into here
	extractedText, extractTextError := textRecognition.Base64BytesToText(base64ByteSlice)
	if extractTextError != nil {
		http.Error(w, "Failed to extract text from image", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/image-to-text", handleImageToTextPath)
	http.ListenAndServe(":8080", nil)
}
