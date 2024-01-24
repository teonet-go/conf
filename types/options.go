// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// RadioGroup options.

package types

import "slices"

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

// GetSelectedStr returns the selected option string of the radio group.
func (o RadioGroup) GetSelectedStr() (s string) {
	if o.Selected >= 0 && o.Selected < len(o.Options) {
		s = o.Options[o.Selected]
	}
	return
}

// GetValue returns the options, horizontal layout flag, and selected option
// string.
func (o RadioGroup) GetValue() (options []string, horizontal bool, selected string) {
	return o.GetOptions(), o.GetHorizontal(), o.GetSelectedStr()
}

// SetValue sets the selected option index of the radio group by the given
// string value.
func (o RadioGroup) SetValue(val string) RadioGroup {
	o.Selected = slices.Index(o.Options, val)
	return o
}
