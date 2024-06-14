package main

import "github.com/george012/fltk_go"

func main() {
	win := fltk_go.NewWindow(300, 200)
	tabs := fltk_go.NewTabs(5, 5, 290, 190)
	pack := fltk_go.NewPack(10, 10, 200, 100, "Tab1")
	pack.SetType(fltk_go.VERTICAL)
	fltk_go.NewButton(0, 0, 100, 20, "Button")
	pack.End()
	tabs.Add(pack)
	pack2 := fltk_go.NewPack(10, 10, 200, 100, "Tab2")
	pack2.SetType(fltk_go.VERTICAL)
	pack2.End()
	tabs.Add(pack2)
	tabs.End()
	tabs.Resizable(pack)

	win.End()
	win.Resizable(tabs)

	win.Show()
	fltk_go.Run()

}
