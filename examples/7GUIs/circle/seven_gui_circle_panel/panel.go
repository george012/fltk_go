package seven_gui_circle_panel

import (
	"examples/7GUIs/circle/circle_cfg"
	"examples/7GUIs/circle/seven_gui_circle_context"
	"fmt"
	"github.com/george012/fltk_go"
	"runtime"
)

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

	col := fltk_go.NewFlex(circle_cfg.WIDGET_PADDING, circle_cfg.WIDGET_PADDING, win.W()-circle_cfg.WIDGET_PADDING*2, win.H()-circle_cfg.WIDGET_PADDING*2)
	col.SetGap(circle_cfg.WIDGET_PADDING)

	row := fltk_go.NewFlex(0, 0, 0, 0)
	col.Fixed(row, circle_cfg.WIDGET_HEIGHT)
	row.SetType(fltk_go.ROW)
	row.SetGap(circle_cfg.WIDGET_PADDING)

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

	p.adjustDlg = fltk_go.NewWindow(circle_cfg.WIDGET_WIDTH*2+circle_cfg.WIDGET_PADDING*2, circle_cfg.WIDGET_HEIGHT*2+circle_cfg.WIDGET_PADDING*3)
	p.adjustDlg.SetModal()

	row = fltk_go.NewFlex(circle_cfg.WIDGET_PADDING, circle_cfg.WIDGET_PADDING, circle_cfg.WIDGET_WIDTH*2, circle_cfg.WIDGET_HEIGHT*2)
	p.adjustTips = fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0)
	p.adjustTips.SetLabel("Adjust diameter")
	p.adjustSlider = fltk_go.NewSlider(0, 0, 0, 0)
	p.adjustSlider.SetType(fltk_go.HOR_NICE_SLIDER)
	p.adjustSlider.SetMaximum(seven_gui_circle_context.MAX_RADIUS)
	p.adjustSlider.SetMinimum(seven_gui_circle_context.MIN_RADIUS)
	p.adjustSlider.SetValue(seven_gui_circle_context.DEF_RADIUS)
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

func (p *Panel) Bind(ctx *seven_gui_circle_context.Context) {
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
				ctx.AddOp(seven_gui_circle_context.OP_ADD, c)
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
				p.adjustDlg.SetPosition(fltk_go.EventXRoot()-c.R, fltk_go.EventYRoot()-c.R-p.adjustDlg.H()-circle_cfg.WIDGET_PADDING)
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
			ctx.AddOp(seven_gui_circle_context.OP_UPDATE, c)
			return true
		case fltk_go.HIDE:
			c.R = int(p.adjustSlider.Value())
			ctx.AddOp(seven_gui_circle_context.OP_UPDATE, c)
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

func (p *Panel) Update(ctx *seven_gui_circle_context.Context) {
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
