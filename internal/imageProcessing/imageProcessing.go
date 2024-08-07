package imageProcessing

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"strings"
)

func ConvertToGrayscale(base64String string) ([]byte, error) {
	fmt.Println("Hello from grayscale")

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

			fmt.Println(originalColor.RGBA())

			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	// encode gray img to byte buffer
	// convert to byte slice and return it

	return nil, nil
}
