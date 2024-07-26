package main

import (
	"github.com/george012/fltk_go"
)

func main() {
	win := fltk_go.NewWindow(320, 180, "grid example")
	grid := fltk_go.NewGrid(0, 0, win.W(), win.H())
	grid.SetLayout(3, 3, 10, 10)
	grid.SetColor(fltk_go.WHITE)
	b0 := fltk_go.NewButton(0, 0, 0, 0, "New")
	b1 := fltk_go.NewButton(0, 0, 0, 0, "Options")
	b2 := fltk_go.NewButton(0, 0, 0, 0, "About")
	b3 := fltk_go.NewButton(0, 0, 0, 0, "Help")
	b4 := fltk_go.NewButton(0, 0, 0, 0, "Quit")
	grid.SetWidget(b0, 0, 0, fltk_go.GridFill)
	grid.SetWidget(b1, 0, 2, fltk_go.GridFill)
	grid.SetWidget(b2, 1, 1, fltk_go.GridFill)
	grid.SetWidget(b3, 2, 0, fltk_go.GridFill)
	grid.SetWidget(b4, 2, 2, fltk_go.GridFill)
	grid.SetShowGrid(false)
	grid.End()
	win.End()
	win.Resizable(grid)
	win.Show()
	fltk_go.Run()
}
