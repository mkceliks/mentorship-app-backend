package errorPackage

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrNoSuchKey   = errors.New("NoSuchKey")
)
