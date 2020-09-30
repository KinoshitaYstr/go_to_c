package main

import (
	"flag"
	"fmt"
	"os"
)

// トークンの種類用変数
type TokenKind int

// 種類の定数
const (
	TK_RESERVED TokenKind = iota // 予約時
	TK_NUM                       // 整数血
	TK_EOF                       // 入力終了用
)

// 現状のトークン
type Token struct {
	kind TokenKind // トークンの種類
	next *Token    // 次のトークン
	val  string    // 値
	str  string    // 文字
}

// 現状のトークン
var token *Token

// 入力プログラム
var user_input string

// 現状見ている場所
var now_loc int

// えらー
func error(str string) {
	tmp := ""
	for i := 0; i < now_loc; i++ {
		tmp += " "
	}
	fmt.Fprintln(os.Stderr, user_input)
	fmt.Fprintln(os.Stderr, tmp+"^ "+str)
	os.Exit(1)
}

// 次のトークンが予約されているものか確認
func consume(op string) bool {
	if token.kind != TK_RESERVED || string(token.str[0]) != op {
		return false
	}
	token = token.next
	return true
}

// 演算子確認
func expect(op string) {
	if token.kind != TK_RESERVED || string(token.str[0]) != op {
		error(op + "ではありません")
	}
	token = token.next
}

// 数値確認
func expect_number() string {
	if token.kind != TK_NUM {
		error("数値ではありません")
		return ""
	}
	val := token.val
	token = token.next
	return val
}

// EOF確認
func at_eof() bool {
	return token.kind == TK_EOF
}

// トークン作成とつなげる
func new_token(kind TokenKind, cur *Token, str string) *Token {
	var tok *Token
	tok = new(Token)
	tok.kind = kind
	tok.str = str
	tok.next = nil
	cur.next = tok
	return tok
}

// もじれつpをトーク内図
func tokenize(str string) *Token {
	var head Token
	head.next = nil
	cur := &head
	user_input = str
	now_loc = 0
	for len(str) > 0 {
		// 空白飛ばし
		if str[0] == ' ' {
			str = str[1:]
			now_loc += 1
			continue
		}

		// +/-演算子
		if str[0] == '+' || str[0] == '-' {
			cur = new_token(TK_RESERVED, cur, string(str[0]))
			str = str[1:]
			now_loc += 1
			continue
		}

		// 数値処理
		if '0' <= str[0] && str[0] <= '9' {
			cur = new_token(TK_NUM, cur, "")
			cur.val, str = get_number_string(str)
			continue
		}

		error("トークナイズできません")
	}
	new_token(TK_EOF, cur, str)
	return head.next
}

// 文字列からfloat64を取得して、読み取ったものを飛ばして返すもの
func get_number_string(data string) (string, string) {
	result := ""
	for {
		if len(data) != 0 && '0' <= data[0] && data[0] <= '9' {
			result += string(data[0])
			data = data[1:]
			now_loc += 1
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

	// 実際のとーくないずする
	token = tokenize(args[0])

	// // 実際の処理
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	// 最初の値を置く
	fmt.Println("  mov rax, " + expect_number())

	// その後に+/-があったら動かす
	for !at_eof() {
		// +のとき
		if consume("+") {
			fmt.Println("  add rax, " + expect_number())
			continue
		}

		// -のとき
		expect("-")
		fmt.Println("  sub rax, " + expect_number())
	}

	fmt.Println("  ret")
}
