package contnet

import "errors"

var Errors = struct {
	NotImplemented error
}{
	NotImplemented: errors.New("Not implemented."),
}
