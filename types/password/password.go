// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Password entry.
package password

// Password type.
type Password struct {
	Entry   string `json:"entry"`
	Visible bool   `json:"visible"`
}

// SetValue sets the value of the password.
func (p *Password) SetValue(val string) {
	p.Entry = val
}

// GetValue returns the value of the password.
func (p Password) GetValue() string {
	return p.Entry
}

// PasswordType defines the interface of the password.
type PasswordType interface {
	SetValue(val string)
	GetValue() string
}

// GetValue returns the value of the password.
func GetValue(pas any) (password string) {
	if pas, ok := pas.(PasswordType); ok {
		return pas.GetValue()
	}
	return
}

// SetValue sets the value of the password.
func SetValue(pas any, val string) {
	if pas, ok := pas.(PasswordType); ok {
		pas.SetValue(val)
	}
}

