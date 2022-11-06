package myErrors

import "fmt"

type NotSingle struct {
	Val string
	Err string
}

func (n NotSingle) Error() string {
	return fmt.Sprintf("Find not only one in TABLE \"%s\", error: %s", n.Val, n.Err)
}

func (n NotSingle) Unwrap() error {
	return n
}
