package decodeImage

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"strings"
)

func Decode(file multipart.File, mimeType string) (image.Image, error) {
	// When decoding a file, it moves the file pointer forwards as data is read. So, because we've already decoded with with exif, we have to put the file pointer back to the start of the file
	file.Seek(0, 0)

	// Decode image
	if mimeType == "image/png" {
		return png.Decode(file)
	} else if mimeType == "image/jpg" || mimeType == "image/jpeg" {
		return jpeg.Decode(file)
	} else {
		extension := strings.Split(mimeType, "/")[1]
		return nil, errors.New(extension + " extention not supported")
	}
}
