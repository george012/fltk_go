package main

import (
	"strconv"

	"github.com/george012/fltk_go"
)

const (
	WIDGET_HEIGHT  = 25
	WIDGET_PADDING = 5
	WIDGET_WIDTH   = 70
)

func main() {
	fltk_go.SetScheme("gtk+")

	win := fltk_go.NewWindow(
		WIDGET_WIDTH*2+WIDGET_PADDING*2,
		WIDGET_HEIGHT+WIDGET_PADDING*2)
	win.SetLabel("Counter")

	row := fltk_go.NewFlex(WIDGET_PADDING, WIDGET_PADDING, WIDGET_WIDTH*2, WIDGET_HEIGHT)
	row.SetType(fltk_go.ROW)
	row.SetGap(WIDGET_PADDING)

	text := fltk_go.NewOutput(0, 0, 0, 0)
	text.SetValue("0")

	btn := fltk_go.NewButton(0, 0, 0, 0)
	btn.SetLabel("Count")

	value := 0
	btn.SetCallback(func() {
		value++
		text.SetValue(strconv.Itoa(value))
	})

	row.End()
	win.End()
	win.Show()
	fltk_go.Run()
}
