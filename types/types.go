package types

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/conf"
)

// Types interface type.
type Types[T any] interface {
	GetValue() string
	SetValue(val string) T
	NewWidget() (fyne.CanvasObject, bool)
	GetWidgetValue(field *conf.Field[fyne.CanvasObject]) string
}

// GetValue returns the value of the GetValueType.
func GetValue[T Types[T]](val T) string { return val.GetValue() }

// GetWidgetValue returns the widget value.
func GetWidgetValue[T Types[T]](val T, field *conf.Field[fyne.CanvasObject]) string {
	return val.GetWidgetValue(field)
}

// SetValue sets the value to field.Value.
func SetValue[T Types[T], F any](field *conf.Field[F], val string) {
	field.Value = field.Value.(T).SetValue(val)
}

// NewWidget creates and returns widget and true if the field type is supported.
func NewWidget[T Types[T], F any](field *conf.Field[F]) (fyne.CanvasObject, bool) {
	return field.Value.(T).NewWidget()
}

// CheckWidget creates and returns widget and true if the field type is supported.
func CheckWidget(field *conf.Field[fyne.CanvasObject]) (w fyne.CanvasObject, h, ok bool) {
	switch field.Type {
	case "types.Email":
		w, h = NewWidget[Email](field)
	case "types.RadioGroup":
		w, h = NewWidget[RadioGroup](field)
	case "types.Password":
		w, h = NewWidget[Password](field)
	case "types.Multiline":
		w, h = NewWidget[Multiline](field)
	default:
		return
	}
	ok = true
	return
}

// CheckSave checks if the field type is supported get widget value and set it
// to field using SetValue.
func CheckSave(field *conf.Field[fyne.CanvasObject]) (ok bool) {
	switch field.Type {
	case "types.Email":
		val := field.Entry.(*widget.Entry).Text
		SetValue[Email](field, val)
	case "types.Password":
		val := field.Entry.(*widget.Entry).Text
		SetValue[Password](field, val)
	case "types.Multiline":
		val := field.Entry.(*widget.Entry).Text
		SetValue[Multiline](field, val)
	case "types.RadioGroup":
		val := field.Entry.(*widget.RadioGroup).Selected
		SetValue[RadioGroup](field, val)
	default:
		return
	}
	ok = true
	return
}
