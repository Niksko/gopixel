package sort

import (
	"fmt"
	"image"
	"os"
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

func loadImage(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load image at path %s", path)
	}
	defer reader.Close()

	var img image.Image
	img, _, err = image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode image at path %s with error: %s", path, err)
	}

	return img, nil
}

func TestSort(t *testing.T) {
	testCases := []string{
		"data/single-row.png",
	}
	for _, inputPath := range testCases {
		img, err := loadImage(inputPath)
		if err != nil {
			t.Fatalf("Failed to load input image: %s", err)
		}

		sortedImage := sort(img)

		baseFileName := inputPath[:len(inputPath)-len(path.Ext(inputPath))]
		expectedImageFilename := baseFileName + "-sorted" + path.Ext(inputPath)
		var expectedImage image.Image
		expectedImage, err = loadImage(expectedImageFilename)
		if err != nil {
			t.Fatalf("Failed to load expected image: %s", err)
		}

		if !compareImages(sortedImage, expectedImage) {
			t.Errorf("Image didn't match expected image\nExpected: %v+\nGot     : %v+", expectedImage, sortedImage)
		}
	}
}
