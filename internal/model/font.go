package model

import "time"

type Font struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FontFilter struct {
	Name TextFilter
}

type FontQuery struct {
	Filter     FontFilter
	Pagination Pagination
}
