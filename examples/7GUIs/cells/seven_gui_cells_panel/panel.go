package seven_gui_cells_panel

import (
	"examples/7GUIs/cells/sven_gui_cells_context"
	"github.com/george012/fltk_go"
	"log"
)

type Panel struct {
	tb         *fltk_go.TableRow
	cellValues map[sven_gui_cells_context.CellLoc]string

	editInput *fltk_go.Input               // the input box to show on editing cell
	editCell  *sven_gui_cells_context.Cell // current editing cell meta, nil means not editing
	Ctx       *sven_gui_cells_context.Context
}

func NewPanel(win *fltk_go.Window, rowCount, colCount int, aCtx *sven_gui_cells_context.Context) *Panel {
	p := &Panel{
		Ctx: aCtx,
	}

	p.cellValues = make(map[sven_gui_cells_context.CellLoc]string)

	p.tb = fltk_go.NewTableRow(0, 0, win.W(), win.H())
	p.tb.SetRowCount(rowCount)
	p.tb.SetColumnCount(colCount)
	p.tb.EnableColumnHeaders()
	p.tb.EnableRowHeaders()
	p.tb.AllowColumnResizing()
	p.tb.AllowRowResizing()

	p.tb.Begin()

	p.editInput = fltk_go.NewInput(0, 0, 0, 0)
	p.editInput.Hide()
	p.editInput.SetColor(fltk_go.YELLOW)

	p.tb.End()
	win.Resizable(p.tb)

	return p
}

func (p *Panel) Bind(ctx *sven_gui_cells_context.Context) {
	for row := range ctx.Cells {
		for col := range ctx.Cells[row] {
			cell := ctx.Cells[row][col]
			p.cellValues[cell.Loc] = cell.Data.Display()
		}
	}

	p.tb.SetDrawCellCallback(func(tc fltk_go.TableContext, i, j, x, y, w, h int) {
		row := sven_gui_cells_context.CellRow(i)
		col := sven_gui_cells_context.CellCol(j)

		switch tc {
		case fltk_go.ContextRowHeader:
			fltk_go.SetDrawFont(fltk_go.HELVETICA_BOLD, 14)
			fltk_go.DrawBox(fltk_go.UP_BOX, x, y, w, h, fltk_go.BACKGROUND_COLOR)
			fltk_go.SetDrawColor(fltk_go.BLACK)
			fltk_go.Draw(row.String(), x, y, w, h, fltk_go.ALIGN_CENTER)
		case fltk_go.ContextColHeader:
			fltk_go.SetDrawFont(fltk_go.HELVETICA_BOLD, 14)
			fltk_go.DrawBox(fltk_go.UP_BOX, x, y, w, h, fltk_go.BACKGROUND_COLOR)
			fltk_go.SetDrawColor(fltk_go.BLACK)
			fltk_go.Draw(col.String(), x, y, w, h, fltk_go.ALIGN_CENTER)
		case fltk_go.ContextCell:
			loc := sven_gui_cells_context.CellLoc{Row: row, Col: col}
			if p.IsEditingAt(col, row) {
				p.editInput.Resize(x, y, w, h)
				return
			}
			fltk_go.SetDrawFont(fltk_go.HELVETICA, 14)
			fltk_go.DrawBox(fltk_go.FLAT_BOX, x, y, w, h, fltk_go.BLACK)
			fltk_go.DrawBox(fltk_go.FLAT_BOX, x+1, y+1, w-2, h-2, fltk_go.WHITE)
			fltk_go.SetDrawColor(fltk_go.BLACK)
			fltk_go.Draw(p.cellValues[loc], x, y, w, h, fltk_go.ALIGN_CENTER)
		}
	})

	// p.tb.SetCallbackCondition(fltk.WhenNotChanged)
	p.tb.SetCallback(func() {
		tc := p.tb.CallbackContext()
		if tc != fltk_go.ContextCell {
			p.DoneEditing(ctx)
			return
		}

		if fltk_go.EventClicks() == 0 {
			p.DoneEditing(ctx)
			return
		}

		p.StartEditing(ctx)
	})

	p.editInput.SetCallbackCondition(fltk_go.WhenEnterKeyAlways)
	p.editInput.SetCallback(func() {
		p.DoneEditing(ctx)
	})
}

func (p *Panel) IsEditing() bool {
	return p.editCell != nil
}

func (p *Panel) IsEditingAt(col sven_gui_cells_context.CellCol, row sven_gui_cells_context.CellRow) bool {
	return p.editCell != nil && p.editCell.Loc.Col == col && p.editCell.Loc.Row == row
}

func (p *Panel) StartEditing(ctx *sven_gui_cells_context.Context) {
	if p.IsEditing() {
		p.DoneEditing(ctx)
	}

	row := sven_gui_cells_context.CellRow(p.tb.CallbackRow())
	col := sven_gui_cells_context.CellCol(p.tb.CallbackColumn())
	loc := sven_gui_cells_context.CellLoc{Row: row, Col: col}

	x, y, w, h, err := p.tb.FindCell(fltk_go.ContextCell, int(row), int(col))
	if err != nil {
		log.Panic("should not go here")
		return
	}

	// log.Print("show input:", x, y, w, h)
	p.editCell = ctx.FindCell(loc)
	p.editInput.Resize(x, y, w, h)
	p.editInput.SetValue(p.editCell.RawValue)
	p.editInput.Show()
	p.editInput.TakeFocus()
}

func (p *Panel) DoneEditing(ctx *sven_gui_cells_context.Context) {
	if p.IsEditing() {
		// log.Print("done editing")
		p.editCell.Update(p.editInput.Value(), p.Ctx)
		p.ApplyChangedCells(ctx, p.editCell)
		p.editInput.Hide()
		p.editCell = nil
	}
}

func (p *Panel) ApplyChangedCells(ctx *sven_gui_cells_context.Context, changedCell *sven_gui_cells_context.Cell) {
	scells := ctx.FindAllChangedCells(changedCell)
	for loc, scell := range scells {
		p.cellValues[loc] = scell.Data.Display()
	}
}
