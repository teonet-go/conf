// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// RadioGroup options.
package options

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

// SetSelected sets the selected option index of the radio group by the given
// string value.
func (o *RadioGroup) SetSelected(val string) {
	o.Selected = slices.Index(o.Options, val)
}

// RadioGroupType defines the interface of the radio group.
type RadioGroupType interface {
	GetOptions() []string
	GetHorizontal() bool
	GetSelected() int
	GetSelectedStr() string
	SetSelected(string)
}

// GetRadioGroup retrieves the radio group options from the given value.
//
// val: valuue by any type, the real type of the value may be RadioGroupType
//
// It returns:
//   - options: []string with radio group options
//   - horizontal: bool withtrue if the radio group is horizontal
//   - selected: string with the selected option
func GetRadioGroup(opt any) (options []string, horizontal bool, selected string) {
	if opt, ok := opt.(RadioGroup); ok {
		return opt.GetOptions(), opt.GetHorizontal(), opt.GetSelectedStr()
	}
	return
}

// SetValue sets the selected value for the given option.
//
//  opt - the option to set the value for, can be of type RadioGroup or *RadioGroup.
//  val - the value to be set for the selected option.
func SetValue(o any, val string) RadioGroup {
	opt, ok := o.(RadioGroup)
	if ok {
		opt.SetSelected(val)
	}
	return opt
}
