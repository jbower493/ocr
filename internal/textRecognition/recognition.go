package textRecognition

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"sort"

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

// Contour represents a contour with its bounding box
type Contour struct {
	Rect  image.Rectangle
	Index int
}

// ByY implements sort.Interface for sorting contours by their Y coordinate
type ByY []Contour

func (a ByY) Len() int           { return len(a) }
func (a ByY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByY) Less(i, j int) bool { return a[i].Rect.Min.Y < a[j].Rect.Min.Y }

// ByYX implements sort.Interface for sorting contours by their Y and X coordinates
type ByYX []Contour

func (a ByYX) Len() int      { return len(a) }
func (a ByYX) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByYX) Less(i, j int) bool {
	// Sort by Y coordinate first
	if a[i].Rect.Min.Y == a[j].Rect.Min.Y {
		// If Y coordinates are the same, sort by X coordinate
		return a[i].Rect.Min.X < a[j].Rect.Min.X
	}
	return a[i].Rect.Min.Y < a[j].Rect.Min.Y
}

func GetTextRegions(argImage image.Image) []image.Image {
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

	// Morphological operations to connect text blocks
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(10, 30))
	gocv.Dilate(thresh, &thresh, kernel)

	// Find contours (text regions)
	contours := gocv.FindContours(thresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// Create a slice to hold contours and their bounding boxes
	var contourList []Contour
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		rect := gocv.BoundingRect(contour)

		// Append to contour list
		contourList = append(contourList, Contour{Rect: rect, Index: i})
	}

	// Sort contours by their Y coordinate
	sort.Sort(ByYX(contourList))

	// Create a client for gosseract
	client := gosseract.NewClient()
	defer client.Close()

	var textRegions []image.Image

	for i := 0; i < len(contourList); i++ {
		rect := contourList[i].Rect

		// Draw bounding box (optional, for visualization)
		gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 1)

		// Crop the region of interest (ROI)
		region := img.Region(rect)

		// Convert the cropped region to bytes (e.g., JPEG format)
		// buf := new(bytes.Buffer)
		roiImg, err := region.ToImage()
		if err != nil {
			log.Fatalf("Failed to convert region to image: %v", err)
		}

		// Append text regions to the image.Image array to return from this function
		textRegions = append(textRegions, roiImg)
	}

	// Optionally, save the image with drawn boxes to see the regions
	// gocv.IMWrite("image_with_boxes.png", img)

	return textRegions
}
