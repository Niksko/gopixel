package gopixel

import "testing"
import "image"
import "os"
import "math"
import "bufio"
import _ "image/png"
import "fmt"

var createSegmentTests = []struct {
    angle float64
    numberOfPixels int
}{
    {45., 16},
    {0., 16},
    {90., 16},
    {math.Atan(2) * 180. / math.Pi, 16},
    {math.Atan(1./2.) * 180. / math.Pi, 16},
    {180., 16},
    {120., 16},
    {-90., 16},
    {-45, 16},
}

func Test4x4(t *testing.T) {
    for _, value := range createSegmentTests {
        file, _ := os.Open("test/4x4.png")
        angleRadians := float64(value.angle / 180. * math.Pi)
        defer file.Close()
        testImage, _, _ := image.Decode(bufio.NewReader(file))
        pixelSegments := createSegments(testImage, angleRadians)
        pixelCount := 0
        for _, segment := range pixelSegments {
            for range segment{
                pixelCount += 1
            }
        }
        if pixelCount != value.numberOfPixels {
            t.Errorf("Not enough pixels with angle %v\n", value.angle)
            t.Errorf("Expected %v, got %v\n", value.numberOfPixels, pixelCount)
        } else {
            fmt.Printf("Angle: %v\n", value.angle)
            fmt.Printf("%+v\n", pixelSegments)
        }
    }
}
