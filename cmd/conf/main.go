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
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/conf"
)

// Person structure
type Person struct {
	Name string  `json:"name"`
	Age  float64 `json:"age"`
	Tst  int     `json:"tst"`
	Map  string  `json:"map"`
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

	// Create a form for editing the JSON data
	form := &widget.Form{}

	// Load the JSON data from a file
	filePath := "data.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Decode the JSON into a data structure
	var person Person
	err = json.Unmarshal(jsonData, &person)
	if err != nil {
		log.Fatal(err)
	}

	// Decode the JSON into map
	// var data map[string]any
	// err = json.Unmarshal(jsonData, &data)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var data = person

	// Get fields and its values from struct or map and add it to form
	fields := conf.GetFields(data, func(field *conf.Field[*widget.Entry]) {

		// Add field to form
		entry := widget.NewEntry()
		entry.SetText(field.ValueStr)

		// Add field validation by type
		entry.Validator = func(s string) (err error) {
			err = conf.ValidateValue(field, s)
			return
		}
		form.Append(field.NameDisplay, entry)

		// Add hint text to this forms entry
		form.Items[len(form.Items)-1].HintText = fmt.Sprintf("%s (%s)",
			field.NameDisplay, field.Type)

		// Set field entry to processing it in SetValues
		field.Entry = entry
	})

	// Create a save button
	saveButton := widget.NewButton("Save", func() {

		// Check if the form is valid
		if err := form.Validate(); err != nil {
			msg := fmt.Sprintf("Cannot save this form:\n %s", err)
			dialog.ShowError(fmt.Errorf(msg), w)
			return
		}

		// Update fields values
		fields.SetValues(&data, func(field *conf.Field[*widget.Entry]) string {
			return field.Entry.Text
		})

		// Encode the modified data structure back into JSON
		updatedJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Write the encoded JSON back to the file
		err = os.WriteFile(filePath, updatedJSON, 0644)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation("Success", "JSON file updated successfully!", w)
	})

	// Create a container for the form and save button
	content := container.NewVBox(
		form,
		saveButton,
	)

	// Set the window content
	w.SetContent(content)

	w.Resize(fyne.NewSize(500, 500))

	// Show the window
	w.ShowAndRun()
}
