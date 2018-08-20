// +build ignore

package main

import (
	"fmt"
	"github.com/Al2Klimov/go-generate-deps"
	"os"
)

func main() {
	if errGD := go_generate_deps.GenDeps(); errGD != nil {
		fmt.Println(errGD.Error())
		os.Exit(1)
	}
}
