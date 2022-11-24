package ui

import (
	"errors"
	"strconv"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/gortag/pkg/signal"
	"github.com/itohio/gortag/pkg/signal/base"
)

var _ fyne.Widget = (*signalWidget)(nil)
var _ Signal = &signalWidget{}
var _ fyne.CanvasObject = &signalWidget{}

type signalWidget struct {
	widget.BaseWidget
	content *fyne.Container
	header  *fyne.Container

	name *widget.Label
	gen  signal.Generator

	active         atomic.Value // binding.Bool
	activeListener binding.DataListener
	open           bool
	unbinds        []func()
}

func NewSignal(onDelete func(signal.Generator)) *signalWidget {
	ret := &signalWidget{
		name: widget.NewLabel("-"),
		open: true,
	}
	ret.ExtendBaseWidget(ret)

	deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if onDelete != nil {
			onDelete(ret.gen)
		}
	})
	startBtn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if active := ret.active.Load(); active != nil {
			bound := active.(binding.Bool)
			current, _ := bound.Get()
			bound.Set(!current)
		}
	})
	openBtn := widget.NewButtonWithIcon("", theme.MenuDropUpIcon(), nil)
	openBtn.OnTapped = func() {
		ret.open = !ret.open
		if ret.open {
			openBtn.SetIcon(theme.MenuDropUpIcon())
			ret.content.Show()
		} else {
			openBtn.SetIcon(theme.MenuDropDownIcon())
			ret.content.Hide()
		}
	}

	ret.header = container.NewBorder(nil, nil,
		deleteBtn, container.NewHBox(openBtn, startBtn),
		ret.name,
	)

	ret.content = container.New(
		layout.NewFormLayout(),
	)

	ret.activeListener = binding.NewDataListener(func() {
		if active := ret.active.Load(); active != nil {
			bound := active.(binding.Bool)
			if playing, _ := bound.Get(); playing {
				startBtn.SetIcon(theme.MediaStopIcon())
			} else {
				startBtn.SetIcon(theme.MediaPlayIcon())
			}
		}
	})

	return ret
}

func (s *signalWidget) UpdateGenerator(gen signal.Generator) {
	if s.gen == gen {
		return
	}
	s.gen = gen

	if active := s.active.Load(); active != nil {
		bound := active.(binding.Bool)
		bound.RemoveListener(s.activeListener)
	}

	for _, u := range s.unbinds {
		u()
	}

	if gen == nil {
		return
	}

	s.active.Store(gen.Active())
	gen.Active().AddListener(s.activeListener)
	s.name.SetText(gen.Name())

	objects := s.content.Objects[:0]
	unbinds := s.unbinds[:0]
	// objects = append(objects, widget.NewLabel("Channels"))
	// objects = append(objects, widget.)

	for _, param := range gen.Parameters() {
		objects = append(objects, widget.NewLabel(param.Name))

		switch b := param.Value.(type) {
		case binding.Float:
			o, u := makeSlider(param, b)
			objects = append(objects, o)
			unbinds = append(unbinds, u)
		case binding.Int:
			bound := binding.IntToStringWithFormat(b, param.Format)
			entry := widget.NewEntryWithData(bound)
			entry.Validator = func(s string) error {
				if v, err := strconv.ParseInt(s, 10, 32); err != nil {
					return err
				} else {
					if v > int64(param.Max) {
						return errors.New("too large")
					}
					if v < int64(param.Min) {
						return errors.New("too small")
					}
				}
				return nil
			}
			objects = append(objects, entry)
			unbinds = append(unbinds, func() {
				entry.Unbind()
				// bind.RemoveListener(bind) FIXME: definite memory leak
			})
		case binding.Bool:
			w := widget.NewCheckWithData("", b)
			objects = append(objects, w)
			unbinds = append(unbinds, func() {
				w.Unbind()
			})
		case binding.String:
			w := widget.NewEntryWithData(b)
			objects = append(objects, w)
			unbinds = append(unbinds, func() {
				w.Unbind()
			})
		default:
			panic("unknown parameter type")
		}
	}
	s.content.Objects = objects
	s.unbinds = unbinds
	// s.content.Layout.Layout(s.content.Objects, s.content.Size())
	s.Refresh()
}

func makeSlider(param base.NamedParameter, val binding.Float) (*fyne.Container, func()) {
	bound := binding.FloatToStringWithFormat(val, param.Format)
	entry := widget.NewEntryWithData(bound)
	slider := widget.NewSliderWithData(param.Min, param.Max, val)
	entry.Validator = func(s string) error {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		if f > param.Max {
			return errors.New("too large")
		}
		if f < param.Min {
			return errors.New("too small")
		}
		return nil
	}
	unbind := func() {
		entry.Unbind()
		slider.Unbind()
		// val.RemoveListener(bound) FIXME: Definitely a memory leak!
	}

	return container.New(NewFixedLayout(-1, 100), slider, entry), unbind
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *signalWidget) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	r := &signalRenderer{
		signal: s,
	}
	return r
}

// MinSize returns the size that this widget should not shrink below.
func (s *signalWidget) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	// essentially calls renderer.MinSize
	minSize := s.BaseWidget.MinSize()
	return minSize
}
