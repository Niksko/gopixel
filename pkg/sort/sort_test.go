package sort

import (
	"image"
	"os"
	"testing"

	_ "image/png"
)

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

	sortedImg := sort(img)

	var sortedReader *os.File
	sortedReader, err = os.Open("data/single-row-sorted.png")
	if err != nil {
		t.Errorf("Failed to load result image")
	}
	defer sortedReader.Close()

	var sortedImage image.Image
	sortedImage, _, err = image.Decode(sortedReader)
	if err != nil {
		t.Errorf("Failed to decode result image with error: %s", err)
	}

	if sortedImage != sortedImg {
		t.Errorf("Image didn't match expected image")
	}
}
