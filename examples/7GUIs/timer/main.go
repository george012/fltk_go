package main

import (
	"time"

	"github.com/george012/fltk_go"
)

const (
	WIDGET_HEIGHT      = 25
	WIDGET_PADDING     = 10
	WIDGET_WIDTH       = 180
	WIDGET_LABEL_WIDTH = 110
)

const (
	DURATION_DEFAULT = 15.0
	DURATION_MAXIMUM = 30.0
)

func main() {
	fltk_go.SetScheme("gtk+")
	fltk_go.Lock() // enable the FLTK lock mechanism

	startTime := time.Now()

	win := fltk_go.NewWindow(
		WIDGET_LABEL_WIDTH+WIDGET_WIDTH+WIDGET_PADDING*2,
		WIDGET_HEIGHT*4+WIDGET_PADDING*5)
	win.SetLabel("Timer")

	col := fltk_go.NewFlex(WIDGET_PADDING, WIDGET_PADDING, win.W()-WIDGET_PADDING*2, win.H()-WIDGET_PADDING*2)
	col.SetGap(WIDGET_PADDING)

	row := fltk_go.NewFlex(0, 0, 0, 0)
	row.SetType(fltk_go.ROW)
	label := fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0, "Elapsed Time:")
	label.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
	row.Fixed(label, WIDGET_LABEL_WIDTH)
	elapsedProgess := fltk_go.NewProgress(0, 0, 0, 0)
	elapsedProgess.SetSelectionColor(fltk_go.BLUE)
	elapsedProgess.SetMaximum(DURATION_DEFAULT)
	row.End()

	elapsedLabel := fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0)
	elapsedLabel.SetLabel("0.0s")
	elapsedLabel.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)

	row = fltk_go.NewFlex(0, 0, 0, 0)
	row.SetType(fltk_go.ROW)
	label = fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0, "Duration:")
	label.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
	row.Fixed(label, WIDGET_LABEL_WIDTH)
	durationSlider := fltk_go.NewSlider(0, 0, 0, 0)
	durationSlider.SetType(fltk_go.HOR_SLIDER)
	durationSlider.SetValue(DURATION_DEFAULT)
	durationSlider.SetMaximum(DURATION_MAXIMUM)
	// durationSlider.SetCallbackCondition(fltk_go.WhenChanged)
	durationSlider.SetCallback(func() {
		// log.Printf("change duration: %f", durationSlider.Value())
		elapsedProgess.SetMaximum(durationSlider.Value())
	})
	row.End()

	resetBtn := fltk_go.NewButton(0, 0, 0, 0)
	resetBtn.SetLabel("Reset")
	resetBtn.SetCallback(func() {
		startTime = time.Now()
		elapsedProgess.SetValue(0.0)
	})

	stopCh := make(chan struct{})
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			select {
			case <-stopCh:
				return
			case t := <-ticker.C:
				fltk_go.Awake(func() {
					if durationSlider.Value()-elapsedProgess.Value() >= 0 {
						d := t.Sub(startTime)
						elapsedProgess.SetValue(d.Seconds())
						elapsedLabel.SetLabel(d.String())
					}
				})
			}
		}
	}()

	col.End()
	win.End()
	win.Show()
	fltk_go.Run()

	close(stopCh)
}
