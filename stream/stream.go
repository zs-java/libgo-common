package stream

type Function[T any, R any] func(t T) R

type Consumer[T any] func(t T)

type Predicate[T any] func(t T) bool

func Filter[T any](list []T, predicate Predicate[T]) []T {
	var result []T
	for _, t := range list {
		if predicate(t) {
			result = append(result, t)
		}
	}
	return result
}

func AnyMatch[T any](list []T, predicate Predicate[T]) bool {
	for _, t := range list {
		if predicate(t) {
			return true
		}
	}
	return false
}

func AllMatch[T any](list []T, predicate Predicate[T]) bool {
	for _, t := range list {
		if !predicate(t) {
			return false
		}
	}
	return true
}

func Map[T any, R any](list []T, fn Function[T, R]) []R {
	var result []R
	for _, t := range list {
		result = append(result, fn(t))
	}
	return result
}

func FlatMap[T any, R any](list []T, fn Function[T, []R]) []R {
	var result []R
	for _, t := range list {
		result = append(result, fn(t)...)
	}
	return result
}

func Reduce[T any, V any](list []T, getter func(t T) V, fn func(v1, v2 V) V) (v V) {
	for _, t := range list {
		v = fn(v, getter(t))
	}
	return v
}

func Summing[T int | int8 | int32 | int64 | float32 | float64 | string](list []T) (v T) {
	return SummingFn(list, func(v1, v2 T) T {
		return v1 + v2
	})
}

func SummingFn[T any](list []T, fn func(v1, v2 T) T) (v T) {
	identity := func(t T) T {
		return t
	}
	return Reduce(list, identity, fn)
}

func GroupBy[T any, K int | int8 | int32 | int64 | float32 | float64 | string](list []T, keyFn Function[T, K]) map[K][]T {
	return GroupByMapping(list, keyFn, func(t T) T {
		return t
	})
}

func GroupByMapping[T any, K int | int8 | int32 | int64 | float32 | float64 | string, V any](list []T, keyFn Function[T, K], valFn Function[T, V]) map[K][]V {
	return GroupByReduce(list, keyFn, func(t T) []V {
		return []V{valFn(t)}
	}, func(v1, v2 []V) []V {
		return append(v1, v2...)
	})
}

func GroupBySumming[T any, K int | int8 | int32 | int64 | float32 | float64 | string, V int | int8 | int32 | int64 | float32 | float64 | string](list []T, keyFn Function[T, K], valFn Function[T, V]) map[K]V {
	return GroupByReduce(list, keyFn, valFn, func(v1, v2 V) V {
		return v1 + v2
	})
}

func GroupByReduce[T any, K int | int8 | int32 | int64 | float32 | float64 | string, V any](list []T, keyFn Function[T, K], valFn Function[T, V], reduceFn func(v1, v2 V) V) map[K]V {
	result := make(map[K]V)
	for _, t := range list {
		key := keyFn(t)
		result[key] = reduceFn(result[key], valFn(t))
	}
	return result
}
