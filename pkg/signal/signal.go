package signal

import (
	"errors"

	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gocoustics/pkg/audio/system"
	"github.com/itohio/gortag/pkg/signal/base"
)

// Generator represents every generator
type Generator interface {
	// Name returns the name of the generator
	Name() string
	// Active returns whether this generator is active or not
	Active() binding.Bool

	// Channels returns a Bit array of the channels this generator outputs
	Channels() binding.IntList
	//
	Parameters() []base.NamedParameter

	// FillBuffer is called by the generator engine
	// with a valid buffer. The generator should fill this buffer according
	// to internal generator logic.
	// Modulation argument specifies aggregated modulators output.
	// Internal generator logic should decide what parameters to modulate.
	// Modulation argument may be nil.
	// Length, Channel count and Sample Rate of buf and modulations will
	// always be the same.
	FillBuffer(buf system.SampleBuffer[float64], modulation system.SampleBuffer[float64])

	// Reset will be called after every generator start.
	Reset()
	// Load(map[string]interface{})
	// Save() map[string]interface{}
}

func New(name string) (Generator, error) {
	maker, ok := signals[name]
	if !ok {
		return nil, errors.New("no such generator")
	}
	return maker(), nil
}

func Signals() []string {
	return signalNames
}
