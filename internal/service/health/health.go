package health

import "time"

const checkTimeout = 2 * time.Second

type Status string

const (
	StatusOK       Status = "ok"
	StatusDegraded Status = "degraded"
	StatusDown     Status = "down"
)

type Report struct {
	Status Status
}
