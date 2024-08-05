package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ocr/internal/utils"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

func ImageToText(filename string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImage(filename)
	if err != nil {
		log.Fatalf("Error performing set image on client: %v", err)
	}

	text, err := client.Text()
	if err != nil {
		log.Fatalf("Error performing OCR: %v", err)
	}

	return text, err
}

func handleImageToTextPath(w http.ResponseWriter, r * http.Request) {
	// Read request body
	decoder := json.NewDecoder(r.Body)

	type RequestBody struct {
		Image string `json:"image"`
	}

	var requestBody RequestBody
	decodeErr := decoder.Decode(&requestBody)
	if decodeErr != nil {
		log.Fatal(decodeErr)
		http.Error(w, "Some shit went wrong", http.StatusInternalServerError)
	}

	// Parse out mime type and base64 data from request
	var imageUrlWithoutDataPart string = strings.Split(requestBody.Image, "data:")[1]
	var splitImageUrlWithoutDataPart []string = strings.Split(imageUrlWithoutDataPart, ";base64,")
	var mimeType string = splitImageUrlWithoutDataPart[0]
	var base64String string = splitImageUrlWithoutDataPart[1]

	fmt.Printf("Mime type: %q\n", mimeType)
	fmt.Printf("Base64 string: %q\n", base64String)

	// Extract text from image data

	// Respond
	type ResponseData struct {
		Text string `json:"text"`
		Success bool `json:"success"`
	}

	response := ResponseData {
		Text: "encoded",
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
	sum := utils.Add(4, 5)
	fmt.Println(sum)

	http.HandleFunc("/image-to-text", handleImageToTextPath)
	http.ListenAndServe(":8080", nil)
}