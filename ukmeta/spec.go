package ukmeta

type Parameter[T any] func(index []int) (T, error)

func ParameterConst[T any](val T, err error) Parameter[T] {
	return func([]int) (T, error) { return val, err }
}
