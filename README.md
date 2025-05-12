<div align="center">

# 1. Document

</div>

<div align="center">

[Document](./README.md) | [中文文档](./README_zh-cn.md) | [examples](./examples.md)

</div>

<!-- TOC -->

- [1. Document](#1-document)
- [2. fltk\_go source](#2-fltk_go-source)
- [3. Usage](#3-usage)
	- [3.1. Dependencies](#31-dependencies)
	- [3.2. Usage Example](#32-usage-example)
- [4. Resources](#4-resources)

<!-- /TOC -->

---

# 2. fltk_go source
* Forked from [pwiecz/go-fltk](https://github.com/pwiecz/go-fltk) with commit hash `5313f8a5a643c8b4f71dabd084cefb9437daa8a7` rebased
* A simple wrapper around the FLTK 1.4 library, a lightweight GUI library that allows creating small, standalone and fast GUI applications.

# 3. Usage
## 3.1. Dependencies
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

## 3.2. Usage Example
[![example showcase](./examples.png)](./examples.png)

# 4. Resources
- [Official FLTK 1.4 Documentation](https://www.fltk.org/doc-1.4/index.html)
- [fltk_go Documentation](https://pkg.go.dev/github.com/george012/fltk_go)