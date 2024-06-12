package main

import (
	"fmt"

	"github.com/george012/fltk_go"
)

func main() {
	win := fltk_go.NewWindow(300, 200)
	table := fltk_go.NewTableRow(5, 5, 295, 190)
	table.EnableColumnHeaders()
	table.EnableRowHeaders()
	table.SetColumnCount(3)
	table.SetRowCount(4)
	table.SetDrawCellCallback(func(tc fltk_go.TableContext, row, col, x, y, w, h int) {
		if tc == fltk_go.ContextCell {
			fltk_go.SetDrawFont(fltk_go.HELVETICA, 14)
			fltk_go.DrawBox(fltk_go.FLAT_BOX, x, y, w, h, fltk_go.BLACK)
			fltk_go.DrawBox(fltk_go.FLAT_BOX, x+1, y+1, w-2, h-2, fltk_go.WHITE)
			fltk_go.SetDrawColor(fltk_go.BLACK)
			fltk_go.Draw(fmt.Sprintf("%d", row+col), x, y, w, h, fltk_go.ALIGN_CENTER)
		}
		if tc == fltk_go.ContextRowHeader {
			fltk_go.SetDrawFont(fltk_go.HELVETICA_BOLD, 14)
			fltk_go.DrawBox(fltk_go.UP_BOX, x, y, w, h, fltk_go.BACKGROUND_COLOR)
			fltk_go.SetDrawColor(fltk_go.BLACK)
			fltk_go.Draw(fmt.Sprintf("row %d", row+1), x, y, w, h, fltk_go.ALIGN_CENTER)
		}
		if tc == fltk_go.ContextColHeader {
			fltk_go.SetDrawFont(fltk_go.HELVETICA_BOLD, 14)
			fltk_go.DrawBox(fltk_go.UP_BOX, x, y, w, h, fltk_go.BACKGROUND_COLOR)
			fltk_go.SetDrawColor(fltk_go.BLACK)
			fltk_go.Draw(fmt.Sprintf("col %d", col+1), x, y, w, h, fltk_go.ALIGN_CENTER)
		}
	})
	table.End()
	win.End()
	win.Show()
	fltk_go.Run()
}
