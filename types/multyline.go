// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// MultiLine text entry.

package types

// Password type.
type Multiline struct {
	Value         string `json:"value"`
	MultiLineRows int    `json:"multiLineRows"`
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
