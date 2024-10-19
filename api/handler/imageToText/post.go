package imageToText

import (
	"fmt"
	"net/http"
	"ocr/internal/convertImageToGrayscale"
	"ocr/internal/decodeImage"
	"ocr/internal/httpHelpers"
	"ocr/internal/imageOrientation"
	"ocr/internal/textRecognition"
	"strings"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusNotFound)
		return
	}

	// Bitwise operator, shifts the bits of a number to the left, effectively multiplying it by 2. So this is like doing 10 * 2^20, which is aparently equal to 10MB (arg unit is bytes)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file_to_optimize")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	// param3 := r.FormValue("param3") // get extra form fields out of the request

	// Get orientation
	orientationValue, _ := imageOrientation.GetImageOrientation(file)

	// Decode image
	decodedImage, imageDecodeError := decodeImage.Decode(file, mimeType)
	if imageDecodeError != nil {
		http.Error(w, imageDecodeError.Error(), http.StatusInternalServerError)
		return
	}

	// Rotate
	rotatedImg := imageOrientation.Rotate(decodedImage, orientationValue)

	/***********************************************************/

	// Grayscale
	grayscaleImage := convertImageToGrayscale.Convert(rotatedImg)

	// Thresholding
	// thresholdedImage := imageThresholding.Threshold(grayscaleImage)

	extension := strings.Split(mimeType, "/")[1]
	extractedText, _ := textRecognition.ImageToText(grayscaleImage, extension)

	fmt.Println(extractedText)

	// TEMP
	// Create the output file
	// outputFile, err := os.Create("output.webp")
	// if err != nil {
	// 	log.Fatalf("failed to create output file: %v", err)
	// }
	// defer outputFile.Close()
	// err = webp.Encode(outputFile, grayscaleImage, &webp.Options{Lossless: false, Quality: 80})
	// if err != nil {
	// 	log.Fatalf("failed to encode image to webp: %v", err)
	// }

	/***********************************************************/

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
