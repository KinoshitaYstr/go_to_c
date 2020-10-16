package main

import (
	"fmt"
	"os"
)

func main() {
	// コマンドライン
	args := os.Args

	// 引数確認
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "引数の個数が正しくありません")
		return
	}

	tokens := tokenize(args[1])
	fmt.Println(tokens)

	nodes := tokens.program()
	fmt.Println(len(nodes))
	fmt.Println(nodes)

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println(".main")

	fmt.Println("  push rbp")
	fmt.Println("  mov rbp, rsp")
	fmt.Println("  sub rbp, " + string(locals.offset))
	nodes[0].gen()
}
