// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// RadioGroup options.

package types

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/conf"
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

// SetOptions sets the options of the radio group.
func (o *RadioGroup) SetOptions(options []string) {
	o.Options = options
}

// SetHorizontal sets the horizontal type of the radio group type.
func (o *RadioGroup) SetHorizontal() {
	o.Horizontal = true
}

// SetVertical sets the vertical type of the radio group type.
func (o *RadioGroup) SetVertical() {
	o.Horizontal = false
}

// GetWidgetValue returns the widget value.
func (p RadioGroup) GetWidgetValue(field *conf.Field[fyne.CanvasObject]) string {
	return field.Entry.(*widget.RadioGroup).Selected
}

// NewWidget creates and returns widget and true if hint for this field is
// supported.
func (o RadioGroup) NewWidget() (fyne.CanvasObject, bool) {
	opts, h := o.GetParams()
	w := widget.NewRadioGroup(opts, func(s string) {})
	w.Selected = GetValue(o)
	w.Horizontal = h
	return w, false
}
