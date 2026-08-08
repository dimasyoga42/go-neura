package helper

import (
	"image"

	"github.com/fogleman/gg"
)

func DrawText (img image.Image, fontPath string, fontSize float64, topText, bottomText string ) (image.Image, error) {
	dc := gg.NewContextForImage(img)
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return nil, err
	}

	w := float64(img.Bounds().Dx())
	h := float64(img.Bounds().Dy())
	drawOutlinedText(dc, topText, w/2, 40)
	drawOutlinedText(dc, bottomText, w/2, h-40)
	return  dc.Image(), nil
}

func drawOutlinedText(dc *gg.Context, text string, x, y float64) {
	dc.SetRGB(0, 0, 0)

	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			dc.DrawStringAnchored(text, x+float64(dx), y+float64(dy), 0.5, 0.5)
		}
	}

	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
}
