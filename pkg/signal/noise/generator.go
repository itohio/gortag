package noise

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gortag/pkg/signal/base"
)

const Name = "Noise"

type Generator struct {
	base.GeneratorBase

	noiseType binding.Int
	sigma     binding.Float
	minFreq   binding.Float
	maxFreq   binding.Float
}

func New() *Generator {
	ret := &Generator{
		noiseType: binding.NewInt(),
		sigma:     binding.NewFloat(),
		minFreq:   binding.NewFloat(),
		maxFreq:   binding.NewFloat(),
	}
	ret.Init()
	ret.sigma.Set(.01)
	ret.minFreq.Set(20)
	ret.maxFreq.Set(20000)
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
		{Name: "Type", Value: g.noiseType, Min: 0, Max: 1},
		{Name: "Sigma", Value: g.sigma, Min: 1e-6, Max: 1000, Format: "%0.01f"},
		{Name: "Min Freq, Hz", Value: g.minFreq, Min: 10, Max: 44000, Format: "%0.1f"},
		{Name: "Max Freq, Hz", Value: g.maxFreq, Min: 10, Max: 44000, Format: "%0.1f"},
	}
}
