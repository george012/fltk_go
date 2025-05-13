package main

import "github.com/george012/fltk_go"

func main() {
	win := fltk_go.NewWindow(300, 200)
	table := fltk_go.NewTableRow(5, 5, 295, 190)
	table.SetRowCount(2)
	table.SetColumnCount(3)
	table.SetBox(fltk_go.NO_BOX)
	table.Begin()
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			x, y, w, h, err := table.FindCell(fltk_go.ContextCell, i, j)
			if err == nil {
				if j == 0 {
					fltk_go.NewInput(x, y, w, h, "")
				} else {
					fltk_go.NewButton(x, y, w, h, "button")
				}
			}
		}
	}
	table.End()
	win.End()
	win.Show()
	fltk_go.Run()
}
