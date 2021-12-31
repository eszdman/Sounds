package parameters

import "math"

type VibratoParameters struct {
	Start     float32
	Pitch     float32
	Frequency float32
}

func (parameters *VibratoParameters) GetVibrato() func(x float32) float32 {
	return func(x float32) float32 {
		return float32(math.Sin(float64(x/20.0))) * 10
	}
}
