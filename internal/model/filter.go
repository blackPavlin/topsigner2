package model

import "uuid"

type Filter[T comparable] struct {
	Eq    *T
	Neq   *T
	In    []T
	NotIn []T
}

type IDFilter = Filter[int64]

type UUIDFileter = Filter[uuid.UUID]

type TextFilter = Filter[string]
