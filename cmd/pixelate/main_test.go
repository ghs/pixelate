package main

import (
	"encoding/xml"
	"image"
	"image/color"
	"image/draw"

	"os"
	"testing"
)

func TestGetDominantColor(t *testing.T) {
	width, height := 30, 30
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	rWidth := int(float64(width) * 0.6)
	gWidth := width - rWidth

	draw.Draw(img, image.Rect(0, 0, rWidth, height), &image.Uniform{C: color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(rWidth, 0, rWidth+gWidth, height), &image.Uniform{C: color.RGBA{0, 255, 0, 255}}, image.Point{}, draw.Src)

	redWeight := float64(rWidth) * float64(height)
	greenWeight := float64(gWidth) * float64(height)

	avgR := (255*redWeight + 0*greenWeight) / (redWeight + greenWeight)
	avgG := (0*redWeight + 255*greenWeight) / (redWeight + greenWeight)
	avgB := 0.0
	avgColor := color.RGBA{
		R: uint8(avgR),
		G: uint8(avgG),
		B: uint8(avgB),
		A: 255,
	}

	domColor := getAvrgDominantColor(img, 0, 0, width, height)

	if domColor != avgColor {
		t.Errorf("Expected color %v but got %v", avgColor, domColor)
	}
}

func TestColorToHex(t *testing.T) {

	red := color.RGBA{255, 0, 0, 255}
	hex := colorToHex(red)

	if hex != "#ff0000" {
		t.Errorf("Expected #ff0000 but got: %s", hex)
	}
}

func TestSVGGeneration(t *testing.T) {

	width, height := 300, 300
	blockSize := 30
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y += blockSize {
		for x := 0; x < width; x += blockSize {
			var blockColor color.Color
			if (x/blockSize+y/blockSize)%2 == 0 {
				blockColor = color.RGBA{255, 0, 0, 255}
			} else {
				blockColor = color.RGBA{0, 255, 0, 255}
			}
			draw.Draw(img, image.Rect(x, y, x+blockSize, y+blockSize), &image.Uniform{C: blockColor}, image.Point{}, draw.Src)
		}
	}

	file, err := os.Create("test.svg")
	if err != nil {
		t.Fatalf("Failed to create test SVG file: %v", err)
	}
	defer os.Remove("test.svg")
	defer file.Close()

	imgBounds := img.Bounds()
	svg := SVG{
		Width:   imgBounds.Dx(),
		Height:  imgBounds.Dy(),
		ViewBox: "0 0 300 300",
		XMLNS:   "http://www.w3.org/2000/svg",
		Rects:   []Rect{},
	}

	for y := 0; y < imgBounds.Dy(); y += blockSize {
		for x := 0; x < imgBounds.Dx(); x += blockSize {
			domColor := getAvrgDominantColor(img, x, y, blockSize, blockSize)
			svg.Rects = append(svg.Rects, Rect{
				X:      x,
				Y:      y,
				Width:  blockSize,
				Height: blockSize,
				Fill:   colorToHex(domColor),
			})
		}
	}

	err = xml.NewEncoder(file).Encode(svg)
	if err != nil {
		t.Errorf("Failed to encode SVG: %v", err)
	}
}
