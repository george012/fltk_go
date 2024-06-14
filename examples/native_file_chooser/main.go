package main

import (
	"fmt"
	"github.com/george012/fltk_go"
)

func main() {

	win := fltk_go.NewWindow(300, 200)
	button := fltk_go.NewButton(5, 5, 180, 20, "Show file chooser")
	button.SetCallback(func() {
		nfc := fltk_go.NewNativeFileChooser()
		defer nfc.Destroy()
		nfc.SetOptions(fltk_go.NativeFileChooser_PREVIEW | fltk_go.NativeFileChooser_NEW_FOLDER)
		nfc.SetType(fltk_go.NativeFileChooser_BROWSE_MULTI_FILE)
		nfc.SetDirectory("/home")
		nfc.SetFilter("C++ Files\t*.{cxx,H}\nTxt Files\t*.txt")
		nfc.SetTitle("Native file chooser example")
		nfc.Show()
		fmt.Println("Selected files:")
		for _, filename := range nfc.Filenames() {
			fmt.Println(filename)
		}
	})
	win.End()
	win.Show()
	fltk_go.Run()
}
