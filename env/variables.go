package env

import (
	"github.com/eszdman/Sounds/engine/preview"
	"github.com/eszdman/Sounds/ui/platform"
	"github.com/eszdman/Sounds/ui/wrapper"
)

var NewPlatform *platform.Platform
var ImguiWrapping *wrapper.ImguiWrapping
var VoiceNotes []VoiceNote

func Init() {
	VoiceNotes = make([]VoiceNote, 0)
	preview.Init()
}

func DeInit() {
	preview.DeInit()
}
