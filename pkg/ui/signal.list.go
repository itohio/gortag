package ui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/gortag/pkg/signal"
)

var _ fyne.Widget = (*signalList)(nil)

// signalList displays a list of signalWidgets.
type signalList struct {
	widget.BaseWidget
	mx        sync.RWMutex
	engine    signal.Engine
	items     []*signalWidget
	MultiOpen bool
}

// NewSignalList creates a new accordion widget.
func NewSignalList(engine signal.Engine) *signalList {

	a := &signalList{
		engine: engine,
	}
	a.ExtendBaseWidget(a)

	data := engine.Generators()
	data.AddListener(binding.NewDataListener(func() {
		boundGenerators, err := data.Get()
		if err != nil {
			panic(err)
		}

		defer a.Refresh()
		a.mx.Lock()
		defer a.mx.Unlock()
		for i, boundGen := range boundGenerators {
			if i >= len(a.items) {
				a.items = append(a.items, NewSignal(a.onDelete(data)))
			}
			gen, ok := boundGen.(signal.Generator)
			if !ok {
				panic("wrong signal type")
			}
			a.items[i].UpdateGenerator(gen)
		}

		if len(boundGenerators) < len(a.items) {
			for _, item := range a.items[len(boundGenerators):] {
				item.UpdateGenerator(nil)
			}
			a.items = a.items[:len(boundGenerators)]
		}
	}))
	return a
}

func (a *signalList) onDelete(data binding.UntypedList) func(gen signal.Generator) {
	return func(gen signal.Generator) {
		if gen == nil {
			return
		}
		items, err := data.Get()
		if err != nil {
			panic(err)
		}
		for i, item := range items {
			if item == gen {
				// deletedItem := items[i]
				items = listRemoveItem(items, i)
				err := data.Set(items)
				if err != nil {
					panic(err)
				}
				return
			}
		}
	}
}

func listRemoveItem(items []interface{}, index int) []interface{} {
	if index < 0 || index >= len(items) {
		return items
	}
	return append(items[:index], items[index+1:]...)
}

// Close collapses the item at the given index.
func (a *signalList) Close(index int) {
	if index < 0 || index >= len(a.items) {
		return
	}
	a.items[index].open = false
	a.Refresh()
}

// CloseAll collapses all items.
func (a *signalList) CloseAll() {
	for _, i := range a.items {
		i.open = false
	}
	a.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (a *signalList) CreateRenderer() fyne.WidgetRenderer {
	a.ExtendBaseWidget(a)
	r := &signalListRenderer{
		container: a,
	}
	r.updateObjects()
	return r
}

// MinSize returns the size that this widget should not shrink below.
func (a *signalList) MinSize() fyne.Size {
	a.ExtendBaseWidget(a)
	return a.BaseWidget.MinSize()
}

// Open expands the item at the given index.
func (a *signalList) Open(index int) {
	if index < 0 || index >= len(a.items) {
		return
	}
	for i, ai := range a.items {
		if i == index {
			ai.open = true
		} else if !a.MultiOpen {
			ai.open = false
		}
	}
	a.Refresh()
}

// OpenAll expands all items.
func (a *signalList) OpenAll() {
	if !a.MultiOpen {
		return
	}
	for _, i := range a.items {
		i.open = true
	}
	a.Refresh()
}

type signalListRenderer struct {
	objects   []fyne.CanvasObject
	container *signalList
	// scroller  *container.Scroll
}

func (r *signalListRenderer) Layout(size fyne.Size) {
	var pos fyne.Position
	for _, ai := range r.container.items {
		min := ai.MinSize()
		ai.Move(pos)
		ai.Resize(fyne.NewSize(size.Width, min.Height))
		pos.Y += min.Height + theme.Padding()
	}
}

func (r *signalListRenderer) MinSize() (size fyne.Size) {
	for _, ai := range r.container.items {
		min := ai.MinSize()
		size.Width = fyne.Max(size.Width, min.Width)
		size.Height += min.Height
		size.Height += theme.Padding()
	}
	return
}

func (r *signalListRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.container.Size())
	canvas.Refresh(r.container)
}

// Destroy does nothing in the base implementation.
//
// Implements: fyne.WidgetRenderer
func (r *signalListRenderer) Destroy() {
}

// Objects returns the objects that should be rendered.
//
// Implements: fyne.WidgetRenderer
func (r *signalListRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *signalListRenderer) updateObjects() {
	r.container.mx.RLock()
	defer r.container.mx.RUnlock()

	r.objects = r.objects[:0]
	for _, so := range r.container.items {
		r.objects = append(r.objects, so)
	}
}
