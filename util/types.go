package util

type Either[A any, B any] struct {
	isA bool
	a   A
	b   B
}

func Switch[A any, B any, R any](either Either[A, B],
	onA func(a A) R,
	onB func(b B) R,
) R {
	if either.isA {
		return onA(either.a)
	} else {
		return onB(either.b)
	}
}

func MakeA[A any, B any](a A) Either[A, B] {
	var result Either[A, B]
	result.isA = true
	result.a = a
	return result
}

func MakeB[A any, B any](b B) Either[A, B] {
	var result Either[A, B]
	result.isA = false
	result.b = b
	return result
}
