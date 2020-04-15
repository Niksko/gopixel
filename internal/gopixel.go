package internal

import (
	"fmt"
	"image"
	"os"
)

func LoadImage(path string) (image.Image, error) {
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
