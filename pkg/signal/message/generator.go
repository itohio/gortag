package message

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gocoustics/pkg/audio/system"
	"github.com/itohio/gortag/pkg/signal/base"
)

const Name = "FSK Encoded Message"

type Generator struct {
	base.GeneratorBaseWithBuffer

	markFreq  binding.Float
	spaceFreq binding.Float
	baudRate  binding.Int
	message   binding.String
}

func New() *Generator {
	ret := &Generator{
		markFreq:  binding.NewFloat(),
		spaceFreq: binding.NewFloat(),
		baudRate:  binding.NewInt(),
		message:   binding.NewString(),
	}
	ret.Init(func() system.SampleBuffer[float64] { return nil })
	ret.markFreq.Set(13000)
	ret.spaceFreq.Set(16000)
	ret.baudRate.Set(9600)
	ret.message.Set("Test")
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
		{Name: "Mark Freq, Hz", Value: g.markFreq, Min: 10, Max: 44000, Format: "%0.1f"},
		{Name: "Space Freq, Hz", Value: g.spaceFreq, Min: 10, Max: 44000, Format: "%0.1f"},
		{Name: "Baud", Value: g.baudRate, Min: 1, Max: 11500, Format: "%d"},
		{Name: "Message", Value: g.message},
	}
}
