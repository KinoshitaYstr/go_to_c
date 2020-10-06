// ぱーーさ

package main

import "fmt"

type Label struct {
	if_label    bool
	else_label  bool
	while_label bool
	for_label   bool
	if_num      int
	else_num    int
	while_num   int
	for_num     int
}

func (l *Label) getIf() string {
	if !l.if_label {
		l.if_label = true
		return fmt.Sprintf(".Lif%02d", l.if_num)
	}
	return ""
}

func (l *Label) getElse() string {
	if !l.else_label {
		l.else_label = true
		return fmt.Sprintf(".Lelse%02d", l.else_num)
	}
	return ""
}

func (l *Label) getWhile() string {
	if !l.while_label {
		l.while_label = true
		return fmt.Sprintf(".Lwhile%02d", l.while_num)
	}
	return ""
}

func (l *Label) getFor() string {
	if !l.while_label {
		l.for_label = true
		return fmt.Sprintf(".Lfor%02d", l.for_num)
	}
	return ""
}

func (l *Label) nextIf() {
	if l.if_label {
		l.if_num++
		l.if_label = false
	}
}

func (l *Label) nextElse() {
	if l.else_label {
		l.else_num++
		l.else_label = false
	}
}

func (l *Label) nextWhile() {
	if l.while_label {
		l.while_num++
		l.while_label = false
	}
}

func (l *Label) nextFor() {
	if l.for_label {
		l.for_num++
		l.for_label = false
	}
}

// ラベル管理
var labels Label

type LVar struct {
	next   *LVar  // 次の変数
	name   string // 名前
	offset int    // PBRからのオフセット
}

// ローカル変数
var locals *LVar

func InitLocals() {
	locals = new(LVar)
}

func GetLocalSpace() int {
	return locals.offset
}

func find_lbar(tok *Token) *LVar {
	for v := locals; v != nil; v = v.next {
		if v.name == tok.str {
			return v
		}
	}
	return nil
}

// ノード種類変数
type NodeKind int

// 抽象木のノード種類
const (
	ND_ADD    NodeKind = iota // +
	ND_SUB                    // -
	ND_MUL                    // *
	ND_DIV                    // /
	ND_NUM                    // 整数
	ND_EQU                    // ==
	ND_NEQ                    // !=
	ND_BIG                    // >
	ND_SML                    // <
	ND_EBG                    // >=
	ND_ESM                    // <=
	ND_ASSIGN                 // =
	ND_LVAR                   // ローカル変数
	ND_RETURN                 // return
	ND_IF                     // if
	ND_ELSE                   // else
	ND_WHILE                  // while
	ND_FOR                    // for
	ND_BLOCK                  // block
	ND_FUNC                   // 関数
	ND_ARG                    // 引数
	ND_NONE                   // none
)

// 抽象木の型
type Node struct {
	kind   NodeKind // ノードの方
	lhs    *Node    // 左辺
	rhs    *Node    // 右辺
	val    string   // せいすうち
	offset int      // ローカル変数のオフセット
	label  string   // ラベル保持
}

// 新しいノードの生成
func new_node(kind NodeKind, lhs *Node, rhs *Node) *Node {
	var node *Node
	node = new(Node)
	node.kind = kind
	node.lhs = lhs
	node.rhs = rhs
	return node
}

// せいすうよう
func new_node_num(val string) *Node {
	var node *Node
	node = new(Node)
	node.kind = ND_NUM
	node.val = val
	return node
}

func new_node_none() *Node {
	var node *Node
	node = new(Node)
	node.kind = ND_NONE
	return node
}

// ローカル変数よう
func new_node_local_variance() *Node {
	var node *Node
	node = new(Node)
	node.kind = ND_LVAR
	return node
}

func Program() []*Node {
	var node []*Node
	for !at_eof() {
		node = append(node, stmt())
	}
	return node
}

