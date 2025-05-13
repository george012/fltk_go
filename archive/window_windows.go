package fltk_go

/*
#include "window.h"
*/
import "C"

func (w *Window) RawHandleWithWin32() uintptr {
	return uintptr(C.go_fltk_Window_win32_xid((*C.Fl_Window)(w.ptr())))
}
