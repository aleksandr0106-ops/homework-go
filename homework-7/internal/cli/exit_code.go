package cli

import (
	"errors"
	"flag"
)

func ExitCode(err error) int {
	switch {
	case err == nil, errors.Is(err, flag.ErrHelp):
		return 0
	case errors.Is(err, ErrUsage):
		return 2
	default:
		return 1
	}
}
