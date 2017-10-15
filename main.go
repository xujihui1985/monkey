package main

import (
	"monkey/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
