package main

import (
	"github.com/george012/fltk_go"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
)

// A port of fltk CubeView test program to Go

type MainWindow struct {
	*fltk_go.Window
	vrot *fltk_go.Roller
	ypan *fltk_go.Slider
	hrot *fltk_go.Roller
	xpan *fltk_go.Slider
	zoom *fltk_go.ValueSlider
	cube *CubeView
}
type CubeView struct {
	*fltk_go.GlWindow
	initialized    bool
	xshift, yshift float32
	hAng, vAng     float32
	size           float32
}

func main() {
	runtime.LockOSThread()
	// Disable screen scaling, as we don't handle it well.
	for i := 0; i < fltk_go.ScreenCount(); i++ {
		fltk_go.SetScreenScale(i, 1.0)
	}
	fltk_go.SetKeyboardScreenScaling(false)
	win := &MainWindow{}
	win.Window = fltk_go.NewWindow(415, 405)
	win.SetBox(fltk_go.UP_BOX)
	win.SetLabelSize(12)
	o := fltk_go.NewGroup(5, 3, 374, 399)
	{
		vchange := fltk_go.NewGroup(5, 100, 37, 192)
		win.vrot = fltk_go.NewRoller(5, 100, 17, 186)
		win.vrot.SetMinimum(-180)
		win.vrot.SetMaximum(180)
		win.vrot.SetStep(1)
		win.vrot.SetAlign(fltk_go.ALIGN_WRAP)
		win.vrot.SetCallback(win.OnVerticalRotation)
		win.ypan = fltk_go.NewSlider(25, 100, 17, 186)
		win.ypan.SetType(4)
		win.ypan.SetMinimum(-25)
		win.ypan.SetMaximum(25)
		win.ypan.SetStep(0.1)
		win.ypan.SetAlign(fltk_go.ALIGN_CENTER)
		win.ypan.SetCallback(win.OnVerticalPan)
		vchange.End()
	}
	{
		hchange := fltk_go.NewGroup(120, 362, 190, 40)
		win.xpan = fltk_go.NewSlider(122, 364, 186, 17)
		win.xpan.SetType(5)
		win.xpan.SetMinimum(25)
		win.xpan.SetMaximum(-25)
		win.xpan.SetStep(0.1)
		win.xpan.SetAlign(fltk_go.ALIGN_CENTER | fltk_go.ALIGN_INSIDE)
		win.xpan.SetCallback(win.OnHorizontalPan)
		win.hrot = fltk_go.NewRoller(122, 383, 186, 17)
		win.hrot.SetType(1)
		win.hrot.SetMinimum(-180)
		win.hrot.SetMaximum(180)
		win.hrot.SetStep(1)
		win.hrot.SetAlign(fltk_go.ALIGN_RIGHT)
		win.hrot.SetCallback(win.OnHorizontalRotation)
		hchange.End()
	}
	{
		mainview := fltk_go.NewGroup(46, 27, 333, 333)
		fltk_go.NewBox(fltk_go.DOWN_FRAME, 46, 27, 333, 333)
		win.cube = &CubeView{size: 10}
		win.cube.GlWindow = fltk_go.NewGlWindow(48, 29, 329, 329, win.cube.draw)
		win.cube.SetBox(fltk_go.NO_BOX)
		win.cube.SetAlign(fltk_go.ALIGN_CENTER | fltk_go.ALIGN_INSIDE)
		win.cube.SetEventHandler(func(event fltk_go.Event) bool { return win.cube.handleEvent(event) })
		win.cube.SetMode(fltk_go.ALPHA | fltk_go.DOUBLE | fltk_go.MULTISAMPLE)
		mainview.End()
		o.Resizable(mainview)
	}
	win.zoom = fltk_go.NewValueSlider(106, 3, 227, 19, "Zoom")
	win.zoom.SetType(5)
	win.zoom.SetLabelFont(1)
	win.zoom.SetLabelSize(12)
	win.zoom.SetMinimum(1)
	win.zoom.SetMaximum(50)
	win.zoom.SetStep(0.1)
	win.zoom.SetValue(10)
	win.zoom.SetTextFont(fltk_go.HELVETICA_BOLD)
	win.zoom.SetAlign(fltk_go.ALIGN_LEFT)
	win.zoom.SetCallback(win.OnZoom)
	o.End()
	win.End()
	win.Resizable(win)
	fltk_go.Lock()
	win.Show()
	fltk_go.Run()

}

