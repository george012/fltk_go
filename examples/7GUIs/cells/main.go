package main

import (
	"errors"
	"fmt"
	"github.com/george012/fltk_go"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	WIDGET_HEIGHT = 200
	WIDGET_WIDTH  = 450
)

const (
	MAX_ROW_COUNT = 100 // 0~99
	MAX_COL_COUNT = 26  // A~Z
)

var reSumFormula = regexp.MustCompile(`(?i)SUM\(([A-Z]\d+):([A-Z]\d+)\)`)

type Context struct {
	Cells [][]*Cell
}

func NewContext(maxRow, maxCol int) *Context {
	ctx := &Context{}
	ctx.Cells = make([][]*Cell, maxRow)
	for row := range ctx.Cells {
		ctx.Cells[row] = make([]*Cell, maxCol)
		for col := range ctx.Cells[row] {
			cell := &Cell{}
			cell.Loc.Row = CellRow(row)
			cell.Loc.Col = CellCol(col)
			cell.Data = &CellDataText{}
			ctx.Cells[row][col] = cell
		}
	}

	return ctx
}

func (ctx *Context) FindCell(loc CellLoc) *Cell {
	if int(loc.Row) >= len(ctx.Cells) || int(loc.Col) >= len(ctx.Cells[loc.Row]) {
		return nil
	}
	return ctx.Cells[loc.Row][loc.Col]
}

func ParseCellData(value string) CellData {
	if strings.HasPrefix(value, "=") {
		formula := value[1:]
		matches := reSumFormula.FindStringSubmatch(formula)
		if matches == nil {
			return &CellDataInvalid{Reason: "SUPPORT FORMULA"}
		}

		strStartLoc := matches[1]
		startLoc, err := ParseCellLoc(strStartLoc)
		if err != nil {
			return &CellDataInvalid{Reason: "INVALID START CELL"}
		}

		strEndLoc := matches[2]
		endLoc, err := ParseCellLoc(strEndLoc)
		if err != nil {
			return &CellDataInvalid{Reason: "INVALID END CELL"}
		}

		formulaData := &CellDataFormula{}
		formulaData.Formula = formula
		formulaData.CalCells = make(map[CellLoc]*Cell)
		for row := startLoc.Row; row <= endLoc.Row; row++ {
			for col := startLoc.Col; col <= endLoc.Col; col++ {
				loc := CellLoc{Row: row, Col: col}
				formulaData.CalCells[loc] = ctx.FindCell(loc)
			}
		}
		return formulaData
	}

	num, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return &CellDataNumber{Number: num}
	}

	return &CellDataText{Text: value}
}

func (cell *Cell) Update(value string) {
	cellData := ParseCellData(value)
	cell.Data = cellData
	cell.RawValue = value
}

func (ctx *Context) FindAllChangedCells(cell *Cell) map[CellLoc]*Cell {
	changedCells := make(map[CellLoc]*Cell)
	changedCells[cell.Loc] = cell

	ctx.findAllChangedCells(cell, changedCells)

	return changedCells
}

func (ctx *Context) findAllChangedCells(cell *Cell, changedCells map[CellLoc]*Cell) {
	for row := range ctx.Cells {
		for col := range ctx.Cells[row] {
			changedCell := ctx.Cells[row][col]
			if changedCell.Data.IsDependOn(cell) {
				if _, exist := changedCells[changedCell.Loc]; !exist {
					changedCells[changedCell.Loc] = changedCell
					ctx.findAllChangedCells(changedCell, changedCells)
				} else {
					changedCell.Data = &CellDataInvalid{Reason: "RECURSIVE FORMULA"}
				}
			}
		}
	}
}

func (ctx *Context) UpdateCellAtLoc(locStr string, value string) {
	loc, err := ParseCellLoc(locStr)
	if err != nil {
		return
	}
	cell := ctx.FindCell(loc)
	if cell == nil {
		return
	}
	cell.Update(value)
}

type Cell struct {
	Loc      CellLoc
	Data     CellData
	RawValue string
}

type CellRow int

func (row CellRow) String() string {
	return strconv.Itoa(int(row))
}

type CellCol int

func reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (col CellCol) String() string {
	//  'A', ... 'Z', 'AA', 'AB', ... 'AZ', 'BA', ... 'BZ', ... 'ZZ', 'AAA', ...
	var stringBytes []byte
	for {
		stringBytes = append(stringBytes, 'A'+byte(col%26))
		col /= 26
		if col <= 0 {
			break
		}
		col -= 1
	}
	reverse(stringBytes)
	return string(stringBytes)
}

type CellLoc struct {
	Row CellRow
	Col CellCol
}

func (loc CellLoc) String() string {
	return loc.Col.String() + loc.Row.String()
}

var ZERO_CELL_LOC = CellLoc{0, 0}

func ParseCellLoc(str string) (CellLoc, error) {
	loc := ZERO_CELL_LOC

	if len(str) < 2 {
		return loc, errors.New("invalid cell location: len(str) < 2")
	}

	str = strings.ToUpper(str)

	runeCol := str[0]
	if runeCol > 'Z' || runeCol < 'A' {
		return ZERO_CELL_LOC, fmt.Errorf("invalid cell column: '%c'", runeCol)
	}
	loc.Col = CellCol(runeCol - 'A')

	strRow := str[1:]
	irow, err := strconv.Atoi(strRow)
	if err != nil {
		return ZERO_CELL_LOC, fmt.Errorf("invalid cell row: '%s'", strRow)
	}
	loc.Row = CellRow(irow)

	return loc, nil
}

