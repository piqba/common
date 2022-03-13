package tools

import (
	"errors"
	"time"
)

type (
	Any interface{}
)

var (
	ErrElementNotFound = errors.New("element not present in Set")
)

// Set - our representation of set data structure
type Set struct {
	Elements map[Any]struct{}
}

// TimeChecker ...
type TimeChecker struct {
	Start  time.Time
	End    time.Time
	Check  time.Time
	Tz     string
	Layout string
}
