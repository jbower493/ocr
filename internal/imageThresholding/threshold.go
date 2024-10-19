package imageThresholding

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func matToImage(mat gocv.Mat) image.Image {
	// Get the Mat's size and type
	width := mat.Cols()
	height := mat.Rows()
	channels := mat.Channels()

	// Create an RGBA image with the same size as the Mat
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate over the Mat and set pixels in the RGBA image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			bgr := mat.GetUCharAt(y, x*channels)
			g := mat.GetUCharAt(y, x*channels+1)
			r := mat.GetUCharAt(y, x*channels+2)
			img.Set(x, y, color.RGBA{R: r, G: g, B: bgr, A: 255})
		}
	}

	return img
}

func imageToMat(img image.Image) gocv.Mat {
	mat := gocv.NewMatWithSize(img.Bounds().Dy(), img.Bounds().Dx(), gocv.MatTypeCV8UC3)
	defer mat.Close()

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			mat.SetUCharAt(y, x*3, uint8(b>>8))
			mat.SetUCharAt(y, x*3+1, uint8(g>>8))
			mat.SetUCharAt(y, x*3+2, uint8(r>>8))
		}
	}

	return mat
}

func Threshold(grayImg image.Image) image.Image {
	mat := imageToMat(grayImg)

	binaryImg := gocv.NewMat()
	defer binaryImg.Close()

	gocv.Threshold(mat, &binaryImg, 127, 255, gocv.ThresholdBinary)

	// convert it back to "image.Image" type and return it
	thresholdedImg := matToImage(binaryImg)
	return thresholdedImg
}
