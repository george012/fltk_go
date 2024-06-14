package main

import (
	"fmt"
	"github.com/george012/fltk_go"
	"os"
)

const (
	WIDGET_HEIGHT  = 400
	WIDGET_PADDING = 0
	WIDGET_WIDTH   = 800
)

type EditorApp struct {
	Win        *fltk_go.Window
	TextBuffer *fltk_go.TextBuffer
	TextEditor *fltk_go.TextEditor
	FileName   string
	IsChanged  bool
}

func (app *EditorApp) BuildGUI() {
	fltk_go.SetScheme("gtk+")

	app.Win = fltk_go.NewWindow(WIDGET_WIDTH, WIDGET_HEIGHT)

	app.Win.SetLabel("TextEditor")
	app.Win.Resizable(app.Win)

	col := fltk_go.NewFlex(WIDGET_PADDING, WIDGET_PADDING, app.Win.W(), WIDGET_HEIGHT)
	col.SetType(fltk_go.COLUMN)
	col.SetSpacing(WIDGET_PADDING)

	menuBar := fltk_go.NewMenuBar(0, 0, 0, 0)
	col.Fixed(menuBar, 20)
	menuBar.SetType(uint8(fltk_go.FLAT_BOX))
	menuBar.Activate()
	menuBar.AddEx("File", fltk_go.ALT+'f', nil, fltk_go.SUBMENU)
	menuBar.AddEx("File/&Open", fltk_go.CTRL+'o', app.callbackMenuFileOpen, fltk_go.MENU_VALUE)
	menuBar.AddEx("File/&Save", fltk_go.CTRL+'s', app.callbackMenuFileSave, fltk_go.MENU_VALUE)
	menuBar.Add("File/Save &As", app.callbackMenuFileSaveAs)
	menuBar.AddEx("Help", 0, nil, fltk_go.SUBMENU)
	menuBar.Add("Help/&About", app.callbackMenuHelpAbout)

	app.TextBuffer = fltk_go.NewTextBuffer()
	app.TextEditor = fltk_go.NewTextEditor(0, 0, 0, 0)
	app.TextEditor.SetBuffer(app.TextBuffer)

	app.TextEditor.SetCallbackCondition(fltk_go.WhenChanged)
	app.TextEditor.SetCallback(func() {
		app.IsChanged = true
	})
	app.TextEditor.Parent().Resizable(app.TextEditor)
	col.End()
	app.Win.End()
	app.IsChanged = false
}

func main() {
	myapp := EditorApp{}
	myapp.BuildGUI()
	myapp.Win.Show()
	fltk_go.Run()
}

func (app *EditorApp) callbackMenuFileOpen() {
	fChooser := fltk_go.NewFileChooser("./", "*.*", fltk_go.FileChooser_SINGLE, "Open text file")
	defer fChooser.Destroy()
	fChooser.Popup()
	fnames := fChooser.Selection()
	if len(fnames) == 1 {
		fmt.Printf("select file name %s\n", fnames[0])
		textByte, err := os.ReadFile(fnames[0])
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		app.TextBuffer.SetText(string(textByte))
		app.FileName = fnames[0]
	}
}

func (app *EditorApp) callbackMenuFileSave() {
	if app.IsChanged {
		info, _ := os.Stat(app.FileName)
		os.WriteFile(app.FileName, []byte(app.TextBuffer.Text()), info.Mode())
		app.IsChanged = false
	}
}

func (app *EditorApp) callbackMenuFileSaveAs() {
	fChooser := fltk_go.NewFileChooser("./", "*.*", fltk_go.FileChooser_CREATE, "Enter/Select file name")
	defer fChooser.Destroy()
	fChooser.Popup()
	fnames := fChooser.Selection()
	if len(fnames) == 1 {
		os.WriteFile(fnames[0], []byte(app.TextBuffer.Text()), 0640)
		app.IsChanged = false
		app.FileName = fnames[0]
	}
}

func (app *EditorApp) callbackMenuHelpAbout() {
	fltk_go.MessageBox("About", "Sample Text Editor")
}
