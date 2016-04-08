package gopixel

import "image"
import "math"
import "github.com/disintegration/gift"
import "sort"

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

    // Sort the points in these segments based on the segment angle
    // Wrap the segment angle into the range (-pi, pi]
    segmentAngle = wrapValue(segmentAngle, -math.Pi, math.Pi)

    // Closures that order points
    xAscending := func(c1, c2 *image.Point) bool {
        return c1.X < c2.X
    }
    xDescending := func(c1, c2 *image.Point) bool {
        return c1.X > c2.X
    }
    yAscending := func(c1, c2 *image.Point) bool {
        return c1.Y < c2.Y
    }
    yDescending := func(c1, c2 *image.Point) bool {
        return c1.Y > c2.Y
    }

    // Perform the sort
    for _, segment := range nonEdgeMappedSegments {
        if segmentAngle >= math.Pi/2 {
            OrderedBy(xDescending, yDescending).Sort(segment)
        } else if segmentAngle >= 0 {
            OrderedBy(xAscending, yDescending).Sort(segment)
        } else if segmentAngle >= -math.Pi/2 {
            OrderedBy(xAscending, yAscending).Sort(segment)
        } else {
            OrderedBy(xDescending, yAscending).Sort(segment)
        }
    }

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

func updateEquivalenceMap(eqMap map[int]int, region1, region2 int) {
    // Retrieve the current value from the map
    current1, ok := eqMap[region1]
    // If there currently is a value
    if ok {
        // The new value is the minimum of the old and the region 2 value
        eqMap[region1] = Min(current1, region2)
    } else {
        // Otherwise, the new value is just the region 2 value
        eqMap[region1] = region2
    }

    // Do the same for the other region
    current2, ok := eqMap[region2]
    if ok {
        eqMap[region2] = Min(current2, region1)
    } else {
        eqMap[region2] = region1
    }
}

// This function produces a map from points to integers, where the integers represent contiguous regions of the image
// The 'edges' or 'background' of the image are marked as -1 region
func findConnectedComponents(im image.Image) map[image.Point]int {
    componentMap := make(map[image.Point]int)
    // Count the regions
    regionCounter := 1
    // Need to record which regions are equivalent. Mark this in a map from regions to slices of equivalent regions
    equivalenceMap := make(map[int]int)
    // To find connected components, we iterate over the image, pixel by pixel
    for y := 0; y < im.Bounds().Max.Y; y++ {
        for x := 0; x < im.Bounds().Max.X; x++ {
            // Check to see if the pixel is black
            R, G, B, _ := im.At(x, y).RGBA()
            if R + G + B == 0 {
                northEastRegion := checkRegion(x + 1, y - 1, componentMap)
                northRegion := checkRegion(x, y - 1, componentMap)
                northWestRegion := checkRegion(x - 1, y - 1, componentMap)
                westRegion := checkRegion(x - 1, y, componentMap)
                if northEastRegion != -1 {
                    componentMap[image.Point{x, y}] = northEastRegion
                }
                if northRegion != -1 {
                    // Check the region of the current pixel
                    currentRegion := checkRegion(x, y, componentMap)
                    // Add an entry in the equivalence map between the northRegion and the current region
                    if currentRegion != -1 {
                        updateEquivalenceMap(equivalenceMap, currentRegion, northRegion)
                    }
                    componentMap[image.Point{x, y}] = northRegion
                }
                if northWestRegion != -1 {
                    // Check the region of the current pixel
                    currentRegion := checkRegion(x, y, componentMap)
                    // Add an entry in the equivalence map between the northWestRegion and the current region
                    if currentRegion != -1 {
                        updateEquivalenceMap(equivalenceMap, currentRegion, northWestRegion)
                    }
                    componentMap[image.Point{x, y}] = northWestRegion

                }
                if westRegion != -1 {
                    // Check the region of the current pixel
                    currentRegion := checkRegion(x, y, componentMap)
                    // Add an entry in the equivalence map between the westRegion and the current region
                    if currentRegion != -1 {
                        updateEquivalenceMap(equivalenceMap, currentRegion, westRegion)
                    }
                    componentMap[image.Point{x, y}] = westRegion
                }
                // If none of these have set the current region
                if checkRegion(x, y, componentMap) == -1 {
                    // Set the region to a new region
                    componentMap[image.Point{x, y}] = regionCounter
                    // Increment the region counter
                    regionCounter++
                }
            } else {
                // Mark the region is background with a -1
                componentMap[image.Point{x,y}] = -1
            }
        }
    }

    // Now we perform a final pass to merge equivalent components
    for y := 0; y < im.Bounds().Max.Y; y++ {
        for x := 0; x < im.Bounds().Max.X; x++ {
            region := componentMap[image.Point{x, y}]
            if region != -1 {
                // This should give us the lowest equivalent region
                newRegion := equivalenceMap[region]
                for ;newRegion != region; {
                    region = newRegion
                    newRegion = equivalenceMap[region]
                }
                // Put it back into the component map
                componentMap[image.Point{x, y}] = region
            }
        }
    }

    return componentMap
}

// Convenience function that takes ints, passes them as floats to math.Min, then returns the resulting int
func Min(x, y int) int {
    return int(math.Min(float64(x), float64(y)))
}

func checkRegion(x, y int, componentMap map[image.Point]int) int {
    // If the pixel hasn't been labelled, we want to return -1
    region, ok := componentMap[image.Point{x, y}]
    if !ok {
        region = -1
    }
    return region
}

func divideFromComponentMap(nonEdgeMappedSegments []PixelSegment, pointComponentMap map[image.Point]int) []PixelSegment {
    var returnSlice []PixelSegment
    // Iterate over the segments
    for _, segment := range nonEdgeMappedSegments {
        // Set the current region to something that should never be in the component map
        currentRegion := -2
        // Iterate over the points in the segment
        for _, point := range segment {
            // variable to hold the segment we're writing to
            currentSegment := PixelSegment{}
            // Look up the point in the component map
            region := pointComponentMap[point]
            // If the regions do not match, but we didn't hit edge pixels
            if currentRegion != region && region != -1 {
                // Add the old Segment to the returnSlice if it's not empty
                if len(currentSegment) != 0 {
                    returnSlice = append(returnSlice, currentSegment)
                }
                // Make a new PixelSegment to start adding points to it
                currentSegment = PixelSegment{}
            }
            // Add the point to the currentSegment
            currentSegment = append(currentSegment, point)
            // Set the current region to the new region
            currentRegion = region
        }
    }
    return returnSlice
}

type lessFunc func(p1, p2 *image.Point) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	segment PixelSegment
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(segment PixelSegment) {
	ms.segment = segment
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.segment)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.segment[i], ms.segment[j] = ms.segment[j], ms.segment[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call. We
// could change the functions to return -1, 0, 1 and reduce the
// number of calls for greater efficiency: an exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.segment[i], &ms.segment[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}
