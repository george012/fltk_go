package main

import (
	"fmt"
	"github.com/george012/fltk_go"
)

func main() {

	win := fltk_go.NewWindow(300, 200)
	y := 0
	nInput := fltk_go.NewInput(70, y, 220, 20, "Normal:")
	nInput.SetCallback(func() {

	})
	y += 35
	iInput := fltk_go.NewIntInput(70, y, 220, 20, "Int:")
	y += 35
	fInput := fltk_go.NewFloatInput(70, y, 220, 20, "Float:")
	y += 35
	sInput := fltk_go.NewSecretInput(70, y, 220, 20, "Secret:")
	y += 35
	output := fltk_go.NewOutput(70, y, 220, 20, "Output:")

	updateOutput := func() {
		output.SetValue(fmt.Sprintf("%s %s %s %s",
			nInput.Value(),
			iInput.Value(),
			fInput.Value(),
			sInput.Value(),
		))
	}

	// 为所有输入框绑定回调
	nInput.SetCallbackCondition(fltk_go.WhenChanged)
	nInput.SetCallback(updateOutput)

	iInput.SetCallbackCondition(fltk_go.WhenChanged)
	iInput.SetCallback(updateOutput)

	fInput.SetCallbackCondition(fltk_go.WhenChanged)
	fInput.SetCallback(updateOutput)

	sInput.SetCallbackCondition(fltk_go.WhenChanged)
	sInput.SetCallback(updateOutput)

	// 初始更新一次
	updateOutput()

	win.End()
	win.Show()
	fltk_go.Run()
}