type CellData interface {
	Eval() float64
	Display() string
	IsDependOn(cell *Cell) bool
}

type CellDataNumber struct {
	Number float64
}

func (data *CellDataNumber) Eval() float64 {
	return data.Number
}

func (data *CellDataNumber) Display() string {
	return fmt.Sprintf("%f", data.Number)
}

func (data *CellDataNumber) IsDependOn(cell *Cell) bool {
	return false
}

type CellDataFormula struct {
	Formula  string
	CalCells map[CellLoc]*Cell
}

func (data *CellDataFormula) Eval() float64 {
	sum := 0.0
	for _, calCell := range data.CalCells {
		if calCell != nil {
			sum += calCell.Data.Eval()
		}
	}
	return sum
}

func (data *CellDataFormula) Display() string {
	num := data.Eval()
	return fmt.Sprintf("%f", num)
}

func (data *CellDataFormula) IsDependOn(cell *Cell) bool {
	_, ok := data.CalCells[cell.Loc]
	return ok
}

type CellDataInvalid struct {
	Reason string
}

func (data *CellDataInvalid) Eval() float64 {
	return math.NaN()
}

func (data *CellDataInvalid) Display() string {
	return data.Reason
}

func (data *CellDataInvalid) IsDependOn(cell *Cell) bool {
	return false
}

type CellDataText struct {
	Text string
}

func (data *CellDataText) Eval() float64 {
	return 0.0
}

func (data *CellDataText) Display() string {
	return data.Text
}

func (data *CellDataText) IsDependOn(cell *Cell) bool {
	return false
}

type Panel struct {
	tb         *fltk_go.TableRow
	cellValues map[CellLoc]string

	editInput *fltk_go.Input // the input box to show on editing cell
	editCell  *Cell          // current editing cell meta, nil means not editing
}

func NewPanel(win *fltk_go.Window, rowCount, colCount int) *Panel {
	p := &Panel{}

	p.cellValues = make(map[CellLoc]string)

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

func (p *Panel) Bind(ctx *Context) {
	for row := range ctx.Cells {
		for col := range ctx.Cells[row] {
			cell := ctx.Cells[row][col]
			p.cellValues[cell.Loc] = cell.Data.Display()
		}
	}

	p.tb.SetDrawCellCallback(func(tc fltk_go.TableContext, i, j, x, y, w, h int) {
		row := CellRow(i)
		col := CellCol(j)

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
			loc := CellLoc{Row: row, Col: col}
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

func (p *Panel) IsEditingAt(col CellCol, row CellRow) bool {
	return p.editCell != nil && p.editCell.Loc.Col == col && p.editCell.Loc.Row == row
}

func (p *Panel) StartEditing(ctx *Context) {
	if p.IsEditing() {
		p.DoneEditing(ctx)
	}

	row := CellRow(p.tb.CallbackRow())
	col := CellCol(p.tb.CallbackColumn())
	loc := CellLoc{Row: row, Col: col}

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

func (p *Panel) DoneEditing(ctx *Context) {
	if p.IsEditing() {
		// log.Print("done editing")
		p.editCell.Update(p.editInput.Value())
		p.ApplyChangedCells(ctx, p.editCell)
		p.editInput.Hide()
		p.editCell = nil
	}
}

func (p *Panel) ApplyChangedCells(ctx *Context, changedCell *Cell) {
	scells := ctx.FindAllChangedCells(changedCell)
	for loc, scell := range scells {
		p.cellValues[loc] = scell.Data.Display()
	}
}

var ctx = NewContext(MAX_ROW_COUNT, MAX_COL_COUNT)

func init() {
	ctx.UpdateCellAtLoc("B1", "5")
	ctx.UpdateCellAtLoc("B2", "1")
	ctx.UpdateCellAtLoc("B3", "10.3")
	ctx.UpdateCellAtLoc("B4", "22.87")
	ctx.UpdateCellAtLoc("B5", "=SUM(B1:B4)")
	ctx.UpdateCellAtLoc("C1", "6")
	ctx.UpdateCellAtLoc("C2", "7")
	ctx.UpdateCellAtLoc("C3", "2")
	ctx.UpdateCellAtLoc("C4", "5")
	ctx.UpdateCellAtLoc("C5", "=SUM(C1:C4)")
	ctx.UpdateCellAtLoc("A5", "Sum")
	ctx.UpdateCellAtLoc("D5", "=SUM(B5:C5)")
}

func main() {
	fltk_go.SetScheme("gtk+")

	win := fltk_go.NewWindow(
		WIDGET_WIDTH,
		WIDGET_HEIGHT)
	win.SetLabel("Cells")

	p := NewPanel(win, MAX_ROW_COUNT, MAX_COL_COUNT)
	p.Bind(ctx)

	win.End()
	win.Show()
	fltk_go.Run()
}
