package gopixel

import "github.com/disintegration/gift"
import "image"
import "image/draw"

type pixelSortFilter struct {
    angle float64
    sortFunction func(int, int, int, int) int
    edgeMapped bool
}

func (p *pixelSortFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	return
}

func (p *pixelSortFilter) Draw(dst draw.Image, src image.Image, options *gift.Options){
}

// This function takes an angle in degrees and a sort function which can sort RGBA values and returns a Filter which
// can be used to apply the pixel sorting operation on an image
func PixelSort(angle float64, sortFunction func(int, int, int, int) int, edgeMapped bool) gift.Filter{
    return &pixelSortFilter{
        angle: angle,
        sortFunction: sortFunction,
        edgeMapped: edgeMapped,
    }
}
