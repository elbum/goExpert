package main

import (
	"fmt"

	"reflect"

	"github.com/elbum/goExpert/something"
	"rsc.io/quote"
)

func main() {
	fmt.Println(quote.Hello())
	something.SayHello()

	const name = "hello"
	fmt.Println(reflect.TypeOf(name))

	bum := map[string]string{"name": "bum", "type": "man"}

	// bumslice = []map
	fmt.Println(bum)
}
