package main

import (
	"fmt"
	"time"

	"github.com/george012/fltk_go"
)

const (
	WIDGET_HEIGHT  = 25
	WIDGET_PADDING = 5
	WIDGET_WIDTH   = 200
)

const (
	DATE_FORMAT = "02.01.2006"
)

type BookOption int

const (
	BookOptionOneWay BookOption = iota
	BookOptionReturn
)

func main() {
	var (
		option       BookOption
		optionChoice *fltk_go.Choice
		startInput   *fltk_go.Input
		returnInput  *fltk_go.Input
		bookBtn      *fltk_go.Button
	)

	update := func() {
		option = BookOption(optionChoice.Value())
		switch option {
		case BookOptionOneWay:
			returnInput.Deactivate()
			if _, ok := validateInput(startInput); ok {
				bookBtn.Activate()
			} else {
				bookBtn.Deactivate()
			}
		case BookOptionReturn:
			returnInput.Activate()
			t1, ok1 := validateInput(startInput)
			t2, ok2 := validateInput(returnInput)
			if ok1 && ok2 && !t1.After(t2) {
				bookBtn.Activate()
			} else {
				bookBtn.Deactivate()
			}
		}
	}

	fltk_go.SetScheme("gtk+")

	win := fltk_go.NewWindow(
		WIDGET_WIDTH+WIDGET_PADDING*2,
		WIDGET_HEIGHT*4+WIDGET_PADDING*2)
	win.SetLabel("Book Flight")

	col := fltk_go.NewFlex(WIDGET_PADDING, WIDGET_PADDING, WIDGET_WIDTH, WIDGET_HEIGHT*4)
	col.SetType(fltk_go.COLUMN)
	col.SetGap(WIDGET_PADDING)

	option = BookOptionOneWay

	optionChoice = fltk_go.NewChoice(0, 0, 0, 0)
	optionChoice.Add("one-way flight", update)
	optionChoice.Add("return flight", update)
	optionChoice.SetValue(int(option))

	now := time.Now()
	startInput = fltk_go.NewInput(0, 0, 0, 0)
	startInput.SetValue(now.Format(DATE_FORMAT))
	startInput.SetCallbackCondition(fltk_go.WhenChanged)
	startInput.SetCallback(update)

	returnInput = fltk_go.NewInput(0, 0, 0, 0)
	returnInput.SetValue(now.Format(DATE_FORMAT))
	returnInput.SetCallbackCondition(fltk_go.WhenChanged)
	returnInput.SetCallback(update)
	// defaultInputColor = startInput.Color()

	bookBtn = fltk_go.NewButton(0, 0, 0, 0)
	bookBtn.SetLabel("Book")
	bookBtn.SetCallback(func() {
		switch option {
		case BookOptionOneWay:
			fltk_go.MessageBox("Book successful", fmt.Sprintf("You have booked a one-way flight on %s.", startInput.Value()))
		case BookOptionReturn:
			fltk_go.MessageBox("Book successful", fmt.Sprintf("You have booked a return flight on %s and %s.", startInput.Value(), returnInput.Value()))
		}
	})

	update()

	col.End()
	win.End()
	win.Show()
	fltk_go.Run()
}

func validateInput(input *fltk_go.Input) (time.Time, bool) {
	defer input.Redraw()

	t, err := time.Parse(DATE_FORMAT, input.Value())
	if err != nil {
		input.SetColor(fltk_go.RED)
		return t, false
	}

	input.SetColor(fltk_go.BACKGROUND2_COLOR)
	return t, true
}
