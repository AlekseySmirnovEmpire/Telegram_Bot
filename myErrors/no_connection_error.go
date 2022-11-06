package myErrors

import "fmt"

type NoConnection struct {
	Val string
	Err string
}

func (n NoConnection) Error() string {
	return fmt.Sprintf("Error connection to \"%s\", error is: %s", n.Val, n.Err)
}

func (n NoConnection) Unwrap() error {
	return n
}
