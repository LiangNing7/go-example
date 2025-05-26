package main

import (
	"fmt"
	"runtime"
	"strconv"
)

func printStack(skip int) {
	var pcs [30]uintptr
	n := runtime.Callers(skip, pcs[:])

	for i := 0; i < n; i++ {
		pc := pcs[i]
		fn := runtime.FuncForPC(pc - 1)
		file, line := fn.FileLine(pc - 1)
		fmt.Printf("Func Name: %s\n", fn.Name())
		fmt.Printf("File: %s, Line: %s\n\n", file, strconv.Itoa(line))
	}
}

func Print(skip int) {
	printStack(skip)
}

func main() {
	Print(0)

	fmt.Println("============================================")

	Print(3)
}
