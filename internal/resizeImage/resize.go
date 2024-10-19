package resizeImage

import (
	"image"

	"github.com/nfnt/resize"
)

func Resize(img image.Image) image.Image {
	// Resize image so that the width is 500px and the height is auto to maintain the aspect ratio
	newWidth := 500
	newHeight := 0

	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	return resizedImg
}
