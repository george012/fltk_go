package main

import (
	"fmt"

	"github.com/george012/fltk_go"
)

func main() {
	win := fltk_go.NewWindow(400, 300)
	box := fltk_go.NewBox(fltk_go.FLAT_BOX, 0, 0, 400, 300, "")
	image, err := fltk_go.NewJpegImageLoad("image.jpg")
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	} else {
		image.Scale(100, 100, true, true)
		box.SetImage(image)
	}
	win.End()
	win.Show()
	fltk_go.Run()
}
