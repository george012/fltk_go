package main

import (
	"strconv"

	"github.com/george012/fltk_go"
)

var i = 0

func main() {
	win := fltk_go.NewWindow(300, 200)
	column := fltk_go.NewFlex(0, 0, 300, 200)
	column.SetType(fltk_go.COLUMN)
	column.SetGap(5)
	inc := fltk_go.NewButton(0, 0, 0, 0, "Increment")
	column.Fixed(inc, 40)
	box := fltk_go.NewBox(fltk_go.FLAT_BOX, 0, 0, 0, 0, "0")
	dec := fltk_go.NewButton(0, 0, 0, 0, "Decrement")
	inc.SetCallback(func() {
		i++
		box.SetLabel(strconv.Itoa(i))
	})
	dec.SetCallback(func() {
		i--
		box.SetLabel(strconv.Itoa(i))
	})

	column.Fixed(dec, 40)
	column.End()
	win.End()
	win.Resizable(column)
	win.Show()
	fltk_go.Run()
}
