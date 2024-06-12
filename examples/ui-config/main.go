package main

import (
	"strings"

	"github.com/george012/fltk_go"
)

const (
	pad       = 6
	rowHeight = 32
	colWidth  = 60
)

func main() {
	window := makeWindow()
	window.Show()
	fltk_go.Run()
}

func makeWindow() *fltk_go.Window {
	width := 200
	height := 115
	window := fltk_go.NewWindow(width, height)
	window.SetLabel("UI Config")
	makeWidgets(width, height)
	window.End()
	return window
}

func makeWidgets(width, height int) {
	colFlex := fltk_go.NewFlex(0, 0, width, height)
	colFlex.SetSpacing(pad)
	rowFlex := makeScaleRow(width, rowHeight)
	colFlex.Fixed(rowFlex, rowHeight)
	rowFlex = makeThemeRow(rowHeight, rowHeight)
	colFlex.Fixed(rowFlex, rowHeight)
	rowFlex = makeTooltipRow(rowHeight, rowHeight)
	colFlex.Fixed(rowFlex, rowHeight)
	colFlex.End()
}

func makeScaleRow(width, height int) *fltk_go.Flex {
	rowFlex := fltk_go.NewFlex(0, 0, width, height)
	rowFlex.SetType(fltk_go.ROW)
	rowFlex.SetSpacing(pad)
	scaleLabel := makeAccelLabel(colWidth, rowHeight, "&Scale")
	scaleSpinner := makeScaleSpinner()
	scaleLabel.SetCallback(func() { scaleSpinner.TakeFocus() })
	rowFlex.Fixed(scaleLabel, colWidth)
	rowFlex.End()
	scaleSpinner.TakeFocus()
	return rowFlex
}

func makeScaleSpinner() *fltk_go.Spinner {
	spinner := fltk_go.NewSpinner(0, 0, colWidth, rowHeight)
	spinner.SetTooltip("Sets the application's scale.")
	spinner.SetType(fltk_go.SPINNER_FLOAT_INPUT)
	spinner.SetMinimum(0.5)
	spinner.SetMaximum(3.5)
	spinner.SetStep(0.1)
	spinner.SetValue(float64(fltk_go.ScreenScale(0)))
	spinner.SetCallback(func() {
		fltk_go.SetScreenScale(0, float32(spinner.Value()))
	})
	return spinner
}

func makeThemeRow(width, height int) *fltk_go.Flex {
	rowFlex := fltk_go.NewFlex(0, 0, width, height)
	rowFlex.SetType(fltk_go.ROW)
	rowFlex.SetSpacing(pad)
	themeLabel := makeAccelLabel(colWidth, rowHeight, "&Theme")
	themeChoice := makeThemeChoice()
	themeLabel.SetCallback(func() { themeChoice.TakeFocus() })
	rowFlex.Fixed(themeLabel, colWidth)
	rowFlex.End()
	return rowFlex
}

func makeThemeChoice() *fltk_go.Choice {
	choice := fltk_go.NewChoice(0, 0, colWidth, rowHeight)
	choice.SetTooltip("Sets the application's theme.")
	for i, name := range []string{"&Base", "&Gleam", "G&tk", "&Oxy", "&Plastic"} {
		theme := strings.ReplaceAll(name, "&", "")
		if theme == "Oxy" {
			choice.SetValue(i)
			fltk_go.SetScheme(theme)
		}
		choice.Add(name, func() { fltk_go.SetScheme(theme) })
	}
	return choice
}

func makeTooltipRow(width, height int) *fltk_go.Flex {
	rowFlex := fltk_go.NewFlex(0, 0, width, height)
	rowFlex.SetType(fltk_go.ROW)
	rowFlex.SetSpacing(pad)
	padBox := fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, colWidth, rowHeight)
	checkButton := fltk_go.NewCheckButton(colWidth, 0, colWidth, rowHeight,
		"S&how Tooltips")
	checkButton.SetTooltip("If checked the application shows tooltips.")
	checkButton.SetValue(fltk_go.AreTooltipsEnabled())
	checkButton.SetCallback(func() {
		if checkButton.Value() {
			fltk_go.EnableTooltips()
		} else {
			fltk_go.DisableTooltips()
		}
	})
	rowFlex.Fixed(padBox, colWidth)
	rowFlex.End()
	return rowFlex
}

func makeAccelLabel(width, height int, label string) *fltk_go.Button {
	button := fltk_go.NewButton(0, 0, width, height, label)
	button.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
	button.SetBox(fltk_go.NO_BOX)
	button.ClearVisibleFocus()
	return button
}
