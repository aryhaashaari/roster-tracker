package imgx

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"testing"
)

func TestDraw(t *testing.T) {
	path := "fonts/HerrVonMuellerhoff-Regular.ttf"
	d, err := NewDrawer(Params{
		FontPath:        path,
		BackgroundColor: image.NewUniform(color.Transparent),
		TextColor:       image.NewUniform(color.Black),
		Width:           430,
		Height:          230,
		FontSize:        110,
	})

	// for main sig
	//Width:           430,
	//Height:          230,
	//FontSize:        110,

	// for initial
	//Width:           230,
	//Height:          230,
	//FontSize:        72,

	if err != nil {
		t.Fatal(err.Error())
	}

	img, err := d.Draw("text2img generates the image from a text")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err := os.Create("test.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer file.Close()
	if err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatal(err.Error())
	}

	// for png
	//err = png.Encode(file, img)
	//if err != nil {
	//	t.Fatal(err.Error())
	//}
}
