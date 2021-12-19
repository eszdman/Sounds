package forms

import (
	"github.com/inkyblackness/imgui-go/v4"
)

const (
	PianoY       = 30
	PianoX       = 100
	PianoX2      = 150
	Roll         = 10
	PianoOctaves = 7
	PianoCount   = PianoOctaves * 12
)

type PianoRollParameters struct {
	PianoY, PianoX, PianoX2, Roll int32
	PianoOctaves, PianoCount      int32
}

func UseDefaultPianoRoll() {
	currentParameters.PianoY = PianoY
	currentParameters.PianoX = PianoX
	currentParameters.PianoX2 = PianoX2
	currentParameters.Roll = Roll
	currentParameters.PianoOctaves = PianoOctaves
	currentParameters.PianoCount = PianoCount
}
func PianoRollSettings(keep *bool) {
	if !imgui.BeginV("PianoRollSettings", keep, imgui.WindowFlagsNone) {
		imgui.End()
		return
	}

	imgui.SliderInt("Piano X", &currentParameters.PianoX, 1, 400)
	imgui.SliderInt("Piano Y", &currentParameters.PianoY, 1, 400)
	imgui.SliderInt("Piano X2", &currentParameters.PianoX2, 1, 800)
	imgui.InputInt("Roll", &currentParameters.Roll)
	imgui.SliderInt("PianoOctaves", &currentParameters.PianoOctaves, 1, 10)
	currentParameters.PianoCount = currentParameters.PianoOctaves * 12
	imgui.End()
}
