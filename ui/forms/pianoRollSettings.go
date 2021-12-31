package forms

import (
	"github.com/eszdman/Sounds/setting"
	"github.com/inkyblackness/imgui-go/v4"
)

var usePianoRollSettings = false

func PianoRollSettings() {
	if !usePianoRollSettings {
		return
	}
	if !imgui.BeginV("PianoRollSettings", &usePianoRollSettings, imgui.WindowFlagsAlwaysAutoResize) {
		imgui.End()
		return
	}

	imgui.SliderInt("Piano X", &setting.CurrentParameters.PianoX, 1, 400)
	imgui.SliderInt("Piano Y", &setting.CurrentParameters.PianoY, 1, 400)
	imgui.SliderInt("Piano X2", &setting.CurrentParameters.PianoX2, 1, 800)
	imgui.InputInt("Roll", &setting.CurrentParameters.Roll)
	imgui.SliderInt("PianoOctaves", &setting.CurrentParameters.PianoOctaves, 1, 10)
	setting.CurrentParameters.PianoCount = setting.CurrentParameters.PianoOctaves * 12
	imgui.End()
}
