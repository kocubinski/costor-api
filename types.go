package api

import "errors"

func (e *DecodeError) Error() string {
	return e.Err
}

func DecodeErr(err string, node *Node) error {
	return &DecodeError{
		Err:  err,
		Node: node,
	}
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
