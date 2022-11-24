package triangle

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gortag/pkg/signal/base"
	"github.com/itohio/gortag/pkg/signal/rectangle"
)

const Name = "Triangle"

type Generator struct {
	rectangle.Generator

	ratio binding.Float
}

func New() *Generator {
	ret := &Generator{
		Generator: *rectangle.New(),
		ratio:     binding.NewFloat(),
	}
	ret.Init()
	ret.ratio.Set(0)
	return ret
}

func (g *Generator) Name() string {
	return Name
}

func (g *Generator) Reset() {
	g.GeneratorBase.Reset()
}

func (g *Generator) Parameters() []base.NamedParameter {
	rpa := g.Generator.Parameters()
	params := make([]base.NamedParameter, len(rpa)+1)
	copy(params, rpa)
	params[len(rpa)] = base.NamedParameter{
		Name:  "Ratio, %",
		Value: g.ratio,
		Min:   0, Max: 100, Format: "%0.1f",
	}
	return params
}
