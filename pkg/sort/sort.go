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

func generatePointOrder(bounds image.Rectangle, sortAngle uint) [][]image.Point {
	var pointOrder [][]image.Point
	if sortAngle == 270 {
		pointOrder = make([][]image.Point, bounds.Dy())
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pointSegment := make([]image.Point, bounds.Dx())
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				pointSegment[x] = image.Pt(x, y)
			}
			pointOrder[y] = pointSegment
		}
	}
	return pointOrder
}

func Sort(inputImage image.Image, sortAngle uint) image.Image {
	// Set up sorted image bounds
	bounds := inputImage.Bounds()
	imageRectangle := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
	sortedImage := image.NewRGBA(imageRectangle)

	pixels := make([][]color.Color, bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		pixelRow := make([]color.Color, bounds.Dx())
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelRow[x] = inputImage.At(x, y)
		}
		s.Sort(ByBrightness(pixelRow))
		pixels[y] = pixelRow
	}

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			sortedImage.Set(x, y, pixels[y][x])
		}
	}

	return sortedImage
}
