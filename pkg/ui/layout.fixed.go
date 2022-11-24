package ui

import "fyne.io/fyne/v2"

var _ fyne.Layout = (*FixedLayout)(nil)

type FixedLayout struct {
	minSize []float32
}

func NewFixedLayout(minSize ...float32) *FixedLayout {
	ret := &FixedLayout{
		minSize: minSize,
	}

	return ret
}

// Layout will manipulate the listed CanvasObjects Size and Position
// to fit within the specified size.
func (l *FixedLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	sizes := make([]float32, len(objects))
	var totalSize float32
	unknown := 0
	for i := range objects {
		if i < len(l.minSize) {
			if l.minSize[i] > 0 {
				sizes[i] = l.minSize[i]
				totalSize += l.minSize[i]
				continue
			}
		}
		unknown++
	}
	if unknown > 0 {
		avgSize := (size.Width - totalSize) / float32(unknown)
		for i, s := range sizes {
			if s > 0 {
				continue
			}
			sizes[i] = avgSize
		}
	}

	var pos fyne.Position
	for i, o := range objects {
		o.Move(pos)
		size.Width = sizes[i]
		o.Resize(size)
		pos.X += sizes[i]
	}
}

// MinSize calculates the smallest size that will fit the listed
// CanvasObjects using this Layout algorithm.
func (l *FixedLayout) MinSize(objects []fyne.CanvasObject) (size fyne.Size) {
	var minSize float32
	for i, o := range objects {
		if i < len(l.minSize) {
			minSize = l.minSize[i]
		}
		min := o.MinSize()
		size.Width += fyne.Max(minSize, min.Width)
		size.Height = fyne.Max(size.Height, min.Height)
	}
	return
}
