package main

import "fmt"

var labels label
var locals *variance

func (t *token) program() []*node {
	// fmt.Println(t)
	// fmt.Println(t.consume("{"))
	// fmt.Println(t)
	var n []*node
	locals.init()
	fmt.Println(t.stmt())
	for t.atEOF() {
		n = append(n, t.stmt())
	}
	fmt.Println(n)
	return n
}

func (t *token) stmt() *node {
	// fmt.Println(t)
	var n *node
	var errFlag bool

	if t, errFlag = t.consume("{"); errFlag {
		t, errFlag = t.consume("}")
		if errFlag {
			return newNodeNone()
		}
		n = t.stmt()
		t, errFlag = t.consume("}")
		for !errFlag {
			n.updateNode(ndBlk, n, t.stmt())
		}
		return n
	} else if t, errFlag = t.consume("return"); errFlag {
		n = new(node)
		n.kind = ndRtn
		n.lhs = t.expr()
		t = t.expect(";")
		return n
	} else if t, errFlag = t.consume("if"); errFlag {
		t = t.expect("(")
		labels.nextIf()
		n = new(node)
		n.kind = ndIf
		n.label = labels.getIf()
		n.lhs = t.expr()
		t = t.expect(")")
		n.rhs = t.stmt()
		if t, errFlag = t.consume("else"); errFlag {
			labels.nextElse()
			n.updateNode(ndElse, n, t.stmt())
			n.label = labels.getElse()
		}
		return n
	} else if t, errFlag = t.consume("while"); errFlag {
		t = t.expect("(")
		labels.nextWhile()
		n = new(node)
		n.kind = ndWhile
		n.label = labels.getWhile()
		n.lhs = t.expr()
		t = t.expect(")")
		n.rhs = t.stmt()
		return n
	} else if t, errFlag = t.consume("for"); errFlag {
		t = t.expect("(")
		labels.nextFor()
		n = new(node)
		n.kind = ndFor
		n.label = labels.getFor()
		n.lhs = new(node)
		n.rhs = new(node)
		if t, errFlag = t.consume(";"); errFlag {
			n.lhs.lhs = newNodeNone()
		} else {
			n.lhs.lhs = t.expr()
			t = t.expect(";")
		}
		if t, errFlag = t.consume(";"); errFlag {
			n.lhs.rhs = newNodeNone()
		} else {
			n.lhs.rhs = t.expr()
			t = t.expect(";")
		}
		if t, errFlag = t.consume(")"); errFlag {
			n.rhs.lhs = newNodeNone()
		} else {
			n.rhs.lhs = t.expr()
			t = t.expect(")")
		}
		n.rhs.rhs = t.stmt()
		return n
	} else {
		n = t.expr()
		t = t.expect(";")
		return n
	}
}

func (t *token) expr() *node {
	// fmt.Println(t)
	return t.assign()
}

func (t *token) assign() *node {
	// fmt.Println(t)
	var n *node
	var errFlag bool
	n = t.equality()
	if t, errFlag = t.consume("="); errFlag {
		n.updateNode(ndAssign, n, t.assign())
	}
	return n
}

func (t *token) equality() *node {
	// fmt.Println(t)
	var n *node
	var errFlag bool
	n = t.relational()
	for {
		if t, errFlag = t.consume("=="); errFlag {
			n.updateNode(ndEqu, n, t.relational())
		} else if t, errFlag = t.consume("!="); errFlag {
			n.updateNode(ndNeq, n, t.relational())
		} else {
			return n
		}
	}
}

func (t *token) relational() *node {
	// fmt.Println(t)
	var n *node
	var errFlag bool
	n = t.add()
	for {
		if t, errFlag = t.consume("<"); errFlag {
			n.updateNode(ndSml, n, t.add())
		} else if t, errFlag = t.consume("<="); errFlag {
			n.updateNode(ndEsm, n, t.add())
		} else if t, errFlag = t.consume(">"); errFlag {
			n.updateNode(ndBig, n, t.add())
		} else if t, errFlag = t.consume(">="); errFlag {
			n.updateNode(ndEbg, n, t.add())
		} else {
			return n
		}
	}
}

func (t *token) add() *node {
	// fmt.Println(t)
	var n *node
	var errFlag bool
	n = t.mul()
	for {
		if t, errFlag = t.consume("+"); errFlag {
			n.updateNode(ndAdd, n, t.mul())
		} else if t, errFlag = t.consume("-"); errFlag {
			n.updateNode(ndSub, n, t.mul())
		} else {
			return n
		}
	}
}

func (t *token) mul() *node {
	// fmt.Println(t)
	var n *node
	var errFlag bool
	n = t.unary()
	for {
		if t, errFlag = t.consume("*"); errFlag {
			n.updateNode(ndMul, n, t.unary())
		} else if t, errFlag = t.consume("/"); errFlag {
			n.updateNode(ndDiv, n, t.unary())
		} else {
			return n
		}
	}
}

