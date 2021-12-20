package env

import (
	"github.com/eszdman/Sounds/ui/platform"
	"time"
)

var PianoSettings = false
var NewPlatform *platform.Platform
var VoiceNotes []VoiceNote
var FPS = int32(180)
var Ticker = time.NewTicker(time.Second / time.Duration(FPS))

func Init() {
	VoiceNotes = make([]VoiceNote, 0)
}
