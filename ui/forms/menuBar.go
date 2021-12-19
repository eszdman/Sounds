package forms

import (
	"github.com/eszdman/Sounds/env"
	"github.com/inkyblackness/imgui-go/v4"
	"os"
)

func MenuBar(keep *bool, size [2]float32) {
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 0, Y: 0}, imgui.ConditionFirstUseEver, imgui.Vec2{})
	imgui.SetNextWindowSize(imgui.Vec2{X: size[0], Y: size[1]})
	if !imgui.BeginV("Window", keep, imgui.WindowFlagsNoMove|
		imgui.WindowFlagsMenuBar|imgui.WindowFlagsNoTitleBar|
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoBringToFrontOnFocus|
		imgui.WindowFlagsNoBackground) {
		// Early out if the window is collapsed, as an optimization.
		imgui.End()
		return
	}
	imgui.PushItemWidth(imgui.FontSize() * -12)
	// MenuBar
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("New") {
			}
			if imgui.MenuItemV("Open", "Ctrl+O", false, true) {
				println("Opened!")
			}
			if imgui.MenuItemV("Save", "Ctrl+S", false, true) {
				println("Saved!")
			}
			if imgui.MenuItemV("Exit", "", false, true) {
				os.Exit(1)
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Settings") {
			if imgui.MenuItemV("PianoRoll", "", false, true) {
				env.PianoSettings = true
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Tools") {
			imgui.EndMenu()
		}
		imgui.EndMenuBar()
	}
	imgui.End()
	imgui.PopStyleVar()
}
