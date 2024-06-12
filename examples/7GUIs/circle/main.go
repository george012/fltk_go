package main

import (
	"github.com/george012/fltk_go"
)

const (
	WIDGET_HEIGHT  = 25
	WIDGET_PADDING = 10
	WIDGET_WIDTH   = 100
)

var ctx = NewContext()

func main() {
	fltk_go.SetScheme("gtk+")

	win := fltk_go.NewWindow(
		WIDGET_WIDTH*4+WIDGET_PADDING*3,
		WIDGET_HEIGHT*14+WIDGET_PADDING*3)
	win.SetLabel("Circle Drawer")

	p := NewPanel(win)
	p.Bind(ctx)

	win.End()
	win.Show()
	fltk_go.Run()
}
