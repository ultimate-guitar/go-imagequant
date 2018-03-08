/*
Copyright (c) 2016, The go-imagequant author(s)

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND ISC DISCLAIMS ALL WARRANTIES WITH REGARD TO
THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS.
IN NO EVENT SHALL ISC BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR
CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA
OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS
ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS
SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"code.ivysaur.me/imagequant"
)

func GoImageToRgba32(im image.Image) []byte {
	w := im.Bounds().Max.X
	h := im.Bounds().Max.Y
	ret := make([]byte, w*h*4)

	p := 0

	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			r16, g16, b16, a16 := im.At(x, y).RGBA() // Each value ranges within [0, 0xffff]

			ret[p+0] = uint8(r16 >> 8)
			ret[p+1] = uint8(g16 >> 8)
			ret[p+2] = uint8(b16 >> 8)
			ret[p+3] = uint8(a16 >> 8)
			p += 4
		}
	}

	return ret
}

func Rgb8PaletteToGoImage(w, h int, rgb8data []byte, pal color.Palette) image.Image {
	rect := image.Rectangle{
		Max: image.Point{
			X: w,
			Y: h,
		},
	}

	ret := image.NewPaletted(rect, pal)

	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			ret.SetColorIndex(x, y, rgb8data[y*w+x])
		}
	}

	return ret
}

func Crush(sourcefile, destfile string, speed int) error {

	fh, err := os.OpenFile(sourcefile, os.O_RDONLY, 0444)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %s", err.Error())
	}
	defer fh.Close()

	img, err := png.Decode(fh)
	if err != nil {
		return fmt.Errorf("png.Decode: %s", err.Error())
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	attr, err := imagequant.NewAttributes()
	if err != nil {
		return fmt.Errorf("NewAttributes: %s", err.Error())
	}
	defer attr.Release()

	err = attr.SetSpeed(speed)
	if err != nil {
		return fmt.Errorf("SetSpeed: %s", err.Error())
	}

	rgba32data := GoImageToRgba32(img)

	iqm, err := imagequant.NewImage(attr, string(rgba32data), width, height, 0)
	if err != nil {
		return fmt.Errorf("NewImage: %s", err.Error())
	}
	defer iqm.Release()

	res, err := iqm.Quantize(attr)
	if err != nil {
		return fmt.Errorf("Quantize: %s", err.Error())
	}
	defer res.Release()

	rgb8data, err := res.WriteRemappedImage()
	if err != nil {
		return fmt.Errorf("WriteRemappedImage: %s", err.Error())
	}

	im2 := Rgb8PaletteToGoImage(res.GetImageWidth(), res.GetImageHeight(), rgb8data, res.GetPalette())

	fh2, err := os.OpenFile(destfile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %s", err.Error())
	}
	defer fh2.Close()

	penc := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	err = penc.Encode(fh2, im2)
	if err != nil {
		return fmt.Errorf("png.Encode: %s", err.Error())
	}

	return nil
}

func main() {
	ShouldDisplayVersion := flag.Bool("Version", false, "")
	Infile := flag.String("In", "", "Input filename")
	Outfile := flag.String("Out", "", "Output filename")
	Speed := flag.Int("Speed", 3, "Speed (1 slowest, 10 fastest)")

	flag.Parse()

	if *ShouldDisplayVersion {
		fmt.Printf("libimagequant '%s' (%d)\n", imagequant.GetLibraryVersionString(), imagequant.GetLibraryVersion())
		os.Exit(1)
	}

	err := Crush(*Infile, *Outfile, *Speed)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
