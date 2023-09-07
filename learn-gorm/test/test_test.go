package test

import (
	"fmt"
	"testing"
)

func Parent() string {
	return "Parent"
}

func Child() string {
	return "Child"
}

func TestMain(t *testing.T) {
	fmt.Println(Parent())

	t.Run("test child", func(t *testing.T) {
		fmt.Println(Child())
	})
}
