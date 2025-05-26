package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func Foo() error {
	return errors.New("foo error")
}

func Bar() error {
	err := Foo()
	if err != nil {
		return errors.WithMessage(err, "bar")
	}
	return nil
}

func main() {
	err := Bar()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}
