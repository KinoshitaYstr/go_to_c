package main

import (
	"fmt"
)

func main() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")
	fmt.Println("  mov rax, 42")
	fmt.Println("  ret")
}
