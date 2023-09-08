package api

import "errors"

func (e *DecodeError) Error() string {
	return e.Reason
}

func IsDecodeError(err error) bool {
	if err == nil {
		return false
	}
	var decodeError *DecodeError
	ok := errors.As(err, &decodeError)
	return ok
}

type Changeset struct {
	Version int64
	Nodes   []*Node
}
