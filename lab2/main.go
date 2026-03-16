package main

import (
	"fmt"

	kregex "kregex/kregex"
)

func main() {
	fmt.Println("Hello from kregex")
	_, err := kregex.BuildAST("(()(hello thesre)(()))")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("OK")
	}
}
