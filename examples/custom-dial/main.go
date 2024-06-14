package main

import (
	"github.com/george012/fltk_go"
	"strconv"
)

type MyDial struct {
	main_wid  *fltk_go.Group
	value     *int
	value_box *fltk_go.Box
}

func NewMyDial(x, y, w, h int, label string) *MyDial {
	value := 0
	main_wid := fltk_go.NewGroup(x, y, w, h, label)
	main_wid.SetAlign(fltk_go.ALIGN_TOP)
	main_wid.SetLabelSize(22)
	main_wid.SetLabelColor(fltk_go.ColorFromRgb(0x79, 0x79, 0x79))
	value_box :=
		fltk_go.NewBox(fltk_go.NO_BOX, main_wid.X(), main_wid.Y()+80, main_wid.W(), 40, "0")
	value_box.SetLabelSize(26)
	main_wid.End()

	main_wid.SetDrawHandler(func(func()) {
		fltk_go.SetDrawColor(fltk_go.ColorFromRgb(230, 230, 230))
		fltk_go.DrawPie(main_wid.X(), main_wid.Y(), main_wid.W(), main_wid.H(), 0., 180.)
		fltk_go.SetDrawColor(0xb0bf1a00)
		fltk_go.DrawPie(
			main_wid.X(),
			main_wid.Y(),
			main_wid.W(),
			main_wid.H(),
			float64(100-value)*1.8,
			180.,
		)
		fltk_go.SetDrawColor(fltk_go.WHITE)
		fltk_go.DrawPie(
			main_wid.X()-50+main_wid.W()/2,
			main_wid.Y()-50+main_wid.H()/2,
			100,
			100,
			0.,
			360.,
		)
		main_wid.DrawChildren()
	})
	return &MyDial{
		main_wid,
		&value,
		value_box,
	}
}

func (d *MyDial) SetValue(val int) {
	*d.value = val
	d.value_box.SetLabel(strconv.Itoa(val))
	d.main_wid.Redraw()
}

func (d *MyDial) Value() int {
	return *d.value
}

func main() {
	win := fltk_go.NewWindow(400, 300)
	win.SetColor(fltk_go.WHITE)
	dial := NewMyDial(100, 100, 200, 200, "CPU Load %")
	win.End()
	win.Show()
	dial.SetValue(26)
	fltk_go.Run()
}
