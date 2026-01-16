package main

const (
	ErrCounterNotFound = "COUNTER_NOT_FOUND"
)

type Counter struct {
	Value       int   `json:"value"`
	LastUpdated int64 `json:"last_updated"`
}
