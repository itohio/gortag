package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

type signalRenderer struct {
	signal *signalWidget
}

func (r *signalRenderer) Layout(size fyne.Size) {
	h := r.signal.header.MinSize()
	c := r.signal.content.MinSize()
	r.signal.header.Resize(fyne.NewSize(size.Width, h.Height))

	if r.signal.open {
		r.signal.content.Move(fyne.NewPos(0, h.Height+theme.Padding()))
		r.signal.content.Resize(fyne.NewSize(size.Width, c.Height))
	}
}

func (r *signalRenderer) MinSize() (size fyne.Size) {
	size.Width = 300
	min := r.signal.header.MinSize()
	size.Width = fyne.Max(size.Width, min.Width)
	size.Height += min.Height
	min = r.signal.content.MinSize()
	size.Width = fyne.Max(size.Width, min.Width)
	if r.signal.open {
		size.Height += min.Height
		size.Height += theme.Padding()
	}
	return
}

func (r *signalRenderer) Refresh() {
	r.Layout(r.signal.Size())
	canvas.Refresh(r.signal)
}

func (r *signalRenderer) Destroy() {
}

func (r *signalRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.signal.header,
		r.signal.content,
	}
}
