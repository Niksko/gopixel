package GoPixel

import "testing"
import "io"
import "io/ioutil"
import "github.com/niksko/GoPixel/edgeDetect"
import "github.com/stretchr/testify/assert"

func TestCannyEdgeDetection(t *testing.T) {
	assert := assert.New(t)

	// How many pixels different we're allowing between the actual and expected before we
	// can't satisfythe assertion
	const DELTA = 10

	// Test a few different images to make sure that their edge detection results
	// match the expected results from another source
	cases := []struct {
		inputFilePath, expectedFilePath            string
		sigmaValue, lowerThreshold, upperThreshold float32
	}{
		{"foo", "bar", 0.1, 0.2, 0.3},
	}
	for _, c := range cases {
		// We wrap each case in a function so that we can make use of defer for cleaning up
		// the open files
		func() {
			// Read the input file
			inFile, err := ioutil.ReadFile(c.inputFilePath)
			defer inFile.Close()
			// Make sure we found it and read it correctly
			if !assert.Nil(err, "was able to oepn the input file") {
				return
			}

			// Read the expected file
			expectedFile, err := ioutil.ReadFile(c.expectedFilePath)
			defer expectedFile.Close()
			// Make sure we found it and read it correctly
			if !assert.Nil(err, "was able to open the expected file") {
				return
			}

			// Decode the input image
			inputImage, _, err := image.Decode(inFile)
			if !assert.Nil(err, "was able to decode the input file") {
				return
			}

			// Decode the expected image
			expectedImage, _, err := image.Decode(expectedFile)
			if !assert.Nil(err, "Was able to decode the expected file") {
				return
			}

			// Run the edge detection on the input file with the specified parameters and get the result
			got := edgeDetect.CannyEdgeDetect(inputImage, c.sigmaValue, c.lowerThreshold, c.upperThreshold)
			// Assert that the image consists of only black or white pixels
			assert.True(t, isBlackAndWhite(got))

			// Assert that tthe image is the same size as the original and the expected
			assert.True(t, sameDimensions(got, inputImage))
			assert.True(t, sameDimensions(got, expectedImage))

			// Run the expected and actual images through a method that compares every pixel in the
			// images, counting a 0 if the pixels are the same and a 1 otherwise, then returns the
			// accumulated result. Then see whether this value is within some delta of zero
			assert.InDelta(t, countDifferentPixels(expectedImage, got), 0, DELTA)
		}()
	}
}

// sameDimensions returns true if the two given images have the same dimensions, false otherwise
func sameDimensions(imageOne, imageTwo image.Image) bool {
	if imageOne.Bounds() == imageTwo.Bounds() {
		returnVal := true
	} else {
		returnVal := false
	}
}
