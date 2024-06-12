package main

import (
	"bytes"
	_ "embed"
	"github.com/george012/fltk_go"
	"github.com/yuin/goldmark"
)

//go:embed example.md
var exampleText string

const (
	WIDGET_HEIGHT  = 600
	WIDGET_PADDING = 10
	WIDGET_WIDTH   = 400
)

func main() {
	fltk_go.SetScheme("gtk+")

	win := fltk_go.NewWindow(
		WIDGET_WIDTH*2+WIDGET_PADDING*3,
		WIDGET_HEIGHT*1+WIDGET_PADDING*2)
	win.SetLabel("MDViewer")
	win.Resizable(win)

	hpack := fltk_go.NewPack(WIDGET_PADDING, WIDGET_PADDING, win.W(), WIDGET_HEIGHT)
	hpack.SetType(fltk_go.HORIZONTAL)
	hpack.SetSpacing(WIDGET_PADDING)

	mdEditorBuf := fltk_go.NewTextBuffer()
	mdEditorBuf.SetText(exampleText)
	mdEditor := fltk_go.NewTextEditor(0, 0, WIDGET_WIDTH, WIDGET_HEIGHT)
	mdEditor.SetBuffer(mdEditorBuf)

	previewPanel := fltk_go.NewHelpView(0, 0, WIDGET_WIDTH, WIDGET_HEIGHT)
	updatePreview(mdEditorBuf, previewPanel)

	mdEditor.SetCallbackCondition(fltk_go.WhenChanged)
	mdEditor.SetCallback(func() {
		updatePreview(mdEditorBuf, previewPanel)
	})

	win.End()
	win.Show()
	fltk_go.Run()
}

func updatePreview(mdEditorBuf *fltk_go.TextBuffer, previewPanel *fltk_go.HelpView) {
	var buf bytes.Buffer
	source := []byte(mdEditorBuf.Text())
	if err := goldmark.Convert(source, &buf); err != nil {
		return
	}

	previewPanel.SetValue(buf.String())
}
