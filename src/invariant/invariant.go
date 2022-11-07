package invariant

import (
	"errors"
	"fmt"
)

func Assert(condition bool, message string) {
	if !condition {
		fmt.Println("Invariant Assert failed.")
		fmt.Println(message)
		panic(errors.New(message))
	}
}

func AssertNoError(err error) {
	if err != nil {
		fmt.Println("Invariant AssertNoError failed.")
		fmt.Println(err.Error())
		panic(err)
	}
}
