package parameters

import "math"

type PitchMerging struct {
	LengthR   float32
	Frequency float32
}

func (parameters *VibratoParameters) GetMerging() func(x float32) float32 {
	return func(x float32) float32 {
		return float32(math.Sin(float64(x/20.0))) * 10
	}
}
