package imageProcessing

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"
)

func saveImageToFile(filename string, imgData []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(imgData)
	if err != nil {
		return err
	}

	return nil
}

func ConvertToGrayscale(base64String string, saveToFile bool) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64String))
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)

			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	buffer := new(bytes.Buffer)
	pngEncodeErr := png.Encode(buffer, grayImg)
	if pngEncodeErr != nil {
		return nil, pngEncodeErr
	}

	imageByteSlice := buffer.Bytes()

	// Just to verify that the grayscaling worked, optionally save image to file
	if saveToFile {
		saveImageToFile("grayscale.png", imageByteSlice)
	}

	return imageByteSlice, nil
}
