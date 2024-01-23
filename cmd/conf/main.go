// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Config helper go package GUI sample application using Fyne widgets to show
// json config data.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/teonet-go/conf/fyne/form"
	"github.com/teonet-go/conf/types/options"
	"github.com/teonet-go/conf/types/password"
)

// Person structure
type Person struct {
	Name     string             `json:"name"`
	Age      float64            `json:"age"`
	Tst      int                `json:"tst"`
	Map      string             `json:"map"`
	Password password.Password  `json:"password"`
	On       bool               `json:"on"`
	IntArray []int              `json:"int_array"`
	FltArray []float64          `json:"float_array"`
	Option   options.RadioGroup `json:"option"`
}

// main is the entry point of the program.
//
// It creates a new Fyne application, a new window, and a form for editing JSON
// data. It loads the JSON data from a file, decodes it into a data structure,
// and adds the fields and values to the form. It creates a save button that
// validates the form, updates the field values, encodes the modified data
// structure back into JSON, and writes it to a file. Finally, it sets the
// window content, resizes it, and shows the window.
func main() {
	// Create a new Fyne application
	a := app.New()

	// Create a new window
	w := a.NewWindow("JSON Editor")

	// Initialize the data structure with default values. Special form field
	// types must be initialized.
	var person Person

	// Load the JSON data from a file
	err := loadJson(&person)
	if err != nil {
		log.Fatal(err)
	}

	// Decode the JSON into map
	// var data map[string]any
	// err = loadJson(&data)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var data = person

	// Create a form from the struct or map that contains JSON data
	form := form.New(data)

	// Create a save button
	saveButton := form.NewSaveButton(&data,
		// Save button callback
		func() {
			// Write the encoded JSON back to the file and show Info dialog or
			// show error dialog at error.
			if err := saveJson(data); err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation(
				"Success", "JSON file updated successfully!",
				w,
			)
		},
		// Form validation error callback
		func(err error) {
			msg := fmt.Sprintf("Cannot save this form:\n %s", err)
			dialog.ShowError(fmt.Errorf(msg), w)
		},
	)

	// Create a container for the form and save button
	content := container.NewVBox(form, saveButton)

	// Set the window content
	w.SetContent(content)

	// Resize the window
	w.Resize(fyne.NewSize(500, 500))

	// Show the window
	w.ShowAndRun()
}

const filePath = "data.json"

func loadJson(v any) error {
	// Load the JSON data from a file
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Decode the JSON into a data structure
	return json.Unmarshal(jsonData, v)
}

func saveJson(v any) error {

	// Encode the modified data structure back into JSON
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	// Write the encoded JSON back to the file
	return os.WriteFile(filePath, data, 0644)
}
