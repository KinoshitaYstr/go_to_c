// コードジェネレータ

package main

import (
	"fmt"
)

// ノードからアセンブリ生成
func Gen(node *Node) {
	switch node.kind {
	case ND_NONE:
		return
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
		fmt.Println("  je " + node.label)
		Gen(node.rhs)
		fmt.Println(node.label + ":")
		return
	case ND_ELSE:
		Gen(node.lhs.lhs)
		fmt.Println("  pop rax")
		fmt.Println("  cmp rax, 0")
		fmt.Println("  je " + node.label)
		Gen(node.lhs.rhs)
		fmt.Println("  jmp " + node.lhs.label)
		fmt.Println(node.label + ":")
		Gen(node.rhs)
		fmt.Println(node.lhs.label + ":")
		return
	case ND_WHILE:
		fmt.Println(node.label + "begin:")
		Gen(node.lhs)
		fmt.Println("  pop rax")
		fmt.Println("  cmp rax, 0")
		fmt.Println("  je " + node.label + "end")
		Gen(node.rhs)
		fmt.Println("  jmp " + node.label + "begin")
		fmt.Println(node.label + "end:")

		return
	case ND_FOR:
		Gen(node.lhs.lhs)
		fmt.Println(node.label + "begin:")
		Gen(node.lhs.rhs)
		if node.lhs.rhs.kind != ND_NONE {
			fmt.Println("  pop rax")
			fmt.Println("  cmp rax, 0")
			fmt.Println("  je " + node.label + "end")
		}
		Gen(node.rhs.lhs)
		Gen(node.rhs.rhs)
		fmt.Println("  jmp " + node.label + "begin")
		fmt.Println(node.label + "end:")
		return
	case ND_BLOCK:
		Gen(node.lhs)
		Gen(node.rhs)
		return
	case ND_FUNC:
		fmt.Println("=====================")
		if node.lhs != nil {
			Gen(node.lhs)
		}
		fmt.Println(node.val + ":")

		fmt.Println("  nop")
		fmt.Println("  leave")
		fmt.Println("  ret")
		return
	case ND_ARG:
		if node.lhs != nil {
			Gen(node.lhs)
		}
		if node.rhs != nil {
			Gen(node.rhs)
		}
		fmt.Println(node)
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
