package tugboat

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func errf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func wrap(err error, msg string, args ...interface{}) error {
	return errf("%s: %s", fmt.Sprintf(msg, args...), err)
}

func CapturePanic() error {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			return wrap(e, "panic: %s", string(debug.Stack()))
		} else {
			return errf("Unknown panic: %+v", r)
		}
	}
	return nil
}

type MultiError []error

func (me MultiError) Error() string {
	var strs []string
	for _, e := range me {
		strs = append(strs, e.Error())
	}
	return strings.Join(strs, ": ")
}
