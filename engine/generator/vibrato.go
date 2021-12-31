package generator

import (
	"github.com/eszdman/Sounds/engine/parameters"
	"github.com/eszdman/Sounds/env"
)

func FillNote() env.VoiceNote {
	return env.VoiceNote{
		RollStart:     0,
		RollEnd:       0,
		RollPitch:     0,
		Lyrics:        "",
		VibratoParams: parameters.VibratoParameters{},
	}
}
