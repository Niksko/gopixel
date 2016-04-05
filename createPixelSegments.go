package gopixel

import "image"
import "math"

/*
 * This function takes an image and a segment angle, and creates PixelSegments out of that image. This involves
 * splitting the image up into lines of pixels with the gradient given by the segment angle.
 */
func createSegments(sourceImage image.Image, segmentAngle float64) []PixelSegment {
    var returnSlice []PixelSegment

    // Convert our segment angle into a gradient for our lines
    lineGradient := math.Tan(segmentAngle)

    // For positive gradients, start in the top left
    if lineGradient >= 0 {
        endPoint := image.Point{0, 0}
        if lineGradient >= 1 {
            endPoint.X += 1
            endPoint.Y -= int(lineGradient+0.5) // Rounding, since apparently not in the standard library
        } else {
            endPoint.X += int:wq(1 / lineGradient + 0.5)
        }
    } else {
    }

    return returnSlice
}

func createSegmentsWithEdgeMap(sourceImage image.Image, edgeMap image.Image, segmentAngle float64) []PixelSegment {
    var returnSlice []PixelSegment
    return returnSlice
}
