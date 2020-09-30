package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// コマンドライン引数関係
	flag.Parse()
	args := flag.Args()
	// 引数の確認(なかったらエラー)
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "引数の個数が正しくありません")
		return
	}
	// 実際の処理
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")
	fmt.Println("  mov rax, " + args[0])
	fmt.Println("  ret")
	return
}
