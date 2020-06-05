package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type ScrollV struct {
	*widget.ScrollContainer
	MinHeight int
}

func (sv ScrollV) MinSize() fyne.Size {
	size := sv.ScrollContainer.MinSize()
	size.Height = sv.MinHeight
	return size
}

func newScrollWithMinHeight(content fyne.CanvasObject, minHeight int) *ScrollV {
	s := &ScrollV{
		widget.NewVScrollContainer(content),
		minHeight,
	}
	s.ExtendBaseWidget(s)
	return s
}
