package sinesweep

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gortag/pkg/signal/base"
)

const Name = "Log Sine Sweep"

type Generator struct {
	base.GeneratorBase

	minFreq  binding.Float
	maxFreq  binding.Float
	duration binding.Float
}

func New() *Generator {
	ret := &Generator{
		minFreq:  binding.NewFloat(),
		maxFreq:  binding.NewFloat(),
		duration: binding.NewFloat(),
	}
	ret.Init()
	ret.minFreq.Set(20)
	ret.maxFreq.Set(20000)
	ret.duration.Set(5)
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
		{Name: "Amplitude, dB", Value: g.Amplitude},
		{Name: "Min Freq, Hz", Value: g.minFreq, Min: 10, Max: 44000, Format: "%0.1f"},
		{Name: "Max Freq, Hz", Value: g.maxFreq, Min: 10, Max: 44000, Format: "%0.1f"},
		{Name: "Duration, s", Value: g.duration, Min: .01, Max: 30, Format: "%0.1f"},
	}
}
