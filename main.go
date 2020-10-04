package main

import (
	"fmt"
	"os"
)

func main() {
	// コマンドライン引数関係
	args := os.Args

	// 引数の確認(なかったらエラー)
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "引数の個数が正しくありません")
		return
	}

	// 実際のとーくないずする
	token = tokenize(args[1])

	InitLocals()
	nodes := Program()

	// 実際の処理
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	// 変数用のスペーーす確保
	fmt.Println("  push rbp")
	fmt.Println("  mov rbp, rsp")
	fmt.Println("  sub rsp,", GetLocalSpace())

	// ノードの生成とコード書く
	for _, node := range nodes {
		if node == nil {
			fmt.Println("nil !!")
			continue
		}
		Gen(node)
		fmt.Println("  pop rax")
	}

	fmt.Println("  mov rsp, rbp")
	fmt.Println("  pop rbp")
	fmt.Println("  ret")
}
