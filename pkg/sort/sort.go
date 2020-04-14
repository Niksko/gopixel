package sort

import (
	"image"
	"image/color"
	s "sort"
)

func brightness(c color.Color) uint32 {
	r, g, b, _ := c.RGBA()
	return r + g + b
}

type ByBrightness []color.Color

func (c ByBrightness) Len() int {
	return len(c)
}

func (c ByBrightness) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByBrightness) Less(i, j int) bool {
	return brightness(c[i]) > brightness(c[j])
}

func sort(inputImage image.Image) image.Image {
	// Set up sorted image bounds
	bounds := inputImage.Bounds()
	imageRectangle := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
	sortedImage := image.NewRGBA(imageRectangle)

	// Sort image, brightest pixels on the left
	pixels := make([]color.Color, bounds.Dx())
	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		pixels[i] = inputImage.At(i, bounds.Min.Y)
	}

	s.Sort(ByBrightness(pixels))

	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		sortedImage.Set(i, bounds.Min.Y, pixels[i])
	}

	return sortedImage
}
