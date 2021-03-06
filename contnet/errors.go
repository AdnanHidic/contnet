package contnet

import "errors"

var Errors = struct {
	NotImplemented  error
	ContentNotFound error
	ProfileNotFound error
}{
	NotImplemented:  errors.New("Not implemented."),
	ContentNotFound: errors.New("Content not found in content store."),
	ProfileNotFound: errors.New("Profile not found in profile store."),
}
