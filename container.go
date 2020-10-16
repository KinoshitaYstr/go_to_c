package main

import (
	"fmt"
	"os"
)

func error(str string) {
	fmt.Fprintln(os.Stderr, str)
	os.Exit(1)
}
