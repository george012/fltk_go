package browser_test

import (
	"fmt"
	"github.com/george012/fltk_go"
	"testing"
)

func TestBrowsers(t *testing.T) {
	fltk_go.SetScheme("base")
	win := fltk_go.NewWindow(1200, 800)
	win.SetLabel("Browser simple")

	grid := fltk_go.NewGrid(0, 0, win.W(), win.H())
	grid.SetLayout(3, 2, 10, 10) // 设置3行2列，间距为10

	// 创建 Browser 示例
	browser := fltk_go.NewBrowser(0, 0, 0, 0)
	browser.Add("Browser 1")
	browser.Add("Browser 2")
	browser.Add("Browser 3")
	browser.SetCallback(func() {
		item := browser.Value()
		fmt.Printf("Browser selected: %d\n", item)
	})
	grid.SetWidget(browser, 0, 0, fltk_go.GridFill)

	// 创建 CheckBrowser 示例
	checkBrowser := fltk_go.NewCheckBrowser(0, 0, 0, 0)
	checkBrowser.Add("CheckBrowser 1", false)
	checkBrowser.Add("CheckBrowser 2", true)
	checkBrowser.Add("CheckBrowser 3", false)
	checkBrowser.SetCallback(func() {
		itemCount := checkBrowser.ItemCount()
		for i := 1; i <= itemCount; i++ {
			fmt.Printf("CheckBrowser selected %d, status: %v\n", i, checkBrowser.IsChecked(i))
		}
	})
	grid.SetWidget(checkBrowser, 0, 1, fltk_go.GridFill)

	// 创建 SelectBrowser 示例
	selectBrowser := fltk_go.NewSelectBrowser(0, 0, 0, 0)
	selectBrowser.Add("SelectBrowser 1")
	selectBrowser.Add("SelectBrowser 2")
	selectBrowser.Add("SelectBrowser 3")
	selectBrowser.SetCallback(func() {
		item := selectBrowser.Value()
		fmt.Printf("SelectBrowser is selected: %d\n", item)
	})
	grid.SetWidget(selectBrowser, 1, 0, fltk_go.GridFill)

	// 创建 HoldBrowser 示例
	holdBrowser := fltk_go.NewHoldBrowser(0, 0, 0, 0)
	holdBrowser.Add("HoldBrowser 1")
	holdBrowser.Add("HoldBrowser 2")
	holdBrowser.Add("HoldBrowser 3")
	holdBrowser.SetCallback(func() {
		item := holdBrowser.Value()
		fmt.Printf("HoldBrowser is selected: %d\n", item)
	})
	grid.SetWidget(holdBrowser, 1, 1, fltk_go.GridFill)

	// 创建 MultiBrowser 示例
	multiBrowser := fltk_go.NewMultiBrowser(0, 0, 0, 0)
	multiBrowser.Add("MultiBrowser 1")
	multiBrowser.Add("MultiBrowser 2")
	multiBrowser.Add("MultiBrowser 3")
	multiBrowser.SetCallback(func() {
		selected := ""
		for i := 1; i <= multiBrowser.Size(); i++ {
			if multiBrowser.IsSelected(i) {
				selected += fmt.Sprintf("%d ", i)
			}
		}
		fmt.Printf("MultiBrowser is selected: %s\n", selected)
	})
	grid.SetWidget(multiBrowser, 2, 0, fltk_go.GridFill)

	grid.SetShowGrid(true)
	grid.End()
	win.End()
	win.Resizable(grid)
	win.Show()
	fltk_go.Run()
}
