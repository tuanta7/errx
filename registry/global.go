package registry

import (
	"github.com/tuanta7/errx/predefined"
)

var (
	Global = &Registry{
		StatusCodeMap: predefined.DefaultErrorStatusCodeMap,
		MessageMap:    make(map[string]map[string]string),
	}
)

// SetGlobal is ony safe to call during initialization phase.
func SetGlobal(e *Registry) {
	Global = e
}
