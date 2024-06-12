package main

import (
	"math/rand"
	"strconv"

	"github.com/george012/fltk_go"
)

func main() {
	win := fltk_go.NewWindow(301, 440)
	colors := []fltk_go.Color{
		fltk_go.BLUE,
		fltk_go.RED,
		fltk_go.YELLOW,
		fltk_go.GREEN,
		fltk_go.CYAN,
		fltk_go.DARK_MAGENTA,
	}

	ch := fltk_go.NewChart(1, 1, 300, 300, "Example chart")

	ch.SetTextColor(fltk_go.DARK_RED)
	ch.SetType(fltk_go.SPECIALPIE_CHART)
	ch.SetAutosize(true)

	valueNameEditor := fltk_go.NewInput(85, ch.Y()+ch.H()+30, 200, 20, "Value name")
	valueEditor := fltk_go.NewInput(85, valueNameEditor.Y()+valueNameEditor.H()+20, 200, 20, "Value")

	addValueButton := fltk_go.NewButton(222, valueEditor.Y()+valueEditor.H()+10, 70, 20, "Add value")
	addValueButton.SetCallback(func() {
		val, err := strconv.ParseFloat(valueEditor.Value(), 64)
		if err != nil {
			panic(err)
		}
		ch.Add(val, colors[rand.Intn(len(colors))], valueNameEditor.Value())
	})

	win.End()
	win.Show()
	fltk_go.Run()
}
