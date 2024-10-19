package main

import (
	"net/http"
	"ocr/api/handler/optimizeImage"
)

func main() {
	http.HandleFunc("/optimize-image", optimizeImage.Post)
	http.ListenAndServe(":8080", nil)
}
