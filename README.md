# Pixel Art SVG Generator

This program is a Golang application that converts an image into an SVG file with a pixel art style.

## How it works

The program takes an image as input and divides it into small squares of a specified size. Each square is then filled with a solid color that best approximates the average color of the pixels within that square. The resulting SVG image is a mosaic made up of these colored squares.

### Considerations

* Generates pixel art SVG from PNG and JPEG files
* Completely written in Golang with no external dependencies
* The quality of the resulting pixel art will depend on the size of the block specified. Smaller block sizes will result in more detailed pixel art, but will also take bit longer to generate.

## Installation

Since this is a Golang application, you will need to have Golang installed on your system in order to run it. You can download and install Golang from the official website: [https://go.dev/dl/](https://go.dev/dl/)

### Running the program

The program is executed from the command line. The following command shows the usage with required arguments:

`go run . -source_image_path <SOURCE_IMAGE_PATH> -output_image_path <OUTPUT_IMAGE_PATH> -block_size <BLOCK_SIZE>`

#### Arguments

* `-source-image-path`: The path to the image file you want to convert.
* `-output-image-path`: The path where you want to save the generated SVG file.
* `-block-size`: The size (width and height) of each pixel in the resulting SVG image.

#### Example

`go run . -source-image-path ./mona_lisa.jpg -output-image-path ./pixel_mona_lisa.svg -block-size 10`

This command will convert the image "mona_lisa.jpg" into a pixel art SVG file named "pixel_mona_lisa.svg" with each pixel being 10x10 squares.
