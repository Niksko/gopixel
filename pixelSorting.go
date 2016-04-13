package gopixel

import "github.com/disintegration/gift"
import "image"
import "image/draw"
import "image/color"
import "sort"

type colorSortFunction func(*color.Color) int

type pixelSortFilter struct {
    angle float64
    edgeContrast float32
    sortFunction colorSortFunction
    edgeMapped bool
}

func (p *pixelSortFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	return
}

var defaultOptions = gift.Options{
	Parallelization: true,
}

func (p *pixelSortFilter) Draw(dst draw.Image, src image.Image, options *gift.Options){
    if options == nil {
        options = &defaultOptions
    }
    // Call a different function for our segments depending on whether we are edge mapped or not
    var pixelSegments []PixelSegment
    if p.edgeMapped {
        pixelSegments = createSegmentsWithEdgeMap(src, p.angle, p.edgeContrast)
    } else {
        pixelSegments = createSegments(src, p.angle)
    }

    // Iterate over the pixel segments
    for _, segment := range pixelSegments {
        var colorSlice []color.Color
        // Add all of the colors to the colorSlice
        for _, pixel := range segment {
            colorSlice = append(colorSlice, src.At(pixel.X, pixel.Y))
        }
        // Sort the colorslice using the filter sortfunction
        By(p.sortFunction).Sort(colorSlice)
        // Iterate through the pixels again, drawing the color at the given index at the correct position
        for index, pixel := range segment {
            dst.Set(pixel.X, pixel.Y, colorSlice[index])
        }
    }
}

// This function takes an angle in degrees and a sort function which can sort RGBA values and returns a Filter which
// can be used to apply the pixel sorting operation on an image
func PixelSort(angle float64, sortFunction colorSortFunction, edgeMapped bool, edgeContrast float32) gift.Filter{
    return &pixelSortFilter{
        angle: angle,
        sortFunction: sortFunction,
        edgeMapped: edgeMapped,
        edgeContrast: edgeContrast,
    }
}

// Define a by type that is a function that can rank colors
type By colorSortFunction

func (by By) Sort(colorSlice []color.Color) {
    cs := &colorSorter{
        slice: colorSlice,
        by: by,
    }
    sort.Sort(cs)
}

type colorSorter struct {
    slice []color.Color
    by By
}

func (s *colorSorter) Len() int {
    return len(s.slice)
}

func (s *colorSorter) Swap(i, j int) {
    s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
}

func (s *colorSorter) Less(i, j int) bool {
    return s.by(&s.slice[i]) < s.by(&s.slice[j])
}
