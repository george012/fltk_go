package main

import (
	"fmt"
	"github.com/george012/fltk_go"
)

var disp *fltk_go.TextDisplay
var testCnt = 0

func main() {
	win := fltk_go.NewWindow(800, 600)

	fltk_go.NewButton(2, 2, 60, 30, "Test").SetCallback(func() {
		fltk_go.AddTimeout(2.0, timeoutCb)
	})

	buf := fltk_go.NewTextBuffer()
	disp = fltk_go.NewTextDisplay(2, 64, 796, 530)
	disp.SetBuffer(buf)

	win.End()
	win.Show()
	fltk_go.Run()
}

func timeoutCb() {
	disp.Buffer().SetText(fmt.Sprintf("%stimeout cycle run %d\n", disp.Buffer().Text(), testCnt))
	fltk_go.RepeatTimeout(2.0, timeoutCb)
	testCnt++
}
