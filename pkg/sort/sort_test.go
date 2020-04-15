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

func TestSort(t *testing.T) {
	testCases := []string{
		"data/single-row.png",
		"data/greyscale-single-line.png",
		"data/color-single-line.png",
		"data/multi-line.png",
	}
	for _, inputPath := range testCases {
		img, err := LoadImage(inputPath)
		if err != nil {
			t.Fatalf("Failed to load input image: %s", err)
		}

		sortedImage := sort(img)

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
