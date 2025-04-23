package main

import (
	"fmt"
	"github.com/george012/fltk_go"
	"runtime"
)

type OpType int

const (
	OP_ADD OpType = iota
	OP_REMOVE
	OP_UPDATE
)

type Op struct {
	Type   OpType
	Circle *Circle
}

const (
	MAX_RADIUS = 100
	MIN_RADIUS = 10
	DEF_RADIUS = 30
)

type Circle struct {
	X, Y int
	R    int
	ID   int
}

func (c *Circle) Copy() *Circle {
	if c == nil {
		return nil
	}

	return &Circle{
		X:  c.X,
		Y:  c.Y,
		R:  c.R,
		ID: c.ID,
	}
}

type Context struct {
	ops         []*Op
	lastOpIndex int

	circles      []*Circle
	lastCircleID int
	pickedCircle *Circle
}

func NewContext() *Context {
	ctx := &Context{}
	ctx.ops = make([]*Op, 0, 10)
	ctx.circles = make([]*Circle, 0, 10)
	ctx.lastOpIndex = -1
	return ctx
}

func (ctx *Context) NewCircle(x, y int) *Circle {
	c := &Circle{
		X:  x,
		Y:  y,
		R:  DEF_RADIUS,
		ID: ctx.lastCircleID,
	}
	ctx.lastCircleID++
	return c
}

func (ctx *Context) AddCircle(c *Circle) {
	ctx.circles = append(ctx.circles, c)
}

func (ctx *Context) RemoveCircle(c *Circle) {
	for i, cInPool := range ctx.circles {
		if cInPool.ID == c.ID {
			ctx.circles = append(ctx.circles[:i], ctx.circles[i+1:]...)
			if ctx.IsCirclePicked(c) {
				ctx.pickedCircle = nil
			}
			return
		}
	}
}

func (ctx *Context) UpdateCircle(c *Circle) {
	for _, cInPool := range ctx.circles {
		if cInPool.ID == c.ID {
			cInPool.R = c.R
		}
	}
}

func (ctx *Context) PickCircle(pickX, pickY int) *Circle {
	minDoubleR := MAX_RADIUS * MAX_RADIUS
	index := -1
	for i, cInPool := range ctx.circles {
		dx, dy := cInPool.X-pickX, cInPool.Y-pickY
		doubleR := dx*dx + dy*dy
		if doubleR <= cInPool.R*cInPool.R && doubleR < minDoubleR {
			minDoubleR = doubleR
			index = i
		}
	}

	if index < 0 {
		ctx.pickedCircle = nil
		return nil
	}

	ctx.pickedCircle = ctx.circles[index]

	return ctx.circles[index].Copy()
}

func (ctx *Context) UnPickCircle() {
	ctx.pickedCircle = nil
}

func (ctx *Context) IsCirclePicked(c *Circle) bool {
	return ctx.pickedCircle != nil && c.ID == ctx.pickedCircle.ID
}

func (ctx *Context) AddOp(t OpType, c *Circle) {
	if c == nil {
		return
	}

	op := &Op{
		Type:   t,
		Circle: c.Copy(),
	}

	ctx.lastOpIndex++
	ctx.ops = ctx.ops[:ctx.lastOpIndex]
	ctx.ops = append(ctx.ops, op)
}

func (ctx *Context) Undo() bool {
	if ctx.lastOpIndex < 0 {
		return false
	}

	op := ctx.ops[ctx.lastOpIndex]
	ctx.lastOpIndex--
	switch op.Type {
	case OP_ADD:
		ctx.RemoveCircle(op.Circle)
	case OP_REMOVE:
		ctx.AddCircle(op.Circle)
	case OP_UPDATE:
		undoCircle := ctx.ops[ctx.lastOpIndex].Circle
		ctx.lastOpIndex--
		ctx.UpdateCircle(undoCircle)
	}

	return true
}

func (ctx *Context) Redo() bool {
	if ctx.lastOpIndex+1 >= len(ctx.ops) {
		return false
	}

	ctx.lastOpIndex++
	op := ctx.ops[ctx.lastOpIndex]
	switch op.Type {
	case OP_ADD:
		ctx.AddCircle(op.Circle)
	case OP_REMOVE:
		ctx.RemoveCircle(op.Circle)
	case OP_UPDATE:
		ctx.lastOpIndex++
		redoCircle := ctx.ops[ctx.lastOpIndex].Circle
		ctx.UpdateCircle(redoCircle)
	}

	return true
}

func (ctx *Context) HasUndo() bool {
	return ctx.lastOpIndex >= 0
}

func (ctx *Context) HasRedo() bool {
	return ctx.lastOpIndex+1 < len(ctx.ops)
}

func (ctx *Context) Circles() []*Circle {
	return ctx.circles
}

func (ctx *Context) PickedCircle() *Circle {
	return ctx.pickedCircle.Copy()
}

func (ctx *Context) HasPickedCircle() bool {
	return ctx.pickedCircle != nil
}

type Panel struct {
	win          *fltk_go.Window
	undoBtn      *fltk_go.Button
	redoBtn      *fltk_go.Button
	drawBox      *fltk_go.Box
	adjustDlg    *fltk_go.Window
	adjustTips   *fltk_go.Box
	adjustSlider *fltk_go.Slider
	popMenu      *fltk_go.MenuButton
}

