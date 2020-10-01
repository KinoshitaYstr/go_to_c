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
	// ノードの生成
	var node *Node
	node = Expr()

	// // 実際の処理
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	// ノードから汗かく
	Gen(node)

	fmt.Println("  pop rax")
	fmt.Println("  ret")
}
