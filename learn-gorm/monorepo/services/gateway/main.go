// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
)

type test struct {
	Name string
	Age  int
}

func main() {
	// a := test{
	// 	Name: "long",
	// 	Age:  18,
	// }
	b := reflect.TypeOf(test{})

	fmt.Println(b.FieldByIndex([]int{0, 1}))
}
