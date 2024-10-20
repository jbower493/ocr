package imageToText

import (
	"image"
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

	file, header, err := r.FormFile("image_to_text")
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

	// Grayscale
	grayscaleImage := convertImageToGrayscale.Convert(rotatedImg)

	// Thresholding
	// thresholdedImage := imageThresholding.Threshold(grayscaleImage)

	extension := strings.Split(mimeType, "/")[1]
	imageReadyForOcr := grayscaleImage

	// Do segmentation on the image to get each region of text, then feed each region into the OCR separately
	/**************************************************************************************/

	// Get regions of text in image
	var regions []image.Image
	regions = append(regions, imageReadyForOcr)

	// Loop through regions and extract the text from each one
	var textRegions []string

	for i := 0; i < len(regions); i++ {
		extractedText, extractionError := textRecognition.ImageToText(grayscaleImage, extension)
		if extractionError != nil {
			textRegions = append(textRegions, "")
		} else {
			textRegions = append(textRegions, extractedText)
		}
	}

	/**************************************************************************************/

	// Respond
	type ResponseData struct {
		TextRegions []string `json:"text_regions"`
		Success     bool     `json:"success"`
	}
	response := ResponseData{
		TextRegions: textRegions,
		Success:     true,
	}
	httpHelpers.HandleSuccessResponse(w, response)
}