func NewPanel(win *fltk_go.Window) *Panel {
	p := &Panel{}
	p.win = win

	col := fltk_go.NewFlex(WIDGET_PADDING, WIDGET_PADDING, win.W()-WIDGET_PADDING*2, win.H()-WIDGET_PADDING*2)
	col.SetGap(WIDGET_PADDING)

	row := fltk_go.NewFlex(0, 0, 0, 0)
	col.Fixed(row, WIDGET_HEIGHT)
	row.SetType(fltk_go.ROW)
	row.SetGap(WIDGET_PADDING)

	fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0) // invisible

	p.undoBtn = fltk_go.NewButton(0, 0, 0, 0)
	p.undoBtn.SetLabel("Undo")

	p.redoBtn = fltk_go.NewButton(0, 0, 0, 0)
	p.redoBtn.SetLabel("Redo")

	fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0) // invisible

	row.End()

	p.drawBox = fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0)

	col.End()
	win.Resizable(col)

	p.adjustDlg = fltk_go.NewWindow(WIDGET_WIDTH*2+WIDGET_PADDING*2, WIDGET_HEIGHT*2+WIDGET_PADDING*3)
	p.adjustDlg.SetModal()

	row = fltk_go.NewFlex(WIDGET_PADDING, WIDGET_PADDING, WIDGET_WIDTH*2, WIDGET_HEIGHT*2)
	p.adjustTips = fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0)
	p.adjustTips.SetLabel("Adjust diameter")
	p.adjustSlider = fltk_go.NewSlider(0, 0, 0, 0)
	p.adjustSlider.SetType(fltk_go.HOR_NICE_SLIDER)
	p.adjustSlider.SetMaximum(MAX_RADIUS)
	p.adjustSlider.SetMinimum(MIN_RADIUS)
	p.adjustSlider.SetValue(DEF_RADIUS)
	row.End()

	// SetModal makes the dialog's close button disappear on Windows
	// A workaround is to make the dialog resizable
	if runtime.GOOS == "windows" {
		p.adjustDlg.Resizable(row)
	}

	p.adjustDlg.End()

	p.popMenu = fltk_go.NewMenuButton(0, 0, 0, 0)
	p.popMenu.SetType(fltk_go.POPUP2)
	p.popMenu.Add("Adjust diameter..", func() {
		p.adjustDlg.Show()
	})

	return p
}

func (p *Panel) Bind(ctx *Context) {
	p.undoBtn.SetCallback(func() {
		if ctx.Undo() {
			p.Update(ctx)
		}
	})

	p.redoBtn.SetCallback(func() {
		if ctx.Redo() {
			p.Update(ctx)
		}
	})

	p.drawBox.SetDrawHandler(func(func()) {
		fltk_go.DrawRectfWithColor(p.drawBox.X(), p.drawBox.Y(), p.drawBox.W(), p.drawBox.H(), fltk_go.WHITE)

		for _, c := range ctx.Circles() {
			x := c.X + p.drawBox.X() - c.R
			y := c.Y + p.drawBox.Y() - c.R
			w := c.R * 2
			h := c.R * 2

			if ctx.IsCirclePicked(c) {
				fltk_go.SetDrawColor(fltk_go.ColorFromRgb(128, 128, 128))
				fltk_go.DrawPie(x, y, w, h, 0, 360)
			} else {
				fltk_go.SetDrawColor(fltk_go.ColorFromRgb(0, 0, 0))
				fltk_go.DrawArc(x, y, w, h, 0, 360)
			}
		}
	})

	p.drawBox.SetEventHandler(func(e fltk_go.Event) bool {
		if fltk_go.EventIsClick() && e == fltk_go.RELEASE {
			x := (fltk_go.EventX() - p.drawBox.X())
			y := (fltk_go.EventY() - p.drawBox.Y())
			switch fltk_go.EventButton() {
			case fltk_go.LeftMouse:
				c := ctx.NewCircle(x, y)
				ctx.AddCircle(c)
				ctx.AddOp(OP_ADD, c)
				p.Update(ctx)
				return true
			case fltk_go.RightMouse:
				c := ctx.PickCircle(x, y)
				if c == nil {
					return true
				}
				p.Update(ctx)

				p.adjustTips.SetLabel(fmt.Sprintf("Adjust the circle at (%d, %d)", c.X, c.Y))
				p.adjustSlider.SetValue(float64(c.R))
				p.adjustDlg.SetPosition(fltk_go.EventXRoot()-c.R, fltk_go.EventYRoot()-c.R-p.adjustDlg.H()-WIDGET_PADDING)
				p.popMenu.Popup()
			}
		}
		return false
	})

	p.adjustDlg.SetEventHandler(func(e fltk_go.Event) bool {
		c := ctx.PickedCircle()
		if c == nil {
			return false
		}
		switch e {
		case fltk_go.SHOW:
			ctx.UpdateCircle(c)
			ctx.AddOp(OP_UPDATE, c)
			return true
		case fltk_go.HIDE:
			c.R = int(p.adjustSlider.Value())
			ctx.AddOp(OP_UPDATE, c)
			ctx.UnPickCircle()
			p.Update(ctx)
			return true
		}
		return false
	})

	p.adjustSlider.SetCallbackCondition(fltk_go.WhenChanged)
	p.adjustSlider.SetCallback(func() {
		c := ctx.PickedCircle()
		if c == nil {
			return
		}

		c.R = int(p.adjustSlider.Value())
		ctx.UpdateCircle(c)
		p.Update(ctx)
	})

	p.Update(ctx)
}

func (p *Panel) Update(ctx *Context) {
	p.drawBox.Redraw()

	if ctx.HasRedo() {
		p.redoBtn.Activate()
	} else {
		p.redoBtn.Deactivate()
	}

	if ctx.HasUndo() {
		p.undoBtn.Activate()
	} else {
		p.undoBtn.Deactivate()
	}
}

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
