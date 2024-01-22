package options

import "slices"

type RadioGroupType interface {
	GetOptions() []string
	GetHorizontal() bool
	GetSelected() int
	GetSelectedStr() string
	SetSelected(string)
}

type RadioGroup struct {
	Options    []string `json:"options"`
	Horizontal bool     `json:"horizontal"`
	Selected   int      `json:"selected"`
}

// GetOptions returns the options of the radio group.
func (o RadioGroup) GetOptions() []string { return o.Options }
func (o RadioGroup) GetHorizontal() bool  { return o.Horizontal }
func (o RadioGroup) GetSelected() int     { return o.Selected }
func (o RadioGroup) GetSelectedStr() (s string) {
	if o.Selected >= 0 && o.Selected < len(o.Options) {
		s = o.Options[o.Selected]
	}
	return
}
func (o *RadioGroup) SetSelected(val string) {
	o.Selected = slices.Index(o.Options, val)
}

// GetRadioGroup retrieves the radio group options from the given value.
//
// val: valuue by any type, the real type of the value may be RadioGroupType
//
// It returns:
//   - options: []string with radio group options
//   - horizontal: bool withtrue if the radio group is horizontal
//   - selected: string with the selected option
func GetRadioGroup(val any) (options []string, horizontal bool, selected string) {
	if val, ok := val.(RadioGroupType); ok {
		return val.GetOptions(), val.GetHorizontal(), val.GetSelectedStr()
	}
	return
}

func SetSelectedValue(opt any, val string) {

	if opt, ok := opt.(RadioGroup); ok {
		opt.SetSelected(val)
		return
	}

	if opt, ok := opt.(*RadioGroup); ok {
		opt.SetSelected(val)
	}
}
