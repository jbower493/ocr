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
	"gocv.io/x/gocv"
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

func getCorrectlyRotatedJpeg(img image.Image, base64String string) (image.Image, error) {
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

	switch orientation.String() {
	case "1":
		return img, nil
	case "3":
		return imaging.Rotate180(img), nil
	case "6":
		return imaging.Rotate270(img), nil
	case "8":
		return imaging.Rotate90(img), nil
	default:
		return img, nil
	}
}

func convertImgToGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

func performTresholding(img image.Image) image.Image {
	gocv.NewMat()

	return img
}

func PrepareImageForOcr(base64String string, extension string, saveToFile bool) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64String))
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// It seems that only jpeg's have exif data embedded in the image. If it's jpeg, get the orientation and rotate it to it's original orientation
	if extension == "jpeg" {
		rotatedImg, rotateErr := getCorrectlyRotatedJpeg(img, base64String)
		if rotateErr != nil {
			return nil, rotateErr
		}

		img = rotatedImg
	}

	grayImg := convertImgToGrayscale(img)
	thresholdedImg := performTresholding(grayImg)

	buffer := new(bytes.Buffer)

	if extension == "png" {
		pngEncodeErr := png.Encode(buffer, thresholdedImg)
		if pngEncodeErr != nil {
			return nil, pngEncodeErr
		}
	} else if extension == "jpeg" {
		jpegEncodeErr := jpeg.Encode(buffer, thresholdedImg, nil)
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
