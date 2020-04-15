package sort

import (
	"image"
	. "math"
)

func abs(a int) float64 {
	return Abs(float64(a))
}

func reverseLine(line []image.Point) []image.Point {
	for i := len(line)/2 - 1; i >= 0; i-- {
		opp := len(line) - 1 - i
		line[i], line[opp] = line[opp], line[i]
	}
	return line
}

func bresenhamLine(x0, y0, x1, y1 int) []image.Point {
	if abs(y1-y0) < abs(x1-x0) {
		if x0 > x1 {
			return reverseLine(bresenhamLineLowAngle(x1, y1, x0, y0))
		} else {
			return bresenhamLineLowAngle(x0, y0, x1, y1)
		}

	} else {
		if y0 > y1 {
			return reverseLine(bresenhamLineHighAngle(x1, y1, x0, y0))
		} else {
			return bresenhamLineHighAngle(x0, y0, x1, y1)
		}
	}
}

func bresenhamLineLowAngle(x0, y0, x1, y1 int) []image.Point {
	var result []image.Point

	deltaX := x1 - x0
	deltaY := y1 - y0
	yIncrement := 1
	if deltaY < 0 {
		yIncrement = -1
		deltaY = -1 * deltaY
	}

	difference := 2*deltaY - deltaX
	y := y0

	for x := x0; x < x1; x++ {
		result = append(result, image.Pt(x, y))
		if difference > 0 {
			y = y + yIncrement
			difference = difference - 2*deltaX
		}
		difference = difference + 2*deltaY
	}

	return result
}

func bresenhamLineHighAngle(x0, y0, x1, y1 int) []image.Point {
	var result []image.Point

	deltaX := x1 - x0
	deltaY := y1 - y0
	xIncrement := 1
	if deltaX < 0 {
		xIncrement = -1
		deltaX = -1 * deltaX
	}

	difference := 2*deltaX - deltaY
	x := x0

	for y := y0; y < y1; y++ {
		result = append(result, image.Pt(x, y))
		if difference > 0 {
			x = x + xIncrement
			difference = difference - 2*deltaY
		}
		difference = difference + 2*deltaX
	}

	return result
}
