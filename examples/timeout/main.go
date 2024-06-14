package main

import (
	"fmt"
	"github.com/george012/fltk_go"
)

func main() {
	win := fltk_go.NewWindow(400, 300)

	fltk_go.NewButton(2, 2, 60, 30, "Test").SetCallback(func() {
		fltk_go.AddTimeout(2.0, timeoutCb)
	})

	win.End()
	win.Show()
	fltk_go.Run()
}

func timeoutCb() {
	fmt.Println("test")
	fltk_go.RepeatTimeout(2.0, timeoutCb)
}
