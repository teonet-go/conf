package main

import (
	"encoding/json"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Person struct {
	Name string  `json:"name"`
	Age  float64 `json:"age"`
	Tst  int     `json:"tst"`
	Map  string  `json:"map"`
}

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

	// Get fields and its values from struct and add it to form
	fields := GetFields(person, func(field *Field[*widget.Entry]) {
		entry := widget.NewEntry()
		entry.SetText(field.ValueStr)
		form.Append(field.Name, entry)
		field.Entry = entry
	})

	// Create a save button
	saveButton := widget.NewButton("Save", func() {

		// Get the values from the form.
		// Range by fields slice and set value of Entry field to the person
		// struct
		// for _, field := range fields {
		// 	SetValue(&person, field.Name, field.Entry.Text)
		// }
		fields.SetValues(&person, func(field *Field[*widget.Entry]) string {
			return field.Entry.Text
		})

		// Encode the modified data structure back into JSON
		updatedJSON, err := json.MarshalIndent(person, "", "  ")
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
