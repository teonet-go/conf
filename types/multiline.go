// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// MultiLine text entry.

package types

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/conf"
)

// Password type.
type Multiline struct {
	Value         string `json:"value"`
	MultiLineRows int    `json:"multiline_rows"`
}

// GetValue returns the value of the password.
func (m Multiline) GetValue() string {
	return m.Value
}

// GetNumRows returns the number of visible rows without scrolling of the
// multiline entry.
func (m Multiline) GetNumRows() int {
	return m.MultiLineRows
}

// SetValue sets the value of the password.
func (m Multiline) SetValue(val string) Multiline {
	m.Value = val
	return m
}

// SetNumRows sets the number of visible rows without scrolling of the widget.
func (m *Multiline) SetNumRows(num int) {
	m.MultiLineRows = num
}

// GetWidgetValue returns the widget value.
func (p Multiline) GetWidgetValue(field *conf.Field[fyne.CanvasObject]) string {
	return field.Entry.(*widget.Entry).Text
}

// NewWidget creates and returns widget and true if hint for this field is
// supported.
func (m Multiline) NewWidget() (fyne.CanvasObject, bool) {
	w := widget.NewMultiLineEntry()
	w.SetMinRowsVisible(m.GetNumRows())
	w.SetText(GetValue(m))
	return w, false
}
