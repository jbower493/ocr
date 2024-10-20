package textRecognition

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

func imageToBytes(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	// Encode the image based on the specified format
	switch format {
	case "jpeg":
	case "jpg":
		err = jpeg.Encode(&buf, img, nil)
	case "png":
		err = png.Encode(&buf, img)
	default:
		err = errors.New("unsupported image format")
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ImageToText(img image.Image, imgFormat string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	imageBytes, bytesErr := imageToBytes(img, imgFormat)
	if bytesErr != nil {
		return "", bytesErr
	}

	setImageFromBytesErr := client.SetImageFromBytes(imageBytes)
	if setImageFromBytesErr != nil {
		return "", setImageFromBytesErr
	}

	text, getTextErr := client.Text()
	if getTextErr != nil {
		return "", getTextErr
	}

	return text, nil
}

func GetTextRegions(argImage image.Image) {
	img, matErr := gocv.ImageToMatRGB(argImage)
	if matErr != nil {
		log.Fatalf("Failed to convert image to mat", matErr)
	}

	// Convert the image to grayscale
	gray := gocv.NewMat()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	defer gray.Close()

	// Apply thresholding to isolate text regions
	thresh := gocv.NewMat()
	gocv.Threshold(gray, &thresh, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)
	defer thresh.Close()

	// Find contours (text regions)
	contours := gocv.FindContours(thresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// Create a client for gosseract
	client := gosseract.NewClient()
	defer client.Close()

	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)

		// Get bounding box for each contour
		rect := gocv.BoundingRect(contour)

		// Draw bounding box (optional, for visualization)
		gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)

		// Crop the region of interest (ROI)
		// region := img.Region(rect)

		// Save the cropped region to file or pass it to OCR directly
		regionFile := fmt.Sprintf("column_%d.png", i)
		// gocv.IMWrite(regionFile, region)

		// Perform OCR on the cropped region
		client.SetImage(regionFile)
		text, err := client.Text()
		if err != nil {
			log.Fatalf("Failed to recognize text: %v", err)
		}

		fmt.Printf("Text in column %d: %s\n", i, text)
	}

	// Optionally, save the image with drawn boxes to see the regions
	gocv.IMWrite("image_with_boxes.png", img)
}
