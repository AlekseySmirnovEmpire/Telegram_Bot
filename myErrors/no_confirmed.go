package myErrors

import "fmt"

type NotConfirmed struct {
	Val string
	Key string
}

func (n NotConfirmed) Error() string {
	return fmt.Sprintf("Вы ещё не подтвердили %s, чтобы начать пользоваться %s!", n.Val, n.Key)
}

func (n NotConfirmed) Unwrap() error {
	return n
}
