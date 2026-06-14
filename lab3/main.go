package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	intr "lab3/interpreter"
)

func main() {
	filename := "main.yu"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %s", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lexer := &lexerCtx{
		src: reader,
	}
	lexer.current = lexer.getc()
	yyParse(lexer)
	astRoot := GetRoot()
	if astRoot == nil {
		fmt.Printf("nil root\n")
		return
	}

	interpreter := intr.NewInterpreter(astRoot)
	intr.Interpret(interpreter)
}
