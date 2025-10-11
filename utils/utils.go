package utils

var Mouse_X float64
var Mouse_Y float64

func RemoveArrayElement[T any](index_to_remove int, slice *[]T) {
	if len(*slice) == 0 {
		return
	}
	*slice = append((*slice)[:index_to_remove], (*slice)[index_to_remove+1:]...)
}

type Vec2 struct {
	X, Y float64
}
