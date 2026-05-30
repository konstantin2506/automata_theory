package main

import (
	"log"
	"os"
)

func main() {
	filename := "main.yu"

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %s", err)
	}

	lexer := NewLexer(content)

	yyParse(lexer)
}
