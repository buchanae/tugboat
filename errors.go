package tugboat

type Must func(error)
type Finish func(error) error
type Try func(error) error

func Errors() (Try, Must, Finish) {

	var errors MultiError

	finish := func(err error) error {
		if perr := CapturePanic(); perr != nil {
			errors = append(errors, perr)
		}
		if len(errors) > 0 {
			return errors
		}
		return nil
	}

	try := func(err error) error {
		if err != nil {
			errors = append(errors, err)
			return err
		}
		return nil
	}

	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	return try, must, finish
}
