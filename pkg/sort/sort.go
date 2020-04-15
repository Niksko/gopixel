package sort

import (
	"image"
	"image/color"
	. "math"
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
	clippedSortAngle := sortAngle % 360
	if clippedSortAngle >= 270 {
		x0 := bounds.Min.X
		x1 := bounds.Max.X
		y1 := bounds.Min.Y
		deltaX := x1 - x0
		theta := float64(clippedSortAngle-270) / 180.0 * Pi

		y0 := int(Round(float64(y1) - (float64(deltaX) * Tan(theta))))

		pointLine := bresenham(x0, y0, x1, y1)

		for offset := 0; offset < bounds.Max.Y-bounds.Min.Y; offset++ {
			var pointSegment []image.Point
			for _, point := range pointLine {
				offsetPoint := image.Pt(point.X, point.Y+offset)
				if offsetPoint.In(bounds) {
					pointSegment = append(pointSegment, offsetPoint)
				}
			}
			pointOrder = append(pointOrder, pointSegment)
		}
	}

	return pointOrder
}

func sgn(a float64) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

// Naive implementation of Bresenham's line algorithm
func bresenham(x0, y0, x1, y1 int) []image.Point {
	result := make([]image.Point, 0)

	deltaX := x1 - x0
	deltaY := y1 - y0
	deltaErr := Abs(float64(deltaY) / float64(deltaX))
	err := 0.0
	y := y0

	for x := x0; x < x1; x++ {
		result = append(result, image.Pt(x, y))
		err = err + deltaErr
		if err >= 0.5 {
			y = y + sgn(float64(deltaY))
			err = err - 1.0
		}
	}
	return result
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
