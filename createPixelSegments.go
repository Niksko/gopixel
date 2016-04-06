package gopixel

import "image"
import "image/color"
import "math"
import "github.com/disintegration/gift"

/*
 * This function takes an image and a segment angle, and creates PixelSegments out of that image. This involves
 * splitting the image up into lines of pixels with the gradient given by the segment angle.
 */
func createSegments(sourceImage image.Image, segmentAngle float64) []PixelSegment {
    var returnSlice []PixelSegment

    // Wrap the segment angle to the range (-pi/2, pi/2]
    segmentAngle = wrapValue(segmentAngle, -math.Pi/2, math.Pi/2)

    // Convert our segment angle into a gradient for our lines
    lineGradient := math.Tan(segmentAngle)

    // If our line gradient is what we get back from a Tan of 90 degrees, set the gradient to NaN
    if lineGradient == math.Tan(math.Pi / 2) {
        lineGradient = math.NaN()
    }

    // If we have a gradient less than 1 (but non-zero), we perform a coordinate and gradient transform that pretends
    // we have gradient greater than 1
    flipCoords := false
    if math.Abs(lineGradient) < 1 && lineGradient != 0 {
        flipCoords = true
        lineGradient = 1. / lineGradient
    }

    // Pull the bounds of the source image into a local variable for convenience
    bounds := sourceImage.Bounds()

    // We will handle vertical and horizontal lines differently
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
                // Add the current point to this segment, taking into account any coordinate transform from flipCoords
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

    // Call the regular function that creates segments with no edge map
    nonEdgeMappedSegments := createSegments(sourceImage, segmentAngle)

    // Use gift to create an edge map of the source image
    convolutionKernel := []float32{
        -1, -1, -1,
        -1, 8, -1,
        -1, -1, -1,
    }

    g := gift.New(
        gift.Grayscale(),
        gift.Convolution(convolutionKernel,
                         false, false, false, 0.0),
        gift.Contrast(100),
    )

    // Create a new empty image to hold the filtered image
    edgeImage := image.NewRGBA(g.Bounds(sourceImage.Bounds()))

    g.Draw(edgeImage, sourceImage)

    // Now we need to use this edgeImage to split up our nonEdgeMappedSegments
    // First we convert this image into a map from points to connected components
    pointComponentMap := findConnectedComponents(edgeImage)

    // Next we use these to divide the nonEdgeMappedSegments
    returnSlice = divideFromComponentMap(nonEdgeMappedSegments, pointComponentMap)

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

func wrapValue(value, lower, upper float64) float64 {
    // Wrap the angle into the range (lower, upper]
    for ;value <= lower || value > upper; {
        if value > 0 {
            value -= upper - lower
        } else {
            value += upper - lower
        }
    }
    return value
}

func findConnectedComponents(im image.Image) map[image.Point]int {
    componentMap := make(map[image.Point]int)
    // To find connected components, we iterate over the image, pixel by pixel
    for x := 0; x < im.Bounds().Max.X; x++ {
        for y := 0; y < im.Bounds().Max.Y; y++ {
            // Check to see if the pixel is black
            var pixel struct {
                R, G, B uint32
            }
            pixel.R, pixel.G, pixel.B, _ = im.At(x, y).RGBA()
            if pixel.R + pixel.G + pixel.B == 0 {

            }
        }
    }
    return componentMap
}

func divideFromComponentMap(nonEdgeMappedSegments []PixelSegment, pointComponentMap map[image.Point]int) []PixelSegment {
    var returnSlice []PixelSegment
    return returnSlice
}
