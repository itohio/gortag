package base

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gocoustics/pkg/audio/system"
)

type GeneratorBaseWithBuffer struct {
	GeneratorBase

	Buffer    system.SampleBuffer[float64]
	Repeat    binding.Bool
	Delay     binding.Float
	onRebuild func() system.SampleBuffer[float64]
}

func (g *GeneratorBaseWithBuffer) Init(rebuildBuffer func() system.SampleBuffer[float64]) {
	g.GeneratorBase.Init()
	g.Repeat = binding.NewBool()
	g.Delay = binding.NewFloat()
	g.onRebuild = rebuildBuffer
}

func (g *GeneratorBaseWithBuffer) Reset() {
	g.Lock()
	defer g.Unlock()
	g.GeneratorBase.LockedReset()
	if g.Buffer != nil {
		g.Buffer.Release()
		g.Buffer = nil
	}
	if g.onRebuild != nil {
		g.Buffer = g.onRebuild()
	}
}

func (g *GeneratorBaseWithBuffer) FillBuffer(buf system.SampleBuffer[float64]) {

}
