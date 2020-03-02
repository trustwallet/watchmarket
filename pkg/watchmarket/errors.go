package watchmarket

import "errors"

var (
	ErrNotFound = errNotFound()
)

func errNotFound() error { return errors.New("record does not exist") }
