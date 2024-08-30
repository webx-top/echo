package param

func IsNotEmptyString(s *string) bool {
	if s == nil {
		return false
	}
	return len(*s) > 0
}

func IsNotEmptyStringElement(_ int, s string) bool {
	return len(s) > 0
}

func IsTrueBoolElement(_ int, s bool) bool {
	return s
}

func IsFalseBoolElement(_ int, s bool) bool {
	return !s
}

type Number interface {
	~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint | ~int | ~uint64 | ~int64 | ~float32 | ~float64
}

type Scalar interface {
	Number | ~bool | ~string
}

func IsGreaterThanZeroElement[T Number](_ int, v T) bool {
	return v > 0
}

func FilterSlice[T comparable, F func(_ int, v T) bool](p []T, f F) []T {
	filtered := make([]T, 0, len(p))
	for index, value := range p {
		if f(index, value) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func Unique[T comparable](p []T) []T {
	record := map[T]struct{}{}
	result := make([]T, 0, len(p))
	for _, s := range p {
		if _, ok := record[s]; !ok {
			record[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}

func UniqueWithFilter[T comparable, F func(_ int, v T) bool](p []T, f F) []T {
	record := map[T]struct{}{}
	result := make([]T, 0, len(p))
	for i, s := range p {
		if !f(i, s) {
			continue
		}
		if _, ok := record[s]; !ok {
			record[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
