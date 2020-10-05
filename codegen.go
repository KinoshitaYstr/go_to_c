// コードジェネレータ

package main

import (
	"fmt"
)

// ノードからアセンブリ生成
func Gen(node *Node) {
	switch node.kind {
	case ND_NUM:
		fmt.Println("  push " + node.val)
		return
	case ND_LVAR:
		gen_lvar(node)
		fmt.Println("  pop rax")
		fmt.Println("  mov rax, [rax]")
		fmt.Println("  push rax")
		return
	case ND_ASSIGN:
		gen_lvar(node.lhs)
		Gen(node.rhs)

		fmt.Println("  pop rdi")
		fmt.Println("  pop rax")
		fmt.Println("  mov [rax], rdi")
		fmt.Println("  push rdi")
		return
	case ND_RETURN:
		Gen(node.lhs)

		fmt.Println("  pop rax")
		fmt.Println("  mov rsp, rbp")
		fmt.Println("  pop rbp")
		fmt.Println("  ret")
		return
	case ND_IF:
		Gen(node.lhs)

		fmt.Println("  pop rax")
		fmt.Println("  cmp rax, 0")
		fmt.Println("  je .Lend00")
		Gen(node.rhs)
		fmt.Println(".Lend00:")

		return
	case ND_ELSE:
		fmt.Println("  pop rax")

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

// 左辺値用
func gen_lvar(node *Node) {
	if node.kind != ND_LVAR {
		Error("代入の左辺値が変数ではありません")
	}
	fmt.Println("  mov rax, rbp")
	fmt.Println("  sub rax,", node.offset)
	fmt.Println("  push rax")
}
