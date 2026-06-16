package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	game "lab3/game"
	intr "lab3/interpreter"

	"github.com/hajimehoshi/ebiten/v2"
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
	ebiten.SetWindowTitle("Шестиугольный лабиринт 30x30")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled) // разрешить изменение размера
	ebiten.SetFullscreen(true)                                     // полноэкранный режим
	gm := game.NewGame()
	ctrl := game.NewController(gm)

	go intr.Interpret(interpreter, &ctrl)
	if err := ebiten.RunGame(gm); err != nil {
		log.Fatal(err)
	}
}
