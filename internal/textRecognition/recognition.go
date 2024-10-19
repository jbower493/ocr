package textRecognition

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/otiai10/gosseract/v2"
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
