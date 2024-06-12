package main

import (
	"github.com/george012/fltk_go"
)

func main() {
	win := fltk_go.NewWindow(800, 600)
	box := fltk_go.NewBox(fltk_go.NO_BOX, 3, 3, 800-6, 600-6)
	box.SetColor(fltk_go.WHITE)
	win.End()
	win.Show()

	offs := fltk_go.NewOffscreen(box.W(), box.H())
	defer offs.Delete()
	offs.Begin()
	fltk_go.SetDrawColor(fltk_go.WHITE)
	fltk_go.DrawRectf(0, 0, box.W(), box.H())
	offs.End()

	box.SetDrawHandler(func(func()) {
		offs.Copy(3, 3, box.W(), box.H(), 0, 0)
	})

	x := 0
	y := 0

	box.SetEventHandler(func(e fltk_go.Event) bool {
		switch e {
		case fltk_go.PUSH:
			offs.Begin()
			fltk_go.SetDrawColor(fltk_go.RED)
			fltk_go.SetLineStyle(fltk_go.SOLID, 3)
			x = fltk_go.EventX()
			y = fltk_go.EventY()
			fltk_go.DrawPoint(x, y)
			offs.End()
			box.Redraw()
			fltk_go.SetLineStyle(fltk_go.SOLID, 0)
			return true
		case fltk_go.DRAG:
			offs.Begin()
			fltk_go.SetDrawColor(fltk_go.RED)
			fltk_go.SetLineStyle(fltk_go.SOLID, 3)
			nx := fltk_go.EventX()
			ny := fltk_go.EventY()
			fltk_go.DrawLine(x, y, nx, ny)
			x = nx
			y = ny
			offs.End()
			box.Redraw()
			fltk_go.SetLineStyle(fltk_go.SOLID, 0)
			return true
		default:
			return false
		}
	})

	fltk_go.Run()
}
