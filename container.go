// ベクタ、マップ、テストコードとか

package main

import (
	"fmt"
	"os"
)

// えらー
func Error(str string) {
	tmp := ""
	for i := 0; i < now_loc; i++ {
		tmp += " "
	}
	fmt.Fprintln(os.Stderr, user_input)
	fmt.Fprintln(os.Stderr, tmp+"^ "+str)
	os.Exit(1)
}
