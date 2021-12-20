package forms

import (
	"github.com/eszdman/Sounds/env"
	"github.com/inkyblackness/imgui-go/v4"
	"os"
	"time"
)

var useRenderMenu = false

func RenderMenu() {
	if !useRenderMenu {
		return
	}
	defer imgui.End()
	if !imgui.BeginV("Render Settings", &useRenderMenu, imgui.WindowFlagsNone) {
		return
	}
	if env.FPS < 400 {
		imgui.SliderIntV("FPS", &env.FPS, 5, 400, "FPS: %d", imgui.SliderFlagsNone)
		env.Ticker.Reset(time.Second / time.Duration(env.FPS))
	} else {
		imgui.SliderIntV("FPS", &env.FPS, 5, 400, "FPS: UNLIMITED", imgui.SliderFlagsNone)
		env.Ticker.Reset(time.Nanosecond)
	}
}
func MenuBar(size [2]float32) {
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 0)
	defer imgui.PopStyleVar()
	defer imgui.PopStyleVar()
	defer imgui.End()
	imgui.SetNextWindowPosV(imgui.Vec2{X: 0, Y: 0}, imgui.ConditionFirstUseEver, imgui.Vec2{})
	imgui.SetNextWindowSize(imgui.Vec2{X: size[0], Y: size[1]})
	if !imgui.BeginV("Window", nil, imgui.WindowFlagsNoMove|
		imgui.WindowFlagsMenuBar|imgui.WindowFlagsNoTitleBar|
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoBringToFrontOnFocus|
		imgui.WindowFlagsNoBackground) {
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
				usePianoRollSettings = true
			}
			if imgui.MenuItemV("PrintNotes", "", false, true) {
				for i := 0; i < len(env.VoiceNotes); i++ {
					println("Note:", i, " Pitch:", env.VoiceNotes[i].RollPitch, " Start:", env.VoiceNotes[i].RollStart)
				}
			}
			if imgui.MenuItemV("Render", "", false, true) {
				useRenderMenu = true
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Tools") {
			imgui.EndMenu()
		}
		imgui.EndMenuBar()
	}
}
