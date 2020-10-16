package main

import "fmt"

func (n *node) gen() {
	switch n.kind {
	case ndNone:
		return
	case ndNum:
		fmt.Println("  push " + n.val)
		return
	case ndLvar:
		n.genVariance()
		fmt.Println("  pop rax")
		fmt.Println("  mov rax, [rax]")
		fmt.Println("  push rax")
		return
	case ndAssign:
		n.genVariance()
		n.rhs.gen()

		fmt.Println("  pop rdi")
		fmt.Println("  pop rax")
		fmt.Println("  mov [rax], rdi")
		fmt.Println("  push rdi")
		return
	case ndRtn:
		n.lhs.gen()

		fmt.Println("  pop rax")
		fmt.Println("  mov rsp, rbp")
		fmt.Println("  pop rbp")
		fmt.Println("  ret")
		return
	case ndIf:
		n.lhs.gen()

		fmt.Println("  pop rax")
		fmt.Println("  cmp rax, 0")
		fmt.Println("  je " + n.label)
		n.rhs.gen()
		fmt.Println(n.label + ":")
		return
	case ndElse:
		n.lhs.lhs.gen()
		fmt.Println("  pop rax")
		fmt.Println("  cmp rax, 0")
		fmt.Println("  je " + n.label)
		n.lhs.rhs.gen()
		fmt.Println("  jmp " + n.lhs.label)
		fmt.Println(n.label + ":")
		n.rhs.gen()
		fmt.Println(n.lhs.label + ":")
		return
	case ndWhile:
		fmt.Println(n.label + "begin:")
		n.lhs.gen()
		fmt.Println("  pop rax")
		fmt.Println("  cmp rax, 0")
		fmt.Println("  je " + n.label + "end")
		n.rhs.gen()
		fmt.Println("  jmp " + n.label + "begin")
		fmt.Println(n.label + "end:")
		return
	case ndFor:
		n.lhs.lhs.gen()
		fmt.Println(n.label + "begin:")
		n.lhs.rhs.gen()
		if n.lhs.rhs.kind != ndNone {
			fmt.Println("  pop rax")
			fmt.Println("  cmp rax, 0")
			fmt.Println("  je" + n.label + "end")
		}
		n.rhs.lhs.gen()
		n.rhs.rhs.gen()
		fmt.Println("  jmp " + n.label + "begin")
		fmt.Println(n.label + "end:")
		return
	case ndBlk:
		n.lhs.gen()
		n.rhs.gen()
		return
	case ndFunc:
		return
	case ndArg:
		return
	}

	n.lhs.gen()
	n.rhs.gen()

	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	switch n.kind {
	case ndAdd:
		fmt.Println("  add rax, rdi")
	case ndSub:
		fmt.Println("  sub rax, rdi")
	case ndMul:
		fmt.Println("  imul rax, rdi")
	case ndDiv:
		fmt.Println("  cqo")
		fmt.Println("  idiv rdi")
	case ndEqu:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  sete al")
		fmt.Println("  movzb rax, al")
	case ndNeq:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setne al")
		fmt.Println("  movzb rax, al")
	case ndSml:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setl al")
		fmt.Println("  movzb rax, al")
	case ndEsm:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setle al")
		fmt.Println("  movzb rax, al")
	case ndBig:
		fmt.Println("  cmp rdi, rax")
		fmt.Println("  setl al")
		fmt.Println("  movzb rax, al")
	case ndEbg:
		fmt.Println("  cmp rdi, rax")
		fmt.Println("  setle al")
		fmt.Println("  movzb rax, al")
	}
	fmt.Println("  push rax")
}

func (n *node) genVariance() {
	if n.kind != ndLvar {
		error("代入の左辺地が変数ではありません")
	}
	fmt.Println("  mov rax, rbp")
	fmt.Println("  sub rax, ", n.offset)
	fmt.Println("  push rax")
}
