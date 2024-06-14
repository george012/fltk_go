package main

import (
	"fmt"
	"github.com/george012/fltk_go"
)

func main() {
	win := fltk_go.NewWindow(400, 300)
	box := fltk_go.NewBox(fltk_go.FLAT_BOX, 0, 0, 400, 300, "")
	svgImage, err := fltk_go.NewSvgImageFromString(`<svg height="200" width="200">
<polygon points="10,10 20,100 10,190 100,180 190,190 180,100 190,10 100,20" style="stroke:blue;stroke-width:5;fill:red"/>
</svg>`)
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	} else {
		svgImage.Scale(100, 100, true, true)
		box.SetImage(svgImage)
	}
	win.SetIcons([]*fltk_go.RgbImage{&svgImage.RgbImage})
	win.End()
	win.Show()
	fltk_go.Run()
}
