package fonts

import (
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	RegularFont font.Face
	BigFont     font.Face
)

func Init() {
	fontBytes, err := os.ReadFile("resources/fonts/font.TTF")
	if err != nil {
		log.Fatal(err)
	}

	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	RegularFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	BigFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

}
