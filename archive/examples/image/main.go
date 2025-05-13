package main

import (
	_ "embed"
	"fmt"
	"github.com/george012/fltk_go"
)

//go:embed image.jpg
var testJPG []byte

func main() {
	win := fltk_go.NewWindow(400, 400, "jpeg image example")
	box := fltk_go.NewBox(fltk_go.FLAT_BOX, 0, 0, 400, 400, "")

	image, err := fltk_go.NewJpegImageFromData(testJPG)
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	} else {
		image.Scale(360, 360, true, true)
		box.SetImage(image)
	}
	win.End()
	win.Show()
	fltk_go.Run()
}
