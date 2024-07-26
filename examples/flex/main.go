package main

import (
	"github.com/george012/fltk_go"
	"runtime"
	"strconv"
)

var i = 0

func main() {
	// 锁定当前的 goroutine 到操作系统线程
	runtime.LockOSThread()
	// 锁定 FLTK 库
	fltk_go.Lock()

	win := fltk_go.NewWindow(300, 200, "flex example")
	column := fltk_go.NewFlex(0, 0, 300, 200)
	column.SetType(fltk_go.COLUMN)
	column.SetGap(5)
	inc := fltk_go.NewButton(0, 0, 0, 0, "Increment")
	column.Fixed(inc, 40)
	box := fltk_go.NewBox(fltk_go.FLAT_BOX, 0, 0, 0, 0, "0")
	dec := fltk_go.NewButton(0, 0, 0, 0, "Decrement")
	inc.SetCallback(func() {
		i++
		box.SetLabel(strconv.Itoa(i))
	})
	dec.SetCallback(func() {
		i--
		box.SetLabel(strconv.Itoa(i))
	})

	column.Fixed(dec, 40)
	column.End()
	win.End()
	win.Resizable(column)
	win.Show()
	fltk_go.Run()
}
