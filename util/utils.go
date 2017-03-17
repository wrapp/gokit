package util

type ErrFunc func() error

// A handy method to short circuit if any error happens when executing functions.
func ShortCurcuit(funcs ...ErrFunc) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
