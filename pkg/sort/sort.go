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
	if sortAngle == 0 {
		pointOrder = make([][]image.Point, bounds.Dx())
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pointSegment := make([]image.Point, bounds.Dy())
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				pointSegment[y] = image.Pt(x, y)
			}
			pointOrder[x] = pointSegment
		}
	}
	return pointOrder
}

func generateColorSegments(pointOrder [][]image.Point, img image.Image) [][]color.Color {
	result := make([][]color.Color, len(pointOrder))
	for i, points := range pointOrder {
		result[i] = make([]color.Color, len(points))
		for j, point := range points {
			result[i][j] = img.At(point.X, point.Y)
		}
	}
	return result
}

func Sort(inputImage image.Image, sortAngle uint) image.Image {
	// Set up sorted image bounds
	bounds := inputImage.Bounds()
	imageRectangle := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
	sortedImage := image.NewRGBA(imageRectangle)

	pointOrder := generatePointOrder(inputImage.Bounds(), sortAngle)
	colorSegments := generateColorSegments(pointOrder, inputImage)

	for _, segment := range colorSegments {
		s.Sort(ByBrightness(segment))
	}

	for i, points := range pointOrder {
		for j, point := range points {
			sortedImage.Set(point.X, point.Y, colorSegments[i][j])
		}
	}

	return sortedImage
}
