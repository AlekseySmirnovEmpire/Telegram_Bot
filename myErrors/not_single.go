package myErrors

import "fmt"

type NotSingle struct {
	Val string
	Key string
	Err error
}

func (n NotSingle) Error() string {
	return fmt.Sprintf("Find not only one in TABLE \"%s\", error: %s", n.Val, n.Key)
}

func (n NotSingle) Unwrap() error {
	return n.Err
}
