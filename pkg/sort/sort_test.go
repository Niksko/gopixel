package sort

import (
	"image"
	"os"
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
	reader, err := os.Open("data/single-row.png")
	if err != nil {
		t.Errorf("Failed to load test image")
	}
	defer reader.Close()

	var img image.Image
	img, _, err = image.Decode(reader)
	if err != nil {
		t.Errorf("Failed to decode image with error: %s", err)
	}

	sortedImage := sort(img)

	var expectedReader *os.File
	expectedReader, err = os.Open("data/single-row-sorted.png")
	if err != nil {
		t.Errorf("Failed to load result image")
	}
	defer expectedReader.Close()

	var expectedImage image.Image
	expectedImage, _, err = image.Decode(expectedReader)
	if err != nil {
		t.Errorf("Failed to decode result image with error: %s", err)
	}

	if !compareImages(sortedImage, expectedImage) {
		t.Errorf("Image didn't match expected image\nExpected: %v+\nGot     : %v+", expectedImage, sortedImage)
	}
}
