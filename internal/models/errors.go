package models

import (
	"errors"
)

// ErrNoRecord is returned when a requested snippet or record does not exist.
var ErrNoRecord = errors.New("models: no matching record found")