func (t *token) unary() *node {
	// fmt.Println(t)
	var errFlag bool
	if t, errFlag = t.consume("+"); errFlag {
		return t.primary()
	}
	if t, errFlag = t.consume("-"); errFlag {
		var n *node
		n = new(node)
		n.updateNode(ndSub, newNodeNum("0"), t.primary())
		return n
	}
	return t.primary()
}

func (t *token) primary() *node {
	// fmt.Println(t)
	var n *node
	var errFlag bool
	var tk *token
	n = new(node)
	t, tk, errFlag = t.consumeIdent()
	if errFlag {
		n.kind = ndLvar
		lvar := locals.findVar(tk)
		if lvar != nil {
			n.offset = lvar.offset
		} else {
			lvar = new(variance)
			lvar.next = locals
			lvar.name = tk.str
			lvar.offset = locals.offset + 8
			n.offset = lvar.offset
			locals = lvar
		}
		return n
	}
	if t, errFlag = t.consumeNone(); errFlag {
		n = newNodeNone()
		return n
	}
	var str string
	t, str = t.expectNumber()
	return newNodeNum(str)
}

type variance struct {
	next   *variance
	name   string
	offset int
}

func (v *variance) init() {
	v = new(variance)
	v.offset = 0
}

func (v *variance) getSpace() int {
	return v.offset
}

func (v *variance) findVar(t *token) *variance {
	for va := v; va != nil; va = va.next {
		if va.name == t.str {
			return va
		}
	}
	return nil
}

type label struct {
	ifLabel    bool
	elseLabel  bool
	whileLabel bool
	forLabel   bool
	ifNum      int
	elseNum    int
	whileNum   int
	forNum     int
}

func (l *label) getIf() string {
	if l.ifLabel {
		l.ifLabel = true
		return fmt.Sprintf(".Lif%02d", l.ifNum)
	}
	return ""
}

func (l *label) getElse() string {
	if l.elseLabel {
		l.elseLabel = true
		return fmt.Sprintf(".Lelse%02d", l.elseNum)
	}
	return ""
}

func (l *label) getWhile() string {
	if l.whileLabel {
		l.whileLabel = true
		return fmt.Sprintf(".Lwhile%02d", l.whileNum)
	}
	return ""
}

func (l *label) getFor() string {
	if l.forLabel {
		l.forLabel = true
		return fmt.Sprintf(".Lfor%02d", l.forNum)
	}
	return ""
}

func (l *label) nextIf() {
	if l.ifLabel {
		l.ifNum++
		l.ifLabel = false
	}
}

func (l *label) nextElse() {
	if l.elseLabel {
		l.elseNum++
		l.elseLabel = false
	}
}

func (l *label) nextWhile() {
	if l.whileLabel {
		l.whileNum++
		l.whileLabel = false
	}
}

func (l *label) nextFor() {
	if l.forLabel {
		l.forNum++
		l.forLabel = false
	}
}

// ノード種類変数
type nodeKind int

// 抽象木のノードの種類
const (
	ndAdd    nodeKind = iota // +
	ndSub                    // -
	ndMul                    // *
	ndDiv                    // /
	ndNum                    // number
	ndEqu                    // ==
	ndNeq                    // !=
	ndBig                    // >
	ndSml                    // <
	ndEbg                    // >=
	ndEsm                    // <=
	ndAssign                 // =
	ndLvar                   // local var
	ndRtn                    // return
	ndIf                     // if
	ndElse                   // else
	ndWhile                  // while
	ndFor                    // for
	ndBlk                    // block
	ndFunc                   // function
	ndArg                    // arg
	ndNone                   // none
)

type node struct {
	kind   nodeKind // ノードの方
	lhs    *node    // 左辺
	rhs    *node    // 右辺
	val    string   // 値
	offset int      // オフセット
	label  string   // ラベル
}

// ノードの更新
func (n *node) updateNode(kind nodeKind, lhs *node, rhs *node) {
	var newNode *node
	newNode = new(node)
	newNode.kind = kind
	newNode.lhs = lhs
	newNode.rhs = rhs
	n = newNode
}

func newNodeNum(val string) *node {
	var n *node
	n = new(node)
	n.kind = ndNum
	n.val = val
	return n
}

func newNodeNone() *node {
	var n *node
	n = new(node)
	n.kind = ndNone
	return n
}

func newLocalVariance() *node {
	var n *node
	n = new(node)
	n.kind = ndLvar
	return n
}

// トークンの種類用変数
type tokenKind int

// 種類の定位数
const (
	tkReserve tokenKind = iota // 予約時
	tkIdent                    // 識別子
	tkNum                      // 整数値
	tkEOF                      // 入力終了用
)

// 現状のトークン
type token struct {
	kind tokenKind // トークンの種類
	next *token    // 次のトークン
	str  string    // 値
}

// 次のトークン消費
func (t *token) consume(op string) (*token, bool) {
	if t.kind != tkReserve || t.str != op {
		return t, false
	}
	return t.next, true
}

// 変数消費
func (t *token) consumeIdent() (*token, *token, bool) {
	if t.kind != tkIdent {
		return t, nil, false
	}
	return t.next, t, true
}

