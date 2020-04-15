package cmd

import (
	. "github.com/niksko/gopixel/internal"
	gopixel "github.com/niksko/gopixel/pkg/sort"
	log "github.com/sirupsen/logrus"
	"image/png"
	"os"
)

func Sort(filename string, sortAngle uint) (bool, error) {
	image, err := LoadImage(filename)
	if err != nil {
		log.Errorf("Error loading input file: %s", err)
		return false, err
	}

	sorted := gopixel.Sort(image, sortAngle)

	err = png.Encode(os.Stdout, sorted)
	if err != nil {
		log.Errorf("Error encoding image to stdout: %s", err)
		return false, err
	}

	return true, nil
}
