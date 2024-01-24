// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Password entry.

package types

// Password type.
type Password struct {
	Value string `json:"value"`
}

// GetValue returns the value of the password.
func (p Password) GetValue() string {
	return p.Value
}

// SetValue sets the value of the password.
func (p Password) SetValue(val string) Password {
	p.Value = val
	return p
}
