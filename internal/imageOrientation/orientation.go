package imageOrientation

import (
	"image"
	"io"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// Takes in a file and returns the orientation value of the file, according to it's exif data. If orientation can't be determined, returns an error
func GetImageOrientation(file io.Reader) (int, error) {
	x, exifDecodeErr := exif.Decode(file)
	if exifDecodeErr != nil {
		// 0 is the "zero type" for an int, so we have to return that if there is an error
		return 0, exifDecodeErr
	}

	orientation, orientationErr := x.Get(exif.Orientation)
	if orientationErr != nil {
		return 0, orientationErr
	}

	return orientation.Int(0)
}

func Rotate(targetImg image.Image, orientationValue int) image.Image {
	rotatedImage := targetImg

	// Rotate image based on orientation
	switch orientationValue {
	// Rotate 180
	case 3:
		rotatedImage = imaging.Rotate180(targetImg)
	// Rotate 90° CW
	case 6:
		rotatedImage = imaging.Rotate270(targetImg)
	// Rotate 90° CCW
	case 8:
		rotatedImage = imaging.Rotate90(targetImg)
	// Normal
	case 1:
	// Flip horizontal
	case 2:
	// Flip vertical
	case 4:
	// Transpose (Rotate 90° CCW + Flip Horizontal)
	case 5:
	// Transverse (Rotate 90° CW + Flip Horizontal)
	case 7:
	// Unknown orientation
	default:
		rotatedImage = targetImg
	}

	return rotatedImage
}
