package webpImage

import (
	"bytes"
	"image"

	"github.com/chai2010/webp"
)

func ConvertImageToWebp(bufferPointer *bytes.Buffer, img image.Image) error {
	// Setting "quality" to 80% reduces the file size from around 300kb to 50kb!
	return webp.Encode(bufferPointer, img, &webp.Options{Quality: 80})
}
