package main

import (
	"flag"
	"fmt"
	"os"
)

// 文字列からfloat64を取得して、読み取ったものを飛ばして返すもの
func GetNumberString(data string) (string, string) {
	result := ""
	for {
		if len(data) != 0 && '0' <= data[0] && data[0] <= '9' {
			result += string(data[0])
			data = data[1:]
		} else {
			break
		}
	}
	return result, data
}

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
	// 値を取る
	arg := args[0]
	var val string
	for len(arg) > 0 {
		if arg[0] == '+' {
			arg = arg[1:]
			val, arg = GetNumberString(arg)
			fmt.Println("  add rax, " + val)
			continue
		}
		if arg[0] == '-' {
			arg = arg[1:]
			val, arg = GetNumberString(arg)
			fmt.Println("  sub rax, " + val)
			continue
		}
		val, arg = GetNumberString(arg)
		if len(val) != 0 {
			fmt.Println("  mov rax, " + val)
		}
	}
	fmt.Println("  ret")
	return
}
