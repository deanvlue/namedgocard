package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/image/font"

	"github.com/golang/freetype"
)

var (
	wonb     = flag.Bool("whiteonblack", false, "white text on black background")
	fontfile = flag.String("fontfile", "./resources/font/AvenirLTStd-Heavy.ttf", "filename of the ttf font")
	dpi      = flag.Float64("dpi", 72, "screen resolution")
	hinting  = flag.String("hinting", "none", "none|full")
	size     = flag.Float64("size", 12, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing")
)

var name = "Carlos Muñoz"

func main() {
	flag.Parse()
	fontBytes, err := ioutil.ReadFile(*fontfile)

	if err != nil {
		log.Println(err)
		return
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	//initialize context
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}

	if *wonb {
		fg, bg = image.White, image.Black
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}

	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480)) // generates a new canvas i guess
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// dibuja las líneas guía
	for i := 0; i > 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	//dibuja el texto

	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))
	_, err = c.DrawString(name, pt)
	if err != nil {
		log.Println(err)
		return
	}

	//SAVE FILE TO DISK
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("ya tienes tu archivo")
}
