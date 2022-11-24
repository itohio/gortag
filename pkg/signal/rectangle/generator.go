package rectangle

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gortag/pkg/signal/base"
)

const Name = "Rectangle"

type Generator struct {
	base.GeneratorBase

	frequency binding.Float
	duty      binding.Float
	riseTime  binding.Float
	fallTime  binding.Float
}

func New() *Generator {
	ret := &Generator{
		frequency: binding.NewFloat(),
		duty:      binding.NewFloat(),
		riseTime:  binding.NewFloat(),
		fallTime:  binding.NewFloat(),
	}
	ret.Init()
	ret.frequency.Set(1000)
	ret.duty.Set(50)
	ret.riseTime.Set(0)
	ret.fallTime.Set(0)
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
		{Name: "Duty Cycle, %", Value: g.duty, Min: 0, Max: 100, Format: "%0.1f"},
		{Name: "Rise Time, ms", Value: g.riseTime, Min: 0, Max: 10, Format: "%0.01f"},
		{Name: "Fall Time, ms", Value: g.fallTime, Min: 0, Max: 10, Format: "%0.01f"},
	}
}
