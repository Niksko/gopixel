package gopixel

import "image"
import "math"

func createSegments(sourceImage image.Image, segmentAngle float64) []PixelSegment {
    var returnSlice []PixelSegment

    // Convert our segment angle into a gradient for our lines
    lineGradient := math.Tan(segmentAngle)

    // We want to track the furthest right pixel that we visit. That way we can ensure that we assign every pixel in
    // the image to a pixel segment
    rightmostPixelVisited := 0

    var startingY, increment, bound int

    // Set up a starting pixel and an increment, depending on the gradient
    if lineGradient > 0 || math.IsNaN(lineGradient){
        startingY, bound := sourceImage.Bounds().Max.Y, 0
    } else if lineGradient < 0 {
        startingY, bound := 0, sourceImage.Bounds().Max.Y
    }

    // Iterate down the left hand size of the image
    for y := startingY; y != bound; {
        // Set up the segment with starting point
        var segment PixelSegment
        segment.startPoint = image.Point{0, y}
        // Iterate over the pixels in the line determined by the lineGradient

    }
}

func createSegmentsWithEdgeMap(sourceImage image.Image, edgeMap image.Image, segmentAngle float32) []PixelSegment {
}

// This function takes a gradient and starting points and returns a function that generates points along the line 
// defined by those values
func pixelLine(gradient float64, startX int, startY int) func(int) int {
    return func(xValue int) int {
        yValue := 
    }
}
