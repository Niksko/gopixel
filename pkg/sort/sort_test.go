package sort

import (
	. "github.com/niksko/gopixel/internal"
	"image"
	path "path/filepath"
	"testing"

	_ "image/png"
)

func compareImages(img1, img2 image.Image) bool {
	if img1.Bounds() != img2.Bounds() {
		return false
	}
	for i := img1.Bounds().Min.X; i < img1.Bounds().Max.X; i++ {
		for j := img1.Bounds().Min.Y; j < img1.Bounds().Max.Y; j++ {
			r1, g1, b1, a1 := img1.At(i, j).RGBA()
			r2, g2, b2, a2 := img2.At(i, j).RGBA()
			if r1 != r2 {
				return false
			}
			if g1 != g2 {
				return false
			}
			if b1 != b2 {
				return false
			}
			if a1 != a2 {
				return false
			}
		}
	}
	return true
}

func comparePointOrders(pointOrder1, pointOrder2 [][]image.Point) bool {
	for i, points := range pointOrder1 {
		for j, point1 := range points {
			point2 := pointOrder2[i][j]
			if point1.X != point2.X || point1.Y != point2.Y {
				return false
			}
		}
	}
	return true
}

func TestGeneratePointOrder(t *testing.T) {
	testCases := []struct {
		bounds             image.Rectangle
		sortAngle          uint
		expectedPointOrder [][]image.Point
	}{
		{
			image.Rect(0, 0, 4, 4),
			270,
			[][]image.Point{
				[]image.Point{image.Pt(0, 0), image.Pt(1, 0), image.Pt(2, 0), image.Pt(3, 0)},
				[]image.Point{image.Pt(0, 1), image.Pt(1, 1), image.Pt(2, 1), image.Pt(3, 1)},
				[]image.Point{image.Pt(0, 2), image.Pt(1, 2), image.Pt(2, 2), image.Pt(3, 2)},
				[]image.Point{image.Pt(0, 3), image.Pt(1, 3), image.Pt(2, 3), image.Pt(3, 3)},
			},
		},
	}

	for _, testCase := range testCases {
		result := generatePointOrder(testCase.bounds, testCase.sortAngle)
		if !comparePointOrders(result, testCase.expectedPointOrder) {
			t.Fatalf("Test case didn't match.\nGot:      %v\nExpected: %v\n", result, testCase.expectedPointOrder)
		}
	}
}

func TestSort(t *testing.T) {
	testCases := []struct {
		inputPath string
		sortAngle uint
	}{
		{"data/single-row.png", 270},
		{"data/greyscale-single-line.png", 270},
		{"data/color-single-line.png", 270},
		{"data/multi-line.png", 270},
		{"data/multi-line-vertical.png", 0},
	}
	for _, testCase := range testCases {
		inputPath := testCase.inputPath
		img, err := LoadImage(inputPath)
		if err != nil {
			t.Fatalf("Failed to load input image: %s", err)
		}

		sortedImage := Sort(img, testCase.sortAngle)

		baseFileName := inputPath[:len(inputPath)-len(path.Ext(inputPath))]
		expectedImageFilename := baseFileName + "-sorted" + path.Ext(inputPath)
		var expectedImage image.Image
		expectedImage, err = LoadImage(expectedImageFilename)
		if err != nil {
			t.Fatalf("Failed to load expected image: %s", err)
		}

		if !compareImages(sortedImage, expectedImage) {
			t.Errorf("Image didn't match expected image\nExpected: %v+\nGot     : %v+", expectedImage, sortedImage)
		}
	}
}
