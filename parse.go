// ぱーーさ

package main

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
)

// 抽象木の型
type Node struct {
	kind   NodeKind // ノードの方
	lhs    *Node    // 左辺
	rhs    *Node    // 右辺
	val    string   // せいすうち
	offset int      // ローカル変数のオフセット
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

// ローカル変数よう
func new_node_local_variance() *Node {
	var node *Node
	node = new(Node)
	node.kind = ND_LVAR
	// node.offset = int()
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
	node = expr()
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

// 現状見ている場所
var now_loc int

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

		// == / <= / >= 演算子
		if len(str) >= 2 && (str[:2] == "==" || str[:2] == "!=" || str[:2] == "<=" || str[:2] == ">=") {
			cur = new_token(TK_RESERVED, cur, str[:2])
			str = str[2:]
			now_loc += 2
			continue
		}

		// + / - / */ / / () / < / > 演算子
		if str[0] == '+' || str[0] == '-' || str[0] == '*' || str[0] == '/' || str[0] == '(' || str[0] == ')' || str[0] == '>' || str[0] == '<' || str[0] == ';' || str[0] == '=' {
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
