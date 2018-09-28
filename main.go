package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

var (
	wonb     = flag.Bool("whiteonblack", false, "white text on black background")
	fontfile = flag.String("fontfile", "./resources/font/AvenirLTStd-Heavy.ttf", "filename of the ttf font")
	dpi      = flag.Float64("dpi", 72, "screen resolution")
	hinting  = flag.String("hinting", "full", "none|full")
	size     = flag.Float64("size", 86, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing")
	name     = flag.String("name", "Carlos Muñoz", "Nombre del Usuario")
)

//var name = "Carlos Muñoz"

func main() {
	goldCard, err := os.Open("./resources/goldcard.jpg")
	flag.Parse()
	log.Print(*name)

	if err != nil {
		fmt.Println("la imagen no fue encontrada")
		os.Exit(1)
	}

	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		fmt.Println("la fuente no fue encontrada")
		os.Exit(1)
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer goldCard.Close()

	//bg := image.Black

	img, _, err := image.Decode(goldCard) // abre la imagen
	rgba := image.NewRGBA(image.Rect(0, 0, 948, 597))
	//draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	fg := image.White
	c.SetSrc(fg)

	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	//****** dibuja el texto *****//

	pt := freetype.Pt(42, 128+int(c.PointToFixed(*size)>>6))
	_, err = c.DrawString(*name, pt)
	if err != nil {
		log.Println(err)
		return
	}

	/* ** guarda la imagen a disco ** */
	outFile, err := os.Create("out.jpg")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer outFile.Close()
	//b := bufio.NewWriter(outFile)

	/* ENCODING IN PNG ************************
	var Enc png.Encoder
	//set the best compression
	Enc.CompressionLevel = png.BestCompression

	err = Enc.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	******************************************/

	/**********************************
	**** Encoding in JPG
	**********************************/
	//var Quality = 40
	//var b bytes.Buffer
	var opt jpeg.Options
	opt.Quality = 40

	err = jpeg.Encode(outFile, rgba, &opt)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	/*err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}*/

	fmt.Println("archivo escrito")

}
