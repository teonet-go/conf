// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Packages form create fine-go forms widget.
package form

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/conf"
	"github.com/teonet-go/conf/types"
)

// Form is a widget that creates fine-go form widget.
type Form struct {
	*widget.Form
	fields conf.Fields[fyne.CanvasObject]
}

func New(o any) *Form {
	f := &Form{Form: widget.NewForm()}
	f.GetFields(o)
	return f
}

// Append adds a new field to the form.
func (f *Form) Append(field *conf.Field[fyne.CanvasObject]) {

	var w fyne.CanvasObject // Widget
	var d string            // Name to displat
	var h bool              // Hint text is available if true

	switch field.Type {

	// Bool fields
	case "bool":

		// Add checkbox to form
		check := widget.NewCheck(field.NameDisplay, func(bool) {})
		check.Checked = field.Value.(bool)

		w = check

	// Any other simple fields displayed as string: string, int, float, etc.
	default:

		// Check special types and create its widget
		if widget, ok := types.CheckWidget(field); ok {
			w = widget
			d = field.NameDisplay
			break
		}

		// Add text entry field to form
		entry := widget.NewEntry()
		entry.SetText(field.ValueStr)

		// Add field validation by type
		entry.Validator = func(s string) (err error) {
			err = field.ValidateValue(s)
			return
		}

		h = true
		w = entry
		d = field.NameDisplay
	}

	// Append field to form
	f.Form.Append(d, w)

	// Set hint text to this forms entry
	if h {
		f.Form.Items[len(f.Form.Items)-1].HintText = fmt.Sprintf("%s (%s)",
			field.NameDisplay, field.Type)
	}

	// Set field entry to processing it in SetValues
	field.Entry = w
}

func (f *Form) GetFields(o any) {
	f.fields = conf.GetFields(o, func(field *conf.Field[fyne.CanvasObject]) {
		f.Append(field)
	})
}

func (f *Form) NewSaveButton(o any, save func(), valerr func(err error)) *widget.Button {

	return widget.NewButton("Save", func() {

		// Check if the form is valid
		if err := f.Validate(); err != nil {
			valerr(err)
			return
		}

		// Update fields values
		f.fields.SetValues(o, func(field *conf.Field[fyne.CanvasObject]) (string, bool) {
			switch field.Type {

			// Bool fields
			case "bool":
				val := field.Entry.(*widget.Check).Checked
				return fmt.Sprintf("%v", val), true

			// The options.RadioGroup fields
			case "types.RadioGroup":
				val := field.Entry.(*widget.RadioGroup).Selected
				types.SetValue[types.RadioGroup](field, val)

			// The password fields
			case "types.Password":
				val := field.Entry.(*widget.Entry).Text
				types.SetValue[types.Password](field, val)

			// The multi-line text fields
			case "types.Multiline":
				val := field.Entry.(*widget.Entry).Text
				types.SetValue[types.Multiline](field, val)

			// Any other simple fields displayed as string: string, int, float,
			// etc.
			default:
				return field.Entry.(*widget.Entry).Text, true
			}

			return "", false
		})

		// Use save callback to encode json and Write back to the file
		save()
	})
}
