package main

import (
	"os"
	"strings"

	"github.com/george012/fltk_go"
)

const (
	pad       = 6
	rowHeight = 32
	colWidth  = 60
)

func main() {
	debug := false
	if len(os.Args) > 1 {
		debug = true
	}
	window := makeWindow(debug)
	window.Show()
	fltk_go.Run()
}

func makeWindow(debug bool) *fltk_go.Window {
	width := 200
	height := 115
	window := fltk_go.NewWindow(width, height)
	window.SetLabel("UI Config (grid)")
	makeWidgets(width, height, debug)
	window.End()
	return window
}

func makeWidgets(width, height int, debug bool) {
	grid := fltk_go.NewGrid(0, 0, width, height)
	grid.SetLayout(3, 2, pad, pad)
	makeScaleRow(grid, 0)
	makeThemeRow(grid, 1)
	makeTooltipRow(grid, 2)
	if debug {
		grid.SetShowGrid(true)
	}
	grid.End()
}

func makeScaleRow(grid *fltk_go.Grid, row int) {
	scaleLabel := makeAccelLabel(colWidth, rowHeight, "&Scale")
	scaleSpinner := makeScaleSpinner()
	scaleLabel.SetCallback(func() { scaleSpinner.TakeFocus() })
	scaleSpinner.TakeFocus()
	grid.SetWidget(scaleLabel, row, 0, fltk_go.GridLeft)
	grid.SetWidget(scaleSpinner, row, 1, fltk_go.GridFill)
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

func makeThemeRow(grid *fltk_go.Grid, row int) {
	themeLabel := makeAccelLabel(colWidth, rowHeight, "&Theme")
	themeChoice := makeThemeChoice()
	themeLabel.SetCallback(func() { themeChoice.TakeFocus() })
	grid.SetWidget(themeLabel, row, 0, fltk_go.GridLeft)
	grid.SetWidget(themeChoice, row, 1, fltk_go.GridFill)
}

func makeThemeChoice() *fltk_go.Choice {
	choice := fltk_go.NewChoice(0, 0, colWidth, rowHeight)
	choice.SetTooltip("Sets the application's theme.")
	for i, name := range []string{"&Base", "&Gleam", "G&tk", "&Oxy",
		"&Plastic"} {
		theme := strings.ReplaceAll(name, "&", "")
		if theme == "Oxy" {
			choice.SetValue(i)
			fltk_go.SetScheme(theme)
		}
		choice.Add(name, func() { fltk_go.SetScheme(theme) })
	}
	return choice
}

func makeTooltipRow(grid *fltk_go.Grid, row int) {
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
	grid.SetWidgetWithSpan(checkButton, row, 0, 1, 2, fltk_go.GridCenter)
}

func makeAccelLabel(width, height int, label string) *fltk_go.Button {
	button := fltk_go.NewButton(0, 0, width, height, label)
	button.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
	button.SetBox(fltk_go.NO_BOX)
	button.ClearVisibleFocus()
	return button
}
