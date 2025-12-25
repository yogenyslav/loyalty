// Package errs provides error handling utils.
package errs

import (
	"fmt"
	"runtime"
)

// Wrap fits the error in a chain, reports source file and provides optional description.
func Wrap(e error, desc ...string) error {
	if e == nil {
		return nil
	}
	var d string
	if len(desc) > 0 {
		d = desc[0] + " "
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("undefined call %s-> %w", d, e)
	}
	return fmt.Errorf("%s:%d %s-> %w", file, line, d, e)
}