func stmt() *Node {
	var node *Node
	if consume("{") {
		if consume("}") {
			return nil
		}
		node = stmt()
		for !consume("}") {
			node = new_node(ND_BLOCK, node, stmt())
			// expect("}")
		}
		return node
	} else if consume("return") {
		node = new(Node)
		node.kind = ND_RETURN
		node.lhs = expr()
	} else if consume("if") {
		expect("(")
		node = new(Node)
		labels.nextIf()
		node.label = labels.getIf()
		node.kind = ND_IF
		node.lhs = expr()
		expect(")")
		node.rhs = stmt()
		if consume("else") {
			node = new_node(ND_ELSE, node, stmt())
			labels.nextElse()
			node.label = labels.getElse()
		}
		return node
	} else if consume("while") {
		expect("(")
		node = new(Node)
		node.kind = ND_WHILE
		labels.nextWhile()
		node.label = labels.getWhile()
		node.lhs = expr()
		expect(")")
		node.rhs = stmt()
		return node
	} else if consume("for") {
		expect("(")
		node = new(Node)
		node.kind = ND_FOR
		labels.nextFor()
		node.label = labels.getFor()
		node.lhs = new(Node)
		node.rhs = new(Node)
		if !consume(";") {
			node.lhs.lhs = expr()
			expect(";")
		} else {
			node.lhs.lhs = new_node_none()
		}
		if !consume(";") {
			node.lhs.rhs = expr()
			expect(";")
		} else {
			node.lhs.rhs = new_node_none()
		}
		if !consume(")") {
			node.rhs.lhs = expr()
			expect(")")
		} else {
			node.rhs.lhs = new_node_none()
		}
		node.rhs.rhs = stmt()
		return node
	} else {
		node = expr()
	}
	expect(";")
	return node
}

// 式のノード生成
func expr() *Node {
	return assign()
}

func assign() *Node {
	var node *Node
	node = equality()
	if consume("=") {
		node = new_node(ND_ASSIGN, node, assign())
	}
	return node
}

// 等式のノード生成
func equality() *Node {
	var node *Node
	node = relational()
	for {
		if consume("==") {
			node = new_node(ND_EQU, node, relational())
		} else if consume("!=") {
			node = new_node(ND_NEQ, node, relational())
		} else {
			return node
		}
	}
}

// 大小関係のノード生成
func relational() *Node {
	var node *Node
	node = add()
	for {
		if consume("<") {
			node = new_node(ND_SML, node, add())
		} else if consume("<=") {
			node = new_node(ND_ESM, node, add())
		} else if consume(">") {
			node = new_node(ND_BIG, node, add())
		} else if consume(">=") {
			node = new_node(ND_EBG, node, add())
		} else {
			return node
		}
	}
}

// 加算減算のノード生成
func add() *Node {
	var node *Node
	node = mul()
	for {
		if consume("+") {
			node = new_node(ND_ADD, node, mul())
		} else if consume("-") {
			node = new_node(ND_SUB, node, mul())
		} else {
			return node
		}
	}
}

// 乗算除算のノード生成
func mul() *Node {
	var node *Node
	node = unary()
	for {
		if consume("*") {
			node = new_node(ND_MUL, node, unary())
		} else if consume("/") {
			node = new_node(ND_DIV, node, unary())
		} else {
			return node
		}
	}
}

// 単行+/-のノード生成
func unary() *Node {
	if consume("+") {
		return primary()
	}
	if consume("-") {
		return new_node(ND_SUB, new_node_num("0"), primary())
	}
	return primary()
}

// 値等のプライマーのノード生成
func primary() *Node {
	// ()について
	if consume("(") {
		var node *Node
		node = expr()
		expect(")")
		return node
	}

	tok := consume_ident()
	if tok != nil {
		var node *Node
		node = new(Node)
		if consume("(") {
			// 関数
			node.kind = ND_FUNC
			if !consume(")") {
				argNode := expr()
				argNode.kind = ND_ARG
				if consume(",") {
					argNode = new_node(ND_ARG, argNode, expr())
				}
				expect(")")
				node.lhs = argNode
			} else {
				node.lhs = new_node_none()
			}
			node.rhs = stmt()
			return node
		}
		node.kind = ND_LVAR

		lvar := find_lbar(tok)
		if lvar != nil {
			node.offset = lvar.offset
		} else {
			var lvar *LVar
			lvar = new(LVar)
			lvar.next = locals
			lvar.name = tok.str
			lvar.offset = locals.offset + 8
			node.offset = lvar.offset
			locals = lvar
		}
		return node
	}

	if consume_none() {
		var node *Node
		node = new(Node)
		node.kind = ND_NONE
		return node
	}

	// 数値出す
	return new_node_num(expect_number())
}

