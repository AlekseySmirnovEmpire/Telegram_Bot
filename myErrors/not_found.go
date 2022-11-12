package myErrors

import "fmt"

type NotFound struct {
	Val string
	Key string
	Err error
}

func (n NotFound) Error() string {
	return fmt.Sprintf("Cannot found: \"%s\" in \"%s\"", n.Val, n.Key)
}

func (n NotFound) Unwrap() error {
	return n.Err
}
