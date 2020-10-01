// コードジェネレータ

package main

import "fmt"

// ノードからアセンブリ生成
func Gen(node *Node) {
	if node.kind == ND_NUM {
		fmt.Println("  push " + node.val)
		return
	}

	Gen(node.lhs)
	Gen(node.rhs)

	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	switch node.kind {
	case ND_ADD:
		fmt.Println("  add rax, rdi")
	case ND_SUB:
		fmt.Println("  sub rax, rdi")
	case ND_MUL:
		fmt.Println("  imul rax, rdi")
	case ND_DIV:
		fmt.Println("  cqo")
		fmt.Println("  idiv rdi")
	case ND_EQU:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  sete al")
		fmt.Println("  movzb rax, al")
	case ND_NEQ:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setne al")
		fmt.Println("  movzb rax, al")
	case ND_SML:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setl al")
		fmt.Println("  movzb rax, al")
	case ND_ESM:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setle al")
		fmt.Println("  movzb rax, al")
	case ND_BIG:
		fmt.Println("  cmp rdi, rax")
		fmt.Println("  setl al")
		fmt.Println("  movzb rax, al")
	case ND_EBG:
		fmt.Println("  cmp rdi, rax")
		fmt.Println("  setle al")
		fmt.Println("  movzb rax, al")
	}

	fmt.Println("  push rax")
}
