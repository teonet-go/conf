// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Email entry.

package types

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/conf"
)

// Email type.
type Email string

// GetValue returns the value of the pmail.
func (p Email) GetValue() string {
	return string(p)
}

// SetValue sets the value of the pmail.
func (p Email) SetValue(val string) Email {
	return Email(val)
}

// GetWidgetValue returns the widget value.
func (p Email) GetWidgetValue(field *conf.Field[fyne.CanvasObject]) string {
	return field.Entry.(*widget.Entry).Text
}

// NewWidget creates and returns widget and true if hint for this field is
// supported.
func (p Email) NewWidget() (fyne.CanvasObject, bool) {
	w := widget.NewEntry()
	w.SetPlaceHolder("test@example.com")
	w.SetText(GetValue(p))
	w.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")
	return w, true
}
