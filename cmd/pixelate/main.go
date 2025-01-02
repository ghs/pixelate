package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type SVG struct {
	Height  int      `xml:"height,attr"`
	Width   int      `xml:"width,attr"`
	XMLName xml.Name `xml:"svg"`
	XMLNS   string   `xml:"xmlns,attr"`
	Rects   []Rect   `xml:"rect"`
	ViewBox string   `xml:"viewBox,attr"`
}

type Rect struct {
	X      int    `xml:"x,attr"`
	Y      int    `xml:"y,attr"`
	Fill   string `xml:"fill,attr"`
	Height int    `xml:"height,attr"`
	Width  int    `xml:"width,attr"`
}

func main() {

	inputPath := flag.String("source-image-path", "", "Set the value for the image which needs to be converted.")
	outputPath := flag.String("output-image-path", "", "Set the value for the output SVG which will be generated.")
	blockSize := flag.Int("block-size", 0, "Set the size for each 'pixel' block within the SVG.")

	flag.Parse()

	if *inputPath == "" || *outputPath == "" {
		panic("'source-image-path' and 'output-image-path' arguments must be set.")
	}

	if *blockSize <= 0 && *blockSize <= 30 {
		panic("'block-size' argument must be greater than zero and lower than or equal to 30 for optimal generation.")
	}

	file, err := os.Open(*inputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	imgBounds := img.Bounds()
	svg := SVG{
		Width:   imgBounds.Dx(),
		Height:  imgBounds.Dy(),
		ViewBox: fmt.Sprintf("0 0 %d %d", imgBounds.Dx(), imgBounds.Dy()),
		Rects:   []Rect{},
		XMLNS:   "http://www.w3.org/2000/svg",
	}

	for y := 0; y < imgBounds.Dy(); y += *blockSize {
		for x := 0; x < imgBounds.Dx(); x += *blockSize {
			domColor := getAvrgDominantColor(img, x, y, *blockSize, *blockSize)
			svg.Rects = append(svg.Rects, Rect{
				X:      x,
				Y:      y,
				Width:  *blockSize,
				Height: *blockSize,
				Fill:   colorToHex(domColor),
			})
		}
	}

	svgFile, err := os.Create(*outputPath)
	if err != nil {
		panic(err)
	}
	defer svgFile.Close()

	encoder := xml.NewEncoder(svgFile)
	encoder.Indent("", "  ")
	err = encoder.Encode(svg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("SVG generated: %s", *outputPath)
}

func getAvrgDominantColor(img image.Image, startX, startY, width, height int) color.Color {

	var rSum, gSum, bSum, count int
	for y := startY; y < startY+height && y < img.Bounds().Max.Y; y++ {
		for x := startX; x < startX+width && x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rSum += int(r >> 8)
			gSum += int(g >> 8)
			bSum += int(b >> 8)
			count++
		}
	}

	return color.RGBA{
		R: uint8(rSum / count),
		G: uint8(gSum / count),
		B: uint8(bSum / count),
		A: 255,
	}
}

func colorToHex(c color.Color) string {

	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}
