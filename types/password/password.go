// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Password entry.
package password

// Password type.
type Password struct {
	Value string `json:"value"`
}

// SetValue sets the value of the password.
func (p *Password) SetValue(val string) {
	p.Value = val
}

// GetValue returns the value of the password.
func (p Password) GetValue() string {
	return p.Value
}

// GetValue returns the value of the password.
func GetValue(pas any) (password string) {
	if pas, ok := pas.(Password); ok {
		return pas.GetValue()
	}
	return
}

// SetValue sets the value of the password.
func SetValue(pas any, val string) Password {
	return Password{Value: val}
}
