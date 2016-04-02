package gopixel

import "image"
import "image/color"

type PixelSegment struct{
    startPoint image.Point
    pixelSlice []color.Color
}
