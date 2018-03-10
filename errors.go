package tugboat

import (
	"fmt"
	"strings"
)

func errf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func wrap(err error, msg string, args ...interface{}) error {
	return errf("%s: %s", fmt.Sprintf(msg, args...), err)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type MultiError []error

func (me MultiError) Error() string {
	var strs []string
	for _, e := range me {
		strs = append(strs, e.Error())
	}
	return strings.Join(strs, "; ")
}

func (me *MultiError) Try(err error) {
	if err != nil {
		*me = append(*me, err)
	}
}

func (me *MultiError) Finish(err *error) {

	r := recover()
	if r != nil {
		if e, ok := r.(error); ok {
			*me = append(*me, e)
		} else {
			e := errf("Unknown panic: %+v", r)
			*me = append(*me, e)
		}
	}

	if len(*me) > 0 {
		*err = *me
	}
}
