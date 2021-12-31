package env

import (
	"github.com/eszdman/Sounds/engine/parameters"
	"github.com/eszdman/Sounds/setting"
	"math"
)

type VoiceNote struct {
	RollStart, RollEnd, RollPitch int
	Lyrics                        string
	VibratoParams                 parameters.VibratoParameters
	PitchMerging                  parameters.PitchMerging
	PathNext                      *VoiceNote
	PathPrevious                  *VoiceNote
}

func (note *VoiceNote) getFullVibrato() func(x float32) float32 {
	return func(x float32) float32 {
		return note.VibratoParams.GetVibrato()(x)
	}
}

func (note *VoiceNote) GetFrequency() (output float64) {
	output = math.Pow(2.0, float64(int(setting.CurrentParameters.PianoCount)-note.RollPitch-58)/12.0) * 440.0
	println("Frequency:", output)
	return
}
