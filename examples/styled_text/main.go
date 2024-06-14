package main

import (
	"github.com/george012/fltk_go"
	"strings"
)

func main() {
	fltk_go.SetScheme("oxy")

	entries := []fltk_go.StyleTableEntry{
		fltk_go.StyleTableEntry{
			Color: fltk_go.RED,
			Font:  fltk_go.HELVETICA,
			Size:  14,
		},
		fltk_go.StyleTableEntry{
			Color: fltk_go.BLUE,
			Font:  fltk_go.HELVETICA,
			Size:  14,
		},
	}

	buf := fltk_go.NewTextBuffer()
	sbuf := fltk_go.NewTextBuffer()

	window := fltk_go.NewWindow(600, 400)
	disp := fltk_go.NewTextDisplay(0, 0, 600, 400)
	disp.SetBuffer(buf)
	disp.SetHighlightData(sbuf, entries)
	window.End()
	window.Show()

	buf.Append("Hello\n")
	sbuf.Append(strings.Repeat("A", len("Hello\n")))
	buf.Append("World\n")
	sbuf.Append(strings.Repeat("B", len("World\n")))

	fltk_go.Run()
}
