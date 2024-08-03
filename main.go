package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/otiai10/gosseract/v2"
)

func handler(w http.ResponseWriter, r * http.Request) {
	fmt.Fprintf(w, "Hello jamies world\n")
}

func main() {
	var filename string = "./test_text.png"

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

	fmt.Println(text)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}