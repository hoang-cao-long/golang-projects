package main

import (
	"fmt"
	"reflect"

	"github.com/hoang-cao-long/golang-side-projects/learn-gorm/monorepo/services/gateway/utils"
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

	utils.DeprecatedFunc()
}
