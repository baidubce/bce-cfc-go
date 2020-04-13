// Package cfc
package cfc

import (
	"fmt"
)

type MissingHandlerError struct {
	handlerName string
}

func(e MissingHandlerError)Error() string {
	return fmt.Sprintf("Missing handler %s", e.handlerName)
}
