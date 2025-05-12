package main

import (
	"examples/7GUIs/circle/circle_cfg"
	"examples/7GUIs/circle/seven_gui_circle_context"
	"examples/7GUIs/circle/seven_gui_circle_panel"
	"github.com/george012/fltk_go"
)

var ctx = seven_gui_circle_context.NewContext()

func main() {
	fltk_go.SetScheme("gtk+")

	win := fltk_go.NewWindow(
		circle_cfg.WIDGET_WIDTH*4+circle_cfg.WIDGET_PADDING*3,
		circle_cfg.WIDGET_HEIGHT*14+circle_cfg.WIDGET_PADDING*3)
	win.SetLabel("Circle Drawer")

	p := seven_gui_circle_panel.NewPanel(win)
	p.Bind(ctx)

	win.End()
	win.Show()
	fltk_go.Run()
}
