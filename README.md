# go-fltk

---
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/george012/fltk_go)

* [Document](./README.md)
* [中文文档](./README_zh-cn.md)

A simple wrapper around FLTK 1.4 library, which is a lightweight GUI library which allows creating small, self-contained and fast gui applications.

## Requirements
For building go-fltk, besides Golang compiler, you will also need a C++11 compiler, such as GCC or Clang on Linux, MinGW on Windows and XCode on MacOS.

go-fltk comes with prebuilt FLTK libraries for some architectures (linux/amd64, windows/amd64), but you can easily rebuild them yourself 
or build them for other architectures.
To build FLTK libraries for your platform it should be enough to call `go generate` from the root of the go-fltk source tree.

If the build procedure doesn't work for you, you can modify `fltk-build.sh` or `fltk-build.bat` yourself or ask a question on `https://github.com/george012/fltk_go/discussions`.

For running programs built using go-fltk you will need some system libs which are normally available on operating systems with a graphical user interfaces:
- Windows: No external dependencies, besides a (for mingw64)
- MacOS: No external dependencies
- Linux (and other Unix systems - not tested): You need:
    - X11
    - Xrender
    - Xcursor
    - Xfixes
    - Xext
    - Xft
    - Xinerama
    - OpenGL

## Usage
```go
package main

import "github.com/george012/fltk_go"

func main() {
    win := fltk_go.NewWindow(400, 300)
    win.SetLabel("Main Window")
    btn := fltk_go.NewButton(160, 200, 80, 30, "Click")
    btn.SetCallback(func() {
        btn.SetLabel("Clicked")
    })
    win.End()
    win.Show()
    fltk_go.Run()
}
```

Widgets are created using the `fltk_go.New<WidgetType>` functions, modified for whatever widget you're instantiating.
Function and method names resemble the original C++ names, while however, following the Golang convention of PascalCase. 
Setter methods are also preceded by a `Set` prefix.

## Styling
FLTk offers 4 builtin schemes:
- base (default)
- gtk+
- gleam
- plastic
These can be set using `fltk_go.SetScheme("gtk+")` for example.

FLTK also allows custom styling of your widgets:
```go
package main

import (
	"strconv"

	"github.com/george012/fltk_go"
)

// FLTK uses an RGBI color representation, the I is an index into FLTK's color map
// Passing 00 as I will use the RGB part of the value
const GRAY = 0x75757500
const LIGHT_GRAY = 0xeeeeee00
const BLUE = 0x42A5F500
const SEL_BLUE = 0x2196F300
const WIDTH = 600
const HEIGHT = 400

func main() {
	curr := 0
	fltk_go.InitStyles()
	win := fltk_go.NewWindow(WIDTH, HEIGHT)
	win.SetLabel("Flutter-like")
	win.SetColor(fltk_go.WHITE)
	bar := fltk_go.NewBox(fltk_go.FLAT_BOX, 0, 0, WIDTH, 60, "    FLTK App!")
	bar.SetDrawHandler(func() { // Shadow under the bar
		fltk_go.DrawBox(fltk_go.FLAT_BOX, 0, 0, WIDTH, 63, LIGHT_GRAY)
	})
	bar.SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
	bar.SetLabelColor(255) // this uses the index into the color map, here it's white
	bar.SetColor(BLUE)
	bar.SetLabelSize(22)
	text := fltk_go.NewBox(fltk_go.NO_BOX, 250, 180, 100, 40, "You have pushed the button this many times:")
	text.SetLabelSize(18)
	text.SetLabelFont(fltk_go.TIMES)
	count := fltk_go.NewBox(fltk_go.NO_BOX, 250, 180+40, 100, 40, "0")
	count.SetLabelSize(36)
	count.SetLabelColor(GRAY)
	btn := fltk_go.NewButton(WIDTH-100, HEIGHT-100, 60, 60, "@+6plus") // this translates into a plus sign
	btn.SetColor(BLUE)
	btn.SetSelectionColor(SEL_BLUE)
	btn.SetLabelColor(255)
	btn.SetBox(fltk_go.OFLAT_BOX)
	btn.ClearVisibleFocus()
	btn.SetCallback(func() {
		curr += 1
		count.SetLabel(strconv.Itoa(curr))
	})
	win.End()
	win.Show()
	fltk_go.Run()
}
```

![image](https://user-images.githubusercontent.com/37966791/147374840-2d993522-fc86-46fc-9e95-2b3391d31013.png)

Label attributes can be seen [here](https://www.fltk.org/doc-1.3/common.html#common_labels)

## Image support
FLTK supports vector and raster graphics, via several image types:
- SvgImage
- RgbImage
- JpegImage
- PngImage
- BmpImage
- SharedImage

Some of these can be instantiated from an image file or from data:
```go
package main

import (
	"fmt"

	"github.com/george012/fltk_go"
)

func main() {
	win := fltk_go.NewWindow(400, 300)
	box := fltk_go.NewBox(fltk_go.FLAT_BOX, 0, 0, 400, 300, "")
	image, err := fltk.NewJpegImageLoad("image.jpg")
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	} else {
		box.SetImage(image)
	}
	win.End()
	win.Show()
	fltk.Run()
}
```

## Resources
- [Link](https://www.fltk.org/doc-1.4/index.html) to the official FLTK 1.4 documentation.
- [Link](https://pkg.go.dev/github.com/george012/fltk_go) to go-fltk documentation.
