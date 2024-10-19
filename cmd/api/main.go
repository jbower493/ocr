package main

import (
	"net/http"
	optimizeImage "ocr/api/handler"
)

func main() {
	http.HandleFunc("/optimize-image", optimizeImage.Post)
	http.ListenAndServe(":8080", nil)
}
