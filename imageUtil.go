package main

import (
	"image"
	"image/color"
	"io"

	"github.com/disintegration/imaging"
	"github.com/golang/glog"
	"github.com/nfnt/resize"
)

const SCREEN_WIDTH int = 1404
const SCREEN_HEIGHT int = 1872

func ioReaderToImage(reader io.Reader) (image.Image, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// scale, inset image to reMarkable display size
func adjustImage(img image.Image) image.Image {
	glog.Info("Resizing Image")
	img = resize.Resize(uint(SCREEN_WIDTH), 0, img, resize.Bilinear)

	// Crop image overflow
	img = imaging.Crop(img, image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))

	// Center the image
	background := imaging.New(
		SCREEN_WIDTH,
		SCREEN_HEIGHT,
		color.RGBA{255, 255, 255, 255},
	)
	img = imaging.PasteCenter(background, img)
	return img
}
