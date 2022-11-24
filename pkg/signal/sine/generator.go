package sine

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gortag/pkg/signal/base"
)

const Name = "Sine"

type Generator struct {
	base.GeneratorBase
	frequency binding.Float
}

func New() *Generator {
	ret := &Generator{
		frequency: binding.NewFloat(),
	}
	ret.Init()
	ret.frequency.Set(1000)

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
		{Name: "Frequency, Hz", Value: g.frequency, Min: 10, Max: 44000, Format: "%0.1f"},
	}
}
