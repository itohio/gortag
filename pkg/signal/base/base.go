package base

import (
	"sync"

	"fyne.io/fyne/v2/data/binding"
)

type NamedParameter struct {
	Name     string
	Value    binding.DataItem
	Min, Max float64
	Format   string
}

func (nf NamedParameter) GetFloat() float64 {
	if vf, ok := nf.Value.(binding.Float); ok {
		v, _ := vf.Get()
		return v
	}
	return 0
}

func (nf NamedParameter) GetBool() bool {
	if vf, ok := nf.Value.(binding.Bool); ok {
		v, _ := vf.Get()
		return v
	}
	return false
}

func (nf NamedParameter) GetString() string {
	if vf, ok := nf.Value.(binding.String); ok {
		v, _ := vf.Get()
		return v
	}
	return ""
}

type GeneratorBase struct {
	sync.RWMutex
	Amplitude  binding.Float
	channels   binding.IntList
	active     binding.Bool
	modulation binding.Int

	Phase float64
}

func (g *GeneratorBase) Init() {
	g.Amplitude = binding.NewFloat()
	g.channels = binding.NewIntList()
	g.active = binding.NewBool()
	g.modulation = binding.NewInt()
	g.channels.Append(1)
	g.Amplitude.Set(0.0)
}

func (g *GeneratorBase) Name() string {
	panic("must impmlement Name")
}

func (g *GeneratorBase) Channels() binding.IntList {
	return g.channels
}

func (g *GeneratorBase) Active() binding.Bool {
	return g.active
}

func (g *GeneratorBase) Modulation() binding.Int {
	return g.modulation
}

func (g *GeneratorBase) Reset() {
	g.Lock()
	defer g.Unlock()
	g.LockedReset()
	g.Phase = 0
}

func (g *GeneratorBase) LockedReset() {
	g.Phase = 0
}
