package types

import (
	"fyne.io/fyne/v2"
	"github.com/teonet-go/conf"
)

type GetValueType interface {
	GetValue() string
}

type SetValueType[T any] interface {
	SetValue(val string) T
}

type NewWidgetType interface {
	NewWidget() fyne.CanvasObject
}

// GetValue returns the value of the GetValueType.
func GetValue[T GetValueType](val T) string { return val.GetValue() }

func SetValue[T SetValueType[T], F any](field *conf.Field[F], val string) {
	field.Value = field.Value.(T).SetValue(val)
}

func NewWidget[T NewWidgetType, F any](field *conf.Field[F]) fyne.CanvasObject {
	return field.Value.(T).NewWidget()
}

// CheckWidget creates and returns widget and true if the field type is supported.
func CheckWidget(field *conf.Field[fyne.CanvasObject]) (w fyne.CanvasObject, ok bool) {
	switch field.Type {
	case "types.RadioGroup":
		w = NewWidget[RadioGroup](field)
	case "types.Password":
		w = NewWidget[Password](field)
	case "types.Multiline":
		w = NewWidget[Multiline](field)
	default:
		return
	}
	ok = true
	return
}
