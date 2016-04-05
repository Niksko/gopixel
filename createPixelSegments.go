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

    // We will handle vertical and horizontal lines differently, because they're easy
    if !math.IsNaN(lineGradient) && lineGradient != 0 {

        var startX, startY, endX, endY int

        // For positive gradients, our endpoint is at (0, 0) and our start point is on the line y=height
        if lineGradient > 0 {
            endX, endY = 0, 0
            // Compute start x and y points. Note that the start x needs manual rounding with the +0.5 and int cast
            startY = sourceImage.Bounds().Max.Y - 1
            startX = -1 * int(float64(startY) / lineGradient + 0.5)
        } else if lineGradient < 0 {
            // For negative gradients, the end point is at (0, height) and the start is at y=0
            endX, endY = 0, sourceImage.Bounds().Max.Y - 1
            startX = int(float64(endY) / lineGradient - 0.5)
            startY = 0
        }

        // We need to loop until our start point reaches the correct value
        for ;startX < sourceImage.Bounds().Max.X; {
            var currentSegment PixelSegment

            // We need to set the startPixel of the PixelSegment. To do this we set a flag which will be removed upon
            // The first pixel in bounds
            firstPixel := true

            // Use Bresenham's line algorithm to generate pixels between start and end points
            deltaX := float64(endX - startX)
            deltaY := float64(endY - startY)
            err := 0.
            deltaErr := math.Abs(deltaY / deltaX)
            y := startY
            for x := startX; x <= endX; x++ {
               err += deltaErr
               for ;err >= 0.5; {
                   // If the pixel is within the image bounds
                   if inBounds(sourceImage, x, y) {
                       // Write the pixel to the segment
                       currentSegment.pixelSlice = append(currentSegment.pixelSlice, sourceImage.At(x, y))
                       // If this is the first inbounds pixel, write the point to the start point and remove the flag
                       if firstPixel {
                           firstPixel = false
                           currentSegment.startPoint = image.Point{x, y}
                       }
                   }
                   // Adjust y based on sign of y1 - y0
                   if math.Signbit(float64(endY - startY)) {
                       // Negative sign
                       y -= 1
                   } else {
                       // Positive sign
                       y += 1
                   }
                   // Adjust the error
                   err -= 1
               }
            }

            // Add the segment to the returnSlice
            returnSlice = append(returnSlice, currentSegment)

            // Increment the x start and end points. This should ensure that we cover the entire image
            startX += 1
            endX += 1
        }

    } else if lineGradient == 0 {
        // Horizontal lines
        for y := 0; y < sourceImage.Bounds().Max.Y; y++ {
            var currentSegment PixelSegment
            currentSegment.startPoint = image.Point{0, y}
            for x := 0; x < sourceImage.Bounds().Max.X; x++ {
                currentSegment.pixelSlice = append(currentSegment.pixelSlice, sourceImage.At(x, y))
            }
            returnSlice = append(returnSlice, currentSegment)
        }
    } else {
        // lineGradient is NaN, vertical lines
        for x := 0; x < sourceImage.Bounds().Max.X; x++ {
            var currentSegment PixelSegment
            currentSegment.startPoint = image.Point{x, 0}
            for y := 0; y < sourceImage.Bounds().Max.Y; y++ {
                currentSegment.pixelSlice = append(currentSegment.pixelSlice, sourceImage.At(x, y))
            }
            returnSlice = append(returnSlice, currentSegment)
        }
    }

    return returnSlice
}


// Define a method which computes whether a given x and y coordinate is within the bounds of an image
func inBounds(im image.Image, x, y int) bool {
    inBounds := true
    inBounds = inBounds && x >= im.Bounds().Min.X
    inBounds = inBounds && x < im.Bounds().Max.X
    inBounds = inBounds && y >= im.Bounds().Min.Y
    inBounds = inBounds && y < im.Bounds().Max.Y
    return inBounds
}

func createSegmentsWithEdgeMap(sourceImage, edgeMap image.Image, segmentAngle float64) []PixelSegment {
    var returnSlice []PixelSegment
    return returnSlice
}
