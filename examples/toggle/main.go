package main

import "github.com/george012/fltk_go"

type MyToggleButton struct {
	btn *fltk_go.ToggleButton
}

func NewMyToggleButton(x, y, w, h int) *MyToggleButton {
	btn := fltk_go.NewToggleButton(x, y, w, h, "@+9circle")
	btn.SetColor(0x58585800)
	btn.SetSelectionColor(0x00008B00)
	btn.SetBox(fltk_go.RFLAT_BOX)
	btn.SetDownBox(fltk_go.RFLAT_BOX)
	btn.ClearVisibleFocus()
	btn.SetLabelColor(fltk_go.WHITE)
	btn.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
	btn.SetCallback(func() {
		parent := btn.Parent()
		if btn.Value() {
			btn.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_RIGHT)
			parent.Redraw()
		} else {
			btn.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
			parent.Redraw()
		}
	})
	return &MyToggleButton{btn}
}

func main() {
	fltk_go.InitStyles()
	fltk_go.SetBackgroundColor(0, 0, 0)
	win := fltk_go.NewWindow(200, 200)
	NewMyToggleButton(70, 90, 60, 15)
	win.End()
	win.Show()
	fltk_go.Run()
}
