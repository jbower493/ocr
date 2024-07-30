package main

import (
	"fmt"
	"log"

	"github.com/otiai10/gosseract/v2"
)

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
		// log.Fatal(err)
		log.Fatalf("Error performing OCR: %v", err)
	}

	fmt.Println(text)

}