func (c *CubeView) draw() {
	if !c.Valid() {
		gl.LoadIdentity()
		gl.Viewport(0, 0, int32(c.W()), int32(c.H()))
		gl.Ortho(-10, 10, -10, 10, -20050, 10000)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.PushMatrix()
	gl.Translatef(c.xshift, c.yshift, 0)
	gl.Rotatef(c.hAng, 0, 1, 0)
	gl.Rotatef(c.vAng, 1, 0, 0)
	gl.Scalef(c.size, c.size, c.size)
	c.drawCube()
	gl.PopMatrix()
}

var boxv0 [3]float32 = [3]float32{-0.5, -0.5, -0.5}
var boxv1 [3]float32 = [3]float32{0.5, -0.5, -0.5}
var boxv2 [3]float32 = [3]float32{0.5, 0.5, -0.5}
var boxv3 [3]float32 = [3]float32{-0.5, 0.5, -0.5}
var boxv4 [3]float32 = [3]float32{-0.5, -0.5, 0.5}
var boxv5 [3]float32 = [3]float32{0.5, -0.5, 0.5}
var boxv6 [3]float32 = [3]float32{0.5, 0.5, 0.5}
var boxv7 [3]float32 = [3]float32{-0.5, 0.5, 0.5}

const alpha = 0.5

func (c *CubeView) drawCube() {
	gl.ShadeModel(gl.FLAT)
	gl.Begin(gl.QUADS)
	gl.Vertex3fv(&boxv0[0])
	gl.Vertex3fv(&boxv1[0])
	gl.Vertex3fv(&boxv2[0])
	gl.Vertex3fv(&boxv3[0])

	gl.Color4f(1.0, 1.0, 0.0, alpha)
	gl.Vertex3fv(&boxv0[0])
	gl.Vertex3fv(&boxv4[0])
	gl.Vertex3fv(&boxv5[0])
	gl.Vertex3fv(&boxv1[0])

	gl.Color4f(0.0, 1.0, 1.0, alpha)
	gl.Vertex3fv(&boxv2[0])
	gl.Vertex3fv(&boxv6[0])
	gl.Vertex3fv(&boxv7[0])
	gl.Vertex3fv(&boxv3[0])

	gl.Color4f(1.0, 0.0, 0.0, alpha)
	gl.Vertex3fv(&boxv4[0])
	gl.Vertex3fv(&boxv5[0])
	gl.Vertex3fv(&boxv6[0])
	gl.Vertex3fv(&boxv7[0])

	gl.Color4f(1.0, 0.0, 1.0, alpha)
	gl.Vertex3fv(&boxv0[0])
	gl.Vertex3fv(&boxv3[0])
	gl.Vertex3fv(&boxv7[0])
	gl.Vertex3fv(&boxv4[0])

	gl.Color4f(0.0, 1.0, 0.0, alpha)
	gl.Vertex3fv(&boxv1[0])
	gl.Vertex3fv(&boxv5[0])
	gl.Vertex3fv(&boxv6[0])
	gl.Vertex3fv(&boxv2[0])
	gl.End()

	gl.Color3f(1.0, 1.0, 1.0)
	gl.Begin(gl.LINES)
	gl.Vertex3fv(&boxv0[0])
	gl.Vertex3fv(&boxv1[0])

	gl.Vertex3fv(&boxv1[0])
	gl.Vertex3fv(&boxv2[0])

	gl.Vertex3fv(&boxv2[0])
	gl.Vertex3fv(&boxv3[0])

	gl.Vertex3fv(&boxv3[0])
	gl.Vertex3fv(&boxv0[0])

	gl.Vertex3fv(&boxv4[0])
	gl.Vertex3fv(&boxv5[0])

	gl.Vertex3fv(&boxv5[0])
	gl.Vertex3fv(&boxv6[0])

	gl.Vertex3fv(&boxv6[0])
	gl.Vertex3fv(&boxv7[0])

	gl.Vertex3fv(&boxv7[0])
	gl.Vertex3fv(&boxv4[0])

	gl.Vertex3fv(&boxv0[0])
	gl.Vertex3fv(&boxv4[0])

	gl.Vertex3fv(&boxv1[0])
	gl.Vertex3fv(&boxv5[0])

	gl.Vertex3fv(&boxv2[0])
	gl.Vertex3fv(&boxv6[0])

	gl.Vertex3fv(&boxv3[0])
	gl.Vertex3fv(&boxv7[0])
	gl.End()
}

func (c *CubeView) handleEvent(event fltk_go.Event) bool {
	switch event {
	case fltk_go.SHOW:
		if !c.initialized && c.IsShown() {
			c.MakeCurrent()
			if err := gl.Init(); err != nil {
				log.Fatal("Cannot initialize OpenGL", err)
			}
			c.initialized = true
		}
		c.Redraw()
	}
	return false
}

func (w *MainWindow) OnVerticalRotation() {
	w.cube.vAng = float32(w.vrot.Value())
	w.cube.Redraw()
}
func (w *MainWindow) OnHorizontalRotation() {
	w.cube.hAng = float32(w.hrot.Value())
	w.cube.Redraw()
}
func (w *MainWindow) OnVerticalPan() {
	w.cube.yshift = float32(w.ypan.Value())
	w.cube.Redraw()
}
func (w *MainWindow) OnHorizontalPan() {
	w.cube.xshift = float32(w.xpan.Value())
	w.cube.Redraw()
}
func (w *MainWindow) OnZoom() {
	w.cube.size = float32(w.zoom.Value())
	w.cube.Redraw()
}
