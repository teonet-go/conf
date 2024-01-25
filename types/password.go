// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Password entry.

package types

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/conf"
)

// Password type.
type Password string

// GetValue returns the value of the password.
func (p Password) GetValue() string {
	return string(p)
}

// SetValue sets the value of the password.
func (p Password) SetValue(val string) Password {
	return Password(val)
}

// GetWidgetValue returns the widget value.
func (p Password) GetWidgetValue(field *conf.Field[fyne.CanvasObject]) string {
	return field.Entry.(*widget.Entry).Text
}

// NewWidget creates and returns widget and true if hint for this field is
// supported.
func (p Password) NewWidget() (fyne.CanvasObject, bool) {
	w := widget.NewPasswordEntry()
	w.SetText(GetValue(p))
	return w, false
}
