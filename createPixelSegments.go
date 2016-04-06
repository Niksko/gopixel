package gopixel

import "image"
import "math"

/*
 * This function takes an image and a segment angle, and creates PixelSegments out of that image. This involves
 * splitting the image up into lines of pixels with the gradient given by the segment angle.
 */
func createSegments(sourceImage image.Image, segmentAngle float64) []PixelSegment {
    var returnSlice []PixelSegment

    // Wrap the angle into the range (-pi/2, pi/2]
    for ;segmentAngle <= -math.Pi/2 || segmentAngle > math.Pi/2; {
        if segmentAngle > 0 {
            segmentAngle -= math.Pi
        } else {
            segmentAngle += math.Pi
        }
    }

    // Convert our segment angle into a gradient for our lines
    lineGradient := math.Tan(segmentAngle)

    // If our line gradient is larger than that for 90 degrees
    if lineGradient == math.Tan(math.Pi / 2) {
        lineGradient = math.NaN()
    }

    flipCoords := false

    if math.Abs(lineGradient) < 1 && lineGradient != 0 {
        flipCoords = true
        lineGradient = 1. / lineGradient
    }

    // Pull the bounds of the source image into a local variable for convenience
    bounds := sourceImage.Bounds()

    // We will handle vertical lines differently
    if !math.IsNaN(lineGradient)  && lineGradient != 0 {
        // We store our segments in a map, where the keys are the values that are constant over each line
        segmentMap := make(map[float64]PixelSegment)

        // We iterate over each pixel in the image
        for x := 0; x < bounds.Max.X; x++ {
            for y := 0; y < bounds.Max.Y; y++ {
                // Compute the c value for the current pixel
                cValue := float64(y) - lineGradient * float64(x)
                // Round this to the nearest multiple of lineGradient
                roundedC := Round(cValue, lineGradient)
                // Look up this value in the segmentMap
                elem := segmentMap[roundedC]
                // Add the current point to this segment
                if !flipCoords {
                    elem = append(elem, image.Point{x, y})
                } else {
                    elem = append(elem, image.Point{y, x})
                }
                // Store this back in the map
                segmentMap[roundedC] = elem
            }
        }

        // Now iterate over the map, pulling out the values into the returnSlice
        for _, value := range segmentMap {
            returnSlice = append(returnSlice, value)
        }

    } else if math.IsNaN(lineGradient) {
        // lineGradient is NaN, vertical lines
        for x := 0; x < bounds.Max.X; x++ {
            var currentSegment PixelSegment
            for y := 0; y < bounds.Max.Y; y++ {
                currentSegment = append(currentSegment, image.Point{x, y})
            }
            returnSlice = append(returnSlice, currentSegment)
        }
    } else {
        // lineGradient is zero
        for y := 0; y < bounds.Max.Y; y++ {
            var currentSegment PixelSegment
            for x := 0; x < bounds.Max.X; x++ {
                currentSegment = append(currentSegment, image.Point{x, y})
            }
            returnSlice = append(returnSlice, currentSegment)
        }
    }

    return returnSlice
}


func createSegmentsWithEdgeMap(sourceImage, edgeMap image.Image, segmentAngle float64) []PixelSegment {
    var returnSlice []PixelSegment
    return returnSlice
}

// This function takes f, and rounds it to the nearest multiple of m
func Round(f float64, m float64) float64 {
    var returnVal float64
    // Divide by m
    div := f / m
    // Separate the whole and fractional parts
    whole, frac := math.Modf(div)
    if math.Abs(frac) <= 0.5 {
        // If the fractional part is in the range (-0.5, 0.5] we round towards the whole
        returnVal = whole * m
    } else {
        // If the fractional part is outside this range, we round away from the whole
        if math.Signbit(whole) {
            // Negative value
            returnVal = (whole - 1) * m
        } else {
            // Positive value
            returnVal = (whole + 1) * m
        }
    }

    return returnVal
}