// トークンの種類用変数
type TokenKind int

// 種類の定数
const (
	TK_RESERVED TokenKind = iota // 予約時
	TK_IDENT                     // 識別子
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

// 次のトークンが予約されているものか確認
func consume(op string) bool {
	if token.kind != TK_RESERVED || token.str != op {
		return false
	}
	token = token.next
	return true
}

// 変数？
func consume_ident() *Token {
	if token.kind == TK_IDENT {
		var tok *Token
		tok = token
		token = token.next
		return tok
	}
	return nil
}

func consume_none() bool {
	if token != nil {
		return false
	}
	token = token.next
	return true
}

// 演算子確認
func expect(op string) {
	if token.kind != TK_RESERVED || string(token.str[0]) != op {
		Error(op + "ではありません")
	}
	token = token.next
}

// 数値確認
func expect_number() string {
	if token.kind != TK_NUM {
		Error("数値ではありません")
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

		// return文
		if check_key_word("return", str) {
			cur = new_token(TK_RESERVED, cur, "return")
			str = str[6:]
			now_loc += 6
			continue
		}

		// if文
		if check_key_word("if", str) {
			cur = new_token(TK_RESERVED, cur, "if")
			str = str[2:]
			now_loc += 2
			continue
		}

		// else文
		if check_key_word("else", str) {
			cur = new_token(TK_RESERVED, cur, "else")
			str = str[4:]
			now_loc += 4
			continue
		}

		// while文
		if check_key_word("while", str) {
			cur = new_token(TK_RESERVED, cur, "while")
			str = str[5:]
			now_loc += 5
			continue
		}

		// for文
		if check_key_word("for", str) {
			cur = new_token(TK_RESERVED, cur, "for")
			str = str[3:]
			now_loc += 3
			continue
		}

		// == / <= / >= 演算子
		if len(str) >= 2 && (str[:2] == "==" || str[:2] == "!=" || str[:2] == "<=" || str[:2] == ">=") {
			cur = new_token(TK_RESERVED, cur, str[:2])
			str = str[2:]
			now_loc += 2
			continue
		}

		// + / - / */ / / () / < / > 演算子
		if str[0] == '+' || str[0] == '-' || str[0] == '*' || str[0] == '/' || str[0] == '(' || str[0] == ')' || str[0] == '>' || str[0] == '<' || str[0] == ';' || str[0] == '=' || str[0] == '{' || str[0] == '}' || str[0] == ',' {
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

		// ローカル変数
		if check_alphabet(str[0]) {
			name := ""
			for len(str) > 0 {
				if check_alphabet(str[0]) {
					name += string(str[0])
					str = str[1:]
				} else {
					break
				}
			}
			cur = new_token(TK_IDENT, cur, name)
			now_loc += len(name)
			continue
		}

		Error("トークナイズできません")
	}
	new_token(TK_EOF, cur, str)
	return head.next
}

// アルファベット化確認
func check_alphabet(c uint8) bool {
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
		return true
	}
	return false
}

// 数か確認
func check_num(c uint8) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	return false
}

// アンダースコア確認
func check_under_score(c uint8) bool {
	if '_' == c {
		return true
	}
	return false
}

func check_key_word(key string, str string) bool {
	if len(str) == len(key) {
		return str == key
	} else if len(str) < len(key) {
		return false
	} else if check_alphabet(str[len(key)]) || check_num(str[len(key)]) || check_under_score(str[len(key)]) {
		return false
	} else {
		return str[0:len(key)] == key
	}
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