func (t *token) consumeNone() (*token, bool) {
	if t != nil {
		return t, false
	}
	return t.next, true
}

// 確認
func (t *token) expect(op string) *token {
	if t.kind != tkReserve || t.str != op {
		error(op + "ではありません")
	}
	return t.next
}

func (t *token) expectNumber() (*token, string) {
	if t.kind != tkNum {
		error("数値ではありません")
	}
	val := t.str
	return t.next, val
}

func (t *token) atEOF() bool {
	return t.kind == tkEOF
}

// とーくないずする
func tokenize(str string) *token {
	var head token
	head.next = nil
	cur := &head

	var s sentence
	s.init(str)

	for s.isValid() {
		// 空白飛ばし
		s.deleteSpace()

		// return文
		if s.checkKeyword("return") {
			cur = s.createReserveToken(cur, "return")
			continue
		}

		// if文
		if s.checkKeyword("if") {
			cur = s.createReserveToken(cur, "if")
			continue
		}

		// else文
		if s.checkKeyword("else") {
			cur = s.createReserveToken(cur, "else")
			continue
		}

		// while文
		if s.checkKeyword("while") {
			cur = s.createReserveToken(cur, "while")
			continue
		}

		// for文
		if s.checkKeyword("for") {
			cur = s.createReserveToken(cur, "for")
			continue
		}

		// == / <= / >= 演算子
		if s.checkKeyword("==") || s.checkKeyword("<=") || s.checkKeyword(">=") {
			cur = s.createReserveToken(cur, s.nowString[:2])
			continue
		}

		// + / - / * / / / ( / ) / < / > 演算子
		if s.checkKeyword("+") || s.checkKeyword("-") || s.checkKeyword("*") || s.checkKeyword("/") || s.checkKeyword("(") || s.checkKeyword(")") || s.checkKeyword(">") || s.checkKeyword("<") || s.checkKeyword(";") || s.checkKeyword("=") || s.checkKeyword("{") || s.checkKeyword("}") || s.checkKeyword(",") {
			cur = s.createReserveToken(cur, string(s.nowString[0]))
			continue
		}

		// 数値処理
		if s.checkNumber() {
			cur = s.createNumberToken(cur)
			continue
		}

		// 変数
		if s.checkAlphabet() {
			cur = s.createIdentToken(cur)
			continue
		}

		error("トークナイズできません")
	}

	// EOF
	cur = s.createEOFToken(cur)

	return head.next
}

// 文関係
type sentence struct {
	original  string // 元の
	nowString string // 今の読み込み
}

// 初期化
func (s *sentence) init(str string) {
	s.original = str
	s.nowString = str
}

// 余白飛ばし
func (s *sentence) deleteSpace() {
	for s.checkKeyword(" ") {
		s.nowString = s.nowString[1:]
	}
}

// キーワード確認
func (s *sentence) checkKeyword(keyword string) bool {
	if len(keyword) == len(s.nowString) {
		return keyword == s.nowString
	} else if len(s.nowString) > len(keyword) {
		return s.nowString[0:len(keyword)] == keyword
	}
	return false
}

// 数値確認
func (s *sentence) checkNumber() bool {
	if len(s.nowString) > 0 && '0' <= s.nowString[0] && s.nowString[0] <= '9' {
		return true
	}
	return false
}

// アルファベット確認
func (s *sentence) checkAlphabet() bool {
	if len(s.nowString) <= 0 {
		return false
	} else if ('a' <= s.nowString[0] && s.nowString[0] <= 'z') || ('A' <= s.nowString[0] && s.nowString[0] <= 'Z') {
		return true
	}
	return false
}

// トークン作成
// 予約碁盤
func (s *sentence) createReserveToken(cur *token, keyword string) *token {
	var t *token
	t = new(token)
	t.kind = tkReserve
	t.str = keyword
	t.next = nil
	cur.next = t
	s.nowString = s.nowString[len(keyword):]
	return t
}

// 数値版
func (s *sentence) createNumberToken(cur *token) *token {
	var t *token
	t = new(token)
	t.kind = tkNum
	t.str = ""
	for {
		if s.checkNumber() {
			t.str += string(s.nowString[0])
			s.nowString = s.nowString[1:]
		} else {
			break
		}
	}
	t.next = nil
	cur.next = t
	return t
}

// 変数版
func (s *sentence) createIdentToken(cur *token) *token {
	var t *token
	t = new(token)
	t.kind = tkIdent
	t.str = ""
	for len(s.nowString) > 0 {
		if s.checkAlphabet() {
			t.str += string(s.nowString[0])
			s.nowString = s.nowString[1:]
		} else {
			break
		}
	}
	t.next = nil
	cur.next = t
	return t
}

// EOF版
func (s *sentence) createEOFToken(cur *token) *token {
	var t *token
	t = new(token)
	t.kind = tkEOF
	t.next = nil
	cur.next = t
	return t
}

// 生存確認
func (s *sentence) isValid() bool {
	return len(s.nowString) > 0
}
