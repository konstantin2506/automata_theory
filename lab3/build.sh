#!/bin/bash

# Настройки путей (используем $HOME вместо ~ для надежности)
GOLEX="$HOME/go/bin/golex"
GOYACC="$HOME/go/bin/goyacc"

# Функция для проверки существования утилит
check_dependency() {
    if [ ! -x "$1" ]; then
        echo "Ошибка: утилита не найдена или не исполняема: $1" >&2
        exit 1
    fi
}

# Проверяем зависимости перед запуском
check_dependency "$GOLEX"
check_dependency "$GOYACC"

# Генерация файлов (переменные вызываются через $ без скобок)
"$GOLEX" -o lexer.go lexer.l &&
    "$GOYACC" -o parser.go parser.y

# Проверка успешности генерации
if [ $? -eq 0 ]; then
    echo "Success"
    go run .
else
    echo "Error while generating" >&2
    exit 1
fi
