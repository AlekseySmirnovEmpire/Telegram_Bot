package myErrors

import "fmt"

type NoConnection struct {
	Val string
	Key string
	Err error
}

func (n NoConnection) Error() string {
	return fmt.Sprintf("Error connection to \"%s\", error is: %s", n.Val, n.Key)
}

func (n NoConnection) Unwrap() error {
	return n.Err
}
