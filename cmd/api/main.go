package main

import (
	"net/http"
	"ocr/api/handler/imageToText"
	"ocr/api/handler/optimizeImage"
)

func main() {
	http.HandleFunc("/optimize-image", optimizeImage.Post)
	http.HandleFunc("/image-to-text", imageToText.Post)
	http.ListenAndServe(":8080", nil)
}
