package model

type Filter[T comparable] struct {
	Eq    *T
	Neq   *T
	In    []T
	NotIn []T
}

type IDFilter = Filter[int64]

type TextFilter = Filter[string]
