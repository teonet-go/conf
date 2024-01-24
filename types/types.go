package types

import (
	"github.com/teonet-go/conf"
)

type GetValueType interface {
	GetValue() string
}

type SetValueType[T any] interface {
	SetValue(val string) T
}

// GetValue returns the value of the GetValueType.
func GetValue[T GetValueType](val T) string { return val.GetValue() }

func SetValue[T SetValueType[T], F any](field *conf.Field[F], val string) {
	field.Value = field.Value.(T).SetValue(val)
}
