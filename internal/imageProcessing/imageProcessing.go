package imageProcessing

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
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

func rotateImage(img image.Image, orientation string) image.Image {
	switch orientation {
	case "1":
		return img
	case "3":
		return imaging.Rotate180(img)
	case "6":
		return imaging.Rotate270(img)
	case "8":
		return imaging.Rotate90(img)
	default:
		return img
	}
}

func ConvertToGrayscale(base64String string, extension string, saveToFile bool) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64String))
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// It seems that only jpeg's have exif data embedded in the image. If it's jpeg, get the orientation and rotate it to it's original orientation
	if extension == "jpeg" {
		imgBytes, _ := base64.StdEncoding.DecodeString(base64String)
		imgReader := bytes.NewReader(imgBytes)
		imgReader.Seek(0, 0)

		exifData, exifErr := exif.Decode(imgReader)
		if exifErr != nil {
			return nil, exifErr
		}

		orientation, exifGetErr := exifData.Get(exif.Orientation)
		if exifGetErr != nil {
			return nil, exifGetErr
		}

		img = rotateImage(img, orientation.String())
	}

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)

			// Experimenting with getting the colour of each pixel, checking if it's black (or close to black), and if not, set it to white (to remove background)
			// r, g, b, a := originalColor.RGBA()
			// if g < 65535 {
			// 	message := fmt.Sprintf("Original colour at %d, %d: %d, %d, %d, %d", x, y, r, g, b, a)
			// 	fmt.Println(message)
			// }

			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	buffer := new(bytes.Buffer)

	if extension == "png" {
		pngEncodeErr := png.Encode(buffer, grayImg)
		if pngEncodeErr != nil {
			return nil, pngEncodeErr
		}
	} else if extension == "jpeg" {
		jpegEncodeErr := jpeg.Encode(buffer, grayImg, nil)
		if jpegEncodeErr != nil {
			return nil, jpegEncodeErr
		}
	}

	imageByteSlice := buffer.Bytes()

	// Just to verify that the grayscaling worked, optionally save image to file
	if saveToFile {
		saveImageToFile("grayscale."+extension, imageByteSlice)
	}

	return imageByteSlice, nil
}
