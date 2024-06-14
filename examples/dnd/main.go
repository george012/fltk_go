package main

import (
	"fmt"
	"github.com/george012/fltk_go"
)

type Sender struct {
	*fltk_go.Box
}

func (s *Sender) OnEvent(event fltk_go.Event) bool {
	switch event {
	case fltk_go.PUSH:
		fltk_go.CopyToSelectionBuffer("It works!")
		fltk_go.DragAndDrop()
		return true
	}
	return false
}
func NewSender(x, y, w, h int) *Sender {
	s := &Sender{}
	s.Box = fltk_go.NewBox(fltk_go.FLAT_BOX, x, y, w, h, "Drag\nfrom\nhere..")
	s.Box.SetEventHandler(s.OnEvent)
	return s
}

type Receiver struct {
	*fltk_go.Box
}

func (r *Receiver) OnEvent(event fltk_go.Event) bool {
	switch event {
	case fltk_go.DND_ENTER, fltk_go.DND_DRAG, fltk_go.DND_RELEASE:
		return true
	case fltk_go.PASTE:
		r.SetLabel(fltk_go.EventText())
		fmt.Println("Pasted '", fltk_go.EventText(), "'")
		return true
	}
	return false
}
func NewReceiver(x, y, w, h int) *Receiver {
	r := &Receiver{}
	r.Box = fltk_go.NewBox(fltk_go.FLAT_BOX, x, y, w, h, "..to\nhere")
	r.SetEventHandler(r.OnEvent)
	return r
}

type SenderWindow struct {
	*fltk_go.Window
	sender *Sender
}

func NewSenderWindow() *SenderWindow {
	w := &SenderWindow{}
	w.Window = fltk_go.NewWindow(200, 100)
	w.SetPosition(0, 0)
	w.SetLabel("Sender")
	w.sender = NewSender(0, 0, 100, 100)
	w.End()
	return w
}

type ReceiverWindow struct {
	*fltk_go.Window
	receiver *Receiver
}

func NewReceiverWindow() *ReceiverWindow {
	w := &ReceiverWindow{}
	w.Window = fltk_go.NewWindow(200, 100)
	w.SetPosition(400, 0)
	w.SetLabel("Receiver")
	w.receiver = NewReceiver(100, 0, 100, 100)
	w.End()
	return w
}
func main() {
	win_a := NewSenderWindow()
	win_a.Show()
	win_b := NewReceiverWindow()
	win_b.Show()
	fltk_go.Run()
}
