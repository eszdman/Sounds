package setting

const (
	PianoY       = 30
	PianoX       = 100
	PianoX2      = 150
	Roll         = 16
	PianoOctaves = 7
	PianoCount   = PianoOctaves * 12
)

type PianoRollParameters struct {
	PianoY, PianoX, PianoX2, Roll, RollMpy, RollSubMpy int32
	PianoOctaves, PianoCount                           int32
}

var CurrentParameters PianoRollParameters

func UseDefaultPianoRoll() {
	CurrentParameters.PianoY = PianoY
	CurrentParameters.PianoX = PianoX
	CurrentParameters.PianoX2 = PianoX2
	CurrentParameters.Roll = Roll
	CurrentParameters.RollMpy = 4
	CurrentParameters.RollSubMpy = 4
	CurrentParameters.PianoOctaves = PianoOctaves
	CurrentParameters.PianoCount = PianoCount
}
