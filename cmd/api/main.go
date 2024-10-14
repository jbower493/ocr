package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"ocr/internal/resizeImage"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

func handleOptimizeImagePath(w http.ResponseWriter, r *http.Request) {
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
	extension := strings.Split(mimeType, "/")[1]
	// param3 := r.FormValue("param3") // get extra form fields out of the request

	// Convert to webp
	var img image.Image
	var imageDecodeErr error

	// Get orientation
	var orientationPointer *int

	x, exifDecodeErr := exif.Decode(file)
	if exifDecodeErr == nil {
		orientation, orientationErr := x.Get(exif.Orientation)

		if orientationErr == nil {
			tempOrientationValue, tempOrientationValueErr := orientation.Int(0)

			if tempOrientationValueErr == nil {
				orientationPointer = &tempOrientationValue
			}
		}
	}

	// When decoding a file, it moves the file pointer forwards as data is read. So, because we've already decoded with with exif, we have to put the file pointer back to the start of the file
	file.Seek(0, 0)

	// Decode image
	if mimeType == "image/png" {
		img, imageDecodeErr = png.Decode(file)
	} else if mimeType == "image/jpg" || mimeType == "image/jpeg" {
		img, imageDecodeErr = jpeg.Decode(file)
	} else {
		http.Error(w, extension+" extention not supported", http.StatusBadRequest)
		return
	}

	if imageDecodeErr != nil {
		http.Error(w, "Failed to decode "+extension, http.StatusInternalServerError)
		return
	}

	rotatedImg := img

	// If we couldn't get the exif data just dont do any rotation
	if orientationPointer != nil {
		// Rotate image based on orientation
		switch *orientationPointer {
		// Rotate 180
		case 3:
			rotatedImg = imaging.Rotate180(img)
		// Rotate 90° CW
		case 6:
			rotatedImg = imaging.Rotate270(img)
		// Rotate 90° CCW
		case 8:
			rotatedImg = imaging.Rotate90(img)
		// Normal
		case 1:
		// Flip horizontal
		case 2:
		// Flip vertical
		case 4:
		// Transpose (Rotate 90° CCW + Flip Horizontal)
		case 5:
		// Transverse (Rotate 90° CW + Flip Horizontal)
		case 7:
		// Unknown orientation
		default:
			rotatedImg = img
		}
	}

	// Resize image so that the width is 500px and the height is auto to maintain the aspect ratio
	resizedImg := resizeImage.ResizeImage(rotatedImg)

	var buf bytes.Buffer

	// Setting "quality" to 80% reduces the file size from around 300kb to 50kb!
	webpEncodeErr := webp.Encode(&buf, resizedImg, &webp.Options{Quality: 80})
	if webpEncodeErr != nil {
		http.Error(w, "Failed to encode webp", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, "Unable to write filedata to response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/optimize-image", handleOptimizeImagePath)
	http.ListenAndServe(":8080", nil)
}
