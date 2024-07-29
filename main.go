package main

import (
	"github.com/otiai10/gosseract"
)

func main() {
	// var filename string = "./test_text.png"

	client := gosseract.NewClient()
	defer client.Close()

	// client.SetImage(filename)
	// text, err := client.Text()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(text)
}