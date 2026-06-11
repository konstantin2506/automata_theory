package main

import (
	"bufio"
	"log"
	"os"
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
}
