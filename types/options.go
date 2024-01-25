// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// RadioGroup options.

package types

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// RadioGroup type.
type RadioGroup struct {
	Options    []string `json:"options"`
	Horizontal bool     `json:"horizontal"`
	Selected   int      `json:"selected"`
}

// GetOptions returns the options of the radio group.
func (o RadioGroup) GetOptions() []string { return o.Options }

// GetHorizontal returns the horizontal type of the radio group type or false.
func (o RadioGroup) GetHorizontal() bool { return o.Horizontal }

// GetSelected returns the selected option index of the radio group.
func (o RadioGroup) GetSelected() int { return o.Selected }

// GetValue returns the selected option string of the radio group.
func (o RadioGroup) GetValue() (s string) {
	if o.Selected >= 0 && o.Selected < len(o.Options) {
		s = o.Options[o.Selected]
	}
	return
}

// GetParams returns the options, horizontal layout flag, and selected option
// string.
func (o RadioGroup) GetParams() (options []string, horizontal bool) {
	return o.GetOptions(), o.GetHorizontal()
}

// SetValue sets the selected option index of the radio group by the given
// string value.
func (o RadioGroup) SetValue(val string) RadioGroup {
	o.Selected = slices.Index(o.Options, val)
	return o
}

func (o RadioGroup) NewWidget() fyne.CanvasObject {
	opts, h := o.GetParams()
	w := widget.NewRadioGroup(opts, func(s string) {})
	w.Selected = GetValue(o)
	w.Horizontal = h
	return w
}
