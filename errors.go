package emailmanager

import "errors"

var (
	ErrBadRouting          = errors.New("inconsistent mapping between route and handler (programmer error)")
	ErrMailingIDNaN        = errors.New("mailing id is not a number")
	ErrInvalidArgumentType = errors.New("invalid argument type (programmer error)")
)
