package myErrors

import "fmt"

type EmptyFile struct {
	Val     string
	isExist bool
}

func (e EmptyFile) Error() string {
	if e.isExist {
		return fmt.Sprintf("File \"%s\" is empty!", e.Val)
	} else {
		return fmt.Sprintf("There is no \"%s\" file!", e.Val)
	}
}

func (e EmptyFile) Unwrap() error {
	return e
}
