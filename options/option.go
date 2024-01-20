package options

type RadioGroupType interface {
	GetRadioGroup() []string
	GetHorizontal() bool
}

type RadioGroup struct {
	Options    []string `json:"options"`
	Horizontal bool     `json:"horizontal"`
	Selected   int      `json:"selected"`
}

// GetRadioGroup returns the options of the radio group.
func (o RadioGroup) GetRadioGroup() []string { return o.Options }
func (o RadioGroup) GetHorizontal() bool     { return o.Horizontal }

// GetRadioGroup retrieves the radio group options from the given value.
//
// val: valuue by any type, the real type of the value may be RadioGroupType
//
// It returns []string with radio group options.
func GetRadioGroup(val any) (options []string, horizontal bool) {
	if val, ok := val.(RadioGroupType); ok {
		options, horizontal = val.GetRadioGroup(), val.GetHorizontal()
	}
	return
}

// func getRadioGroup(val RadioGroupType) (options []string, horizontal bool) {
// 	return val.GetRadioGroup(), val.GetHorizontal()
// }
