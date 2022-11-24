package normalize

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gocoustics/pkg/audio/system"
	"github.com/itohio/gortag/pkg/signal/base"
)

const Name = "Normalize Output"

type Generator struct {
	base.GeneratorBaseWithBuffer

	decay binding.Float
}

func New() *Generator {
	ret := &Generator{
		decay: binding.NewFloat(),
	}
	ret.Init(func() system.SampleBuffer[float64] { return nil })
	ret.decay.Set(1)
	return ret
}

func (g *Generator) Name() string {
	return Name
}

func (g *Generator) Reset() {
	g.GeneratorBase.Reset()
}

func (g *Generator) Parameters() []base.NamedParameter {
	return []base.NamedParameter{
		{Name: "Amplitude, dB", Value: g.Amplitude, Min: -1000, Max: 10, Format: "%0.1f"},
		{Name: "Decay, s", Value: g.decay},
	}
}
