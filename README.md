<!-- TOC -->

- [1. fltk\_go source](#1-fltk_go-source)
- [2. Usage](#2-usage)
	- [2.1. Dependencies](#21-dependencies)
	- [2.2. Usage](#22-usage)
	- [2.3. Styles](#23-styles)
	- [2.4. Image support](#24-image-support)
- [3. Resources](#3-resources)

<!-- /TOC -->

---
* [Document](./README.md) | [中文文档](./README_zh-cn.md)

# 1. fltk_go source
* Forked from [pwiecz/go-fltk](https://github.com/pwiecz/go-fltk) with commit hash `5313f8a5a643c8b4f71dabd084cefb9437daa8a7` rebased
A simple wrapper around the FLTK 1.4 library, a lightweight GUI library that allows creating small, standalone and fast GUI applications.
# 2. Usage
## 2.1. Dependencies
* To build `fltk_go`, in addition to the `Golang compiler`, you also need a `C++11 compiler`,
* `GCC` or `Clang` on `Linux`
* `MinGW64` on `Windows`
* `XCode` on `MacOS`.

* `fltk_go` comes with prebuilt `FLTK` libraries for some architectures (`linux/amd64`, `windows/amd64`), but you can easily rebuild them yourself, or build them for other architectures.
To build the `FLTK` library for your platform, just run go generate from the root of the `fltk_go` source tree.

* To run programs built with fltk_go, you will need some system libraries that are typically available on operating systems with a graphical user interface:

- Windows: no external dependencies except `mingw64` ([msys2's mingw64 is recommended](./scripts/install_msys2_mingw64.sh))

- MacOS: no external dependencies
- Linux (and other untested Unix systems): you will need:
- X11
- Xrender
- Xcursor
- Xfixes
- Xext
- Xft
- Xinerama
- OpenGL

## 2.2. Usage
* You can use the `fltk_go.New<WidgetType>` function to create widgets and make modifications to the widget you are instantiating.
Function and method names are similar to the original C++ names, but follow the Go language's PascalCase naming convention.

Setter methods are prefixed with `Set`.

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

## 2.3. Styles
FLTK provides 4 built-in styles:
- base (default)
- gtk+
- gleam
- plastic
For example, you can use `fltk_go.SetScheme("gtk+")` to set these styles.

FLTK Also allows customizing the style of the widget:
```go
package main
import (
"strconv"

"github.com/george012/fltk_go"
)
// FLTK uses RGBI color representation, where I is an index into the FLTK color table
// Passing 00 as I will use RGB values
const GRAY = 0x75757500
const LIGHT_GRAY = 0xeeeeee00
const BLUE = 0x42A5F500
const SEL_BLUE = 0x2196F300
const WIDTH = 600
const HEIGHT = 400

func main() {
curr := 0
fltk.InitStyles()
win := fltk.NewWindow(WIDTH, HEIGHT)
win.SetLabel("Flutter-like")
win.SetColor(fltk.WHITE)
bar := fltk.NewBox(fltk.FLAT_BOX, 0, 0, WIDTH, 60, " FLTK App!")
bar.SetDrawHandler(func() { // Shadow under the bar
fltk.DrawBox(fltk.FLAT_BOX, 0, 0, WIDTH, 63, LIGHT_GRAY)
})
bar.SetAlign(fltk.ALIGN_INSIDE | fltk.ALIGN_LEFT)
bar.SetLabelColor(255) // this uses the index into the color map, here it's white
bar.SetColor(BLUE)
bar.SetLabelSize(22)
text := fltk.NewBox(fltk.NO_BOX, 250, 180, 100, 40, "You have pushed the button this many times:")
text.SetLabelSize(18)
text.SetLabelFont(fltk.TIMES)
count := fltk.NewBox(fltk.NO_BOX, 250, 180+40, 100, 40, "0")
count.SetLabelSize(36)
count.SetLabelColor(GRAY)
btn := fltk.NewButton(WIDTH-100, HEIGHT-100, 60, 60, "@+6plus") // This translates to a plus sign
btn.SetColor(BLUE)
btn.SetSelectionColor(SEL_BLUE)
btn.SetLabelColor(255)
btn.SetBox(fltk.OFLAT_BOX)
btn.ClearVisibleFocus()
btn.SetCallback(func() {
curr += 1
count.SetLabel(strconv.Itoa(curr))
})
win.End()
win.Show()
fltk.Run()
}
```

![image](https://user-images.githubusercontent.com/37966791/147374840-2d993522-fc86-46fc-9e95-2b3391d31013.png)

Label properties can be viewed [here](https://www.fltk.org/doc-1.3/common.html#common_labels)

## 2.4. Image support
FLTK supports both vector and raster graphics, through several image types:
- SvgImage
- RgbImage
- JpegImage
- PngImage
- BmpImage
- SharedImage

Some of these can be instantiated from image files or data:
```go
package main

import (
"fmt"

"github.com/george012/fltk_go"

func main() {
win := fltk.NewWindow(400, 300)
box := fltk.NewBox(fltk.FLAT_BOX, 0, 0, 400, 300, "")
image, err := fltk.NewJpegImageLoad("image.jpg")
if err != nil {
fmt.Printf("An error occurred: %s\n", err)
} else {
box.SetImage(image)
}
win.End()
win.Show()
fltk.Run()
}
```

# 3. Resources
- [Official FLTK 1.4 Documentation](https://www.fltk.org/doc-1.4/index.html)
- [fltk_go Documentation](https://pkg.go.dev/github.com/george012/fltk_go)