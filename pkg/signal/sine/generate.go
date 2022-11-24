package sine

import (
	"github.com/itohio/gocoustics/pkg/audio/system"
	"github.com/itohio/gocoustics/pkg/dsp"
)

func (g *Generator) FillBuffer(buf system.SampleBuffer[float64], modulation system.SampleBuffer[float64]) {
	arr, err := g.Channels().Get()
	if err != nil {
		return
	}
	amp, _ := g.Amplitude.Get()
	chNum := buf.Channels()
	sr := buf.SampleRate()
	freq, _ := g.frequency.Get()
	for _, ch := range arr {
		if ch < 1 || ch > chNum {
			continue
		}
		dsp.AddSin(buf.Buffer(), amp, g.Phase, freq, ch, chNum-1, sr)
	}

	g.Phase += buf.Duration().Seconds()
}
