<div align="center">

# 1. Document

</div>

<div align="center">

[Document](./README.md) | [中文文档](./README_zh-cn.md)  | [示例](./examples/README.md)

</div>

<!-- TOC -->

- [1. Document](#1-document)
- [2. fltk\_go来源](#2-fltk_go来源)
- [3. 使用](#3-使用)
	- [3.1. 依赖](#31-依赖)
	- [3.2. 使用样例](#32-使用样例)
- [4. 资源](#4-资源)

<!-- /TOC -->

---
# 2. fltk_go来源
* 从 [pwiecz/go-fltk](https://github.com/pwiecz/go-fltk) fork 以 commit hash `5313f8a5a643c8b4f71dabd084cefb9437daa8a7` 为基础变基修改
* 一个围绕 FLTK 1.4 库的简单封装，FLTK 是一个轻量级 GUI 库，允许创建小型、独立且快速的 GUI 应用程序。

# 3. 使用
## 3.1. 依赖
* 要构建 `fltk_go`，除了 `Golang编译器`，你还需要一个 `C++11 编译器`，
	*	`Linux` 上的 `GCC` 或 `Clang`
	*	`Windows` 上的 `MinGW64`
	*	`MacOS` 上的 `XCode`。

* `fltk_go` 带有一些架构的预构建 `FLTK` 库（`linux/amd64`, `windows/amd64`），但你也可以轻松地自己重建它们，或者为其他架构构建它们。
要为你的平台构建 `FLTK` 库，只需从 `fltk_go` 源代码树的根目录运行 go generate。

*	要运行使用 fltk_go 构建的程序，你将需要一些系统库，这些库通常在带有图形用户界面的操作系统上是可用的：

- Windows: 除了 `mingw64` 没有外部依赖 ([建议使用msys2的mingw64](./scripts/install_msys2_mingw64.sh))
- MacOS: 没有外部依赖
- Linux（和其他未测试的 Unix 系统）: 你需要：
    - X11
    - Xrender
    - Xcursor
    - Xfixes
    - Xext
    - Xft
    - Xinerama
    - OpenGL

## 3.2. 使用样例
![example show case)](./examples.md)

# 4. 资源
- [官方 FLTK 1.4 文档](https://www.fltk.org/doc-1.4/index.html)
- [fltk_go 文档](https://pkg.go.dev/github.com/george012/fltk_go) 