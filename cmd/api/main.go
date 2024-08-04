package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ocr/internal/utils"
	"os"

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

func handleBasePath(w http.ResponseWriter, r * http.Request) {
	fmt.Fprintf(w, "Hello from the Go web server")
}

func handleImagePath(w http.ResponseWriter, r * http.Request) {
	text, err := ImageToText("assets/test_text.png")
	if err != nil {
		fmt.Fprintf(w, "Something went wrong")
	}

	fmt.Fprintf(w, text)
}

func readAndBase64EncodeImage(filepath string) (string, error) {
	contents, err := os.ReadFile(filepath)

	encoded := base64.StdEncoding.EncodeToString(contents)

	return encoded, err
}

func handleEncodePath(w http.ResponseWriter, r * http.Request) {
	encoded, err := readAndBase64EncodeImage("assets/happy_birthday.jpg")
	if err != nil {
		log.Fatal(err)
	}

	type ResponseData struct {
		Data string `json:"data"`
		Success bool `json:"success"`
	}

	response := ResponseData {
		Data: encoded,
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func main() {
	sum := utils.Add(4, 5)
	fmt.Println(sum)

	http.HandleFunc("/", handleBasePath)
	http.HandleFunc("/image", handleImagePath)
	http.HandleFunc("/encode", handleEncodePath)
	http.ListenAndServe(":8080", nil)
}