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
	if clippedSortAngle == 0 || clippedSortAngle == 180 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var pointSegment []image.Point
			if clippedSortAngle == 0 {
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					pointSegment = append(pointSegment, image.Pt(x, y))
				}
			} else {
				for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
					pointSegment = append(pointSegment, image.Pt(x, y))
				}
			}
			pointOrder = append(pointOrder, pointSegment)
		}
	} else {
		var x0, x1, y0, y1 int
		if clippedSortAngle < 90 {
			x0 = bounds.Max.X - 1
			x1 = bounds.Min.X - 1
			y1 = bounds.Min.Y
			deltaX := x0 - x1
			theta := float64(90-clippedSortAngle) / 180.0 * Pi

			y0 = int(Round(float64(y1) - (float64(deltaX) * Tan(theta))))
		} else if clippedSortAngle == 90 {
			x0 = bounds.Max.X
			x1 = bounds.Min.X
			y1 = bounds.Min.Y - 1
			y0 = bounds.Min.Y - 1
		} else if clippedSortAngle < 180 {
			x0 = bounds.Max.X
			y0 = bounds.Min.Y
			x1 = bounds.Min.X
			deltaX := x0 - x1
			theta := float64(clippedSortAngle-90) / 180.0 * Pi

			y1 = int(Round(float64(y0) - (float64(deltaX) * Tan(theta))))
		} else if clippedSortAngle <= 270 {
			x0 = bounds.Min.X
			y0 = bounds.Min.Y
			x1 = bounds.Max.X
			deltaX := x1 - x0
			theta := float64(270-clippedSortAngle) / 180.0 * Pi

			y1 = int(Round(float64(y0) - (float64(deltaX) * Tan(theta))))
		} else {
			x0 = bounds.Min.X
			x1 = bounds.Max.X
			y1 = bounds.Min.Y
			deltaX := x1 - x0
			theta := float64(clippedSortAngle-270) / 180.0 * Pi

			y0 = int(Round(float64(y1) - (float64(deltaX) * Tan(theta))))
		}

		pointLine := bresenhamLine(x0, y0, x1, y1)

		offset := 0
		nonEmptySegments := false
		for {
			var pointSegment []image.Point
			for _, point := range pointLine {
				offsetPoint := image.Pt(point.X, point.Y+offset)
				if offsetPoint.In(bounds) {
					pointSegment = append(pointSegment, offsetPoint)
				}
			}
			if len(pointSegment) > 0 {
				nonEmptySegments = true
				pointOrder = append(pointOrder, pointSegment)
			} else if nonEmptySegments {
				break
			}
			offset++
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
