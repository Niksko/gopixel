package GoPixel

import "image"

// CannyEdgeDetect performs Canny edge detection on a given image, using the passed
// smoothingFactor and threshold. The returned image contains white pixels where an
// edge is found, and black pixels elsewhere.
func CannyEdgeDetect(source image.Image, smoothingFactor, lowerThreshold, upperThreshold float32) image.Image {

}
