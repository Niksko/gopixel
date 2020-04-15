# GoPixel

A library and CLI tool for [pixel sorting](https://www.google.com/search?q=pixel+sorting”)

## CLI usage

1. Clone this repo
1. `cd cmd/gopixel`
1. `go build`
1. `./gopixel <input-file> [sort-angle] > <output-file>`

GoPixel can read files in PNG format, and writes PNG formatted files to stdout.

## Library usage

```golang
package main
import (
    "image"
    "log"
    "os"

    gopixel "github.com/niksko/gopixel/pkg/sort"
)

func main() {
    reader, err := os.Open("somepath/myimage.png")
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    var img image.Image
    img, _, err = image.Decode(reader)
    if err != nil {
        log.Fatal(err)
    }

    sortedImage, err := gopixel.Sort(img, 270)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Feature list

See the [feature list](https://github.com/Niksko/gopixel/issues/1).
