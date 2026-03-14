import random
import re
import os

def generate_valid_string(target_length):
    """
    Генерирует строку, соответствующую паттерну:
    ^(\d+)\s+([a-zA-Z][a-zA-Z0-9]{0,15})\s*=\s*(\-\d+|\d+|[a-zA-Z][a-zA-Z0-9]{0,15})\s*([+\-*/]\s*(\-\d+|\d+|[a-zA-Z][a-zA-Z0-9]{0,15}))?\s*$
    """

    # Вспомогательные функции для генерации компонентов
    def generate_number(allow_negative=True):
        """Генерирует число (возможно отрицательное)"""
        if allow_negative and random.choice([True, False]):
            return f"-{random.randint(1, 999)}"
        return str(random.randint(0, 999))

    def generate_identifier():
        """Генерирует идентификатор: буква + буквы/цифры (0-15 символов)"""
        first = random.choice('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ')
        length = random.randint(0, 15)
        rest = ''.join(random.choices('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', k=length))
        return first + rest

    def generate_operand():
        """Генерирует операнд (число или идентификатор)"""
        if random.choice([True, False]):
            return generate_number()
        else:
            return generate_identifier()

    # Генерируем базовые компоненты
    while True:
        # Первое число
        num1 = generate_number(allow_negative=False)  # только положительное

        # Пробелы после первого числа (минимум 1)
        spaces1 = ' ' * random.randint(1, 5)

        # Первый идентификатор (после числа)
        id1 = generate_identifier()

        # Пробелы вокруг =
        spaces_before_eq = ' ' * random.randint(0, 3)
        spaces_after_eq = ' ' * random.randint(0, 3)

        # Правый операнд
        operand1 = generate_operand()

        # Решаем, добавлять ли операцию
        has_operation = random.choice([True, False])

        if has_operation:
            # Пробелы перед операцией
            spaces_before_op = ' ' * random.randint(0, 3)

            # Операция
            op = random.choice(['+', '-', '*', '/'])

            # Пробелы после операции
            spaces_after_op = ' ' * random.randint(0, 3)

            # Второй операнд
            operand2 = generate_operand()

            # Пробелы в конце
            trailing_spaces = ' ' * random.randint(0, 3)

            # Собираем строку
            line = f"{num1}{spaces1}{id1}{spaces_before_eq}={spaces_after_eq}{operand1}{spaces_before_op}{op}{spaces_after_op}{operand2}{trailing_spaces}"
        else:
            # Пробелы в конце
            trailing_spaces = ' ' * random.randint(0, 3)

            # Собираем строку без операции
            line = f"{num1}{spaces1}{id1}{spaces_before_eq}={spaces_after_eq}{operand1}{trailing_spaces}"

        # Проверяем, что строка соответствует паттерну
        pattern = re.compile(r"^(\d+)\s+([a-zA-Z][a-zA-Z0-9]{0,15})\s*=\s*(\-\d+|\d+|[a-zA-Z][a-zA-Z0-9]{0,15})\s*([+\-*/]\s*(\-\d+|\d+|[a-zA-Z][a-zA-Z0-9]{0,15}))?\s*$")
        if pattern.match(line):
            return line

def generate_file_with_length(filename, target_len):
    """
    Генерирует файл с одной строкой указанной длины
    """
    print(f"Генерация строки длиной {target_len} символов в файл {filename}...")

    # Генерируем строки, пока не получим нужную длину
    attempts = 0
    max_attempts = 1000

    while attempts < max_attempts:
        line = generate_valid_string(target_len)
        current_len = len(line)

        if current_len == target_len:
            with open(filename, 'w') as f:
                f.write(line + '\n')
            print(f"  Успешно! Длина: {current_len}")
            return True
        elif current_len < target_len:
            # Если строка короче, пробуем добавить пробелы
            needed_spaces = target_len - current_len
            # Добавляем пробелы в разные места
            parts = line.split('=')
            if len(parts) == 2:
                # Вставляем пробелы перед =, после = и в конце
                left = parts[0]
                right = parts[1]

                # Распределяем пробелы
                spaces_before = random.randint(0, needed_spaces)
                spaces_after = random.randint(0, needed_spaces - spaces_before)
                spaces_end = needed_spaces - spaces_before - spaces_after

                new_line = f"{left}{' ' * spaces_before}={' ' * spaces_after}{right}{' ' * spaces_end}"

                pattern = re.compile(r"^(\d+)\s+([a-zA-Z][a-zA-Z0-9]{0,15})\s*=\s*(\-\d+|\d+|[a-zA-Z][a-zA-Z0-9]{0,15})\s*([+\-*/]\s*(\-\d+|\d+|[a-zA-Z][a-zA-Z0-9]{0,15}))?\s*$")
                if len(new_line) == target_len and pattern.match(new_line):
                    with open(filename, 'w') as f:
                        f.write(new_line + '\n')
                    print(f"  Успешно после добавления пробелов! Длина: {len(new_line)}")
                    return True

        attempts += 1

    if attempts >= max_attempts:
        print(f"  Не удалось сгенерировать строку длины {target_len} после {max_attempts} попыток")
        # Записываем приблизительную строку
        with open(filename, 'w') as f:
            f.write(line + f" {' ' * (target_len - len(line) - 1)}" + '\n')
        return False

def ensure_directory_exists(directory):
    """Создаёт директорию, если она не существует"""
    if not os.path.exists(directory):
        os.makedirs(directory)

# Основная программа
if __name__ == "__main__":
    lengths = [1000, 5000, 10000, 20000, 30000, 40000, 50000, 60000, 70000, 80000, 90000, 100000]

    # Создаём директорию для выходных файлов
    output_dir = "generated_files"
    ensure_directory_exists(output_dir)

    print(f"Генерация строк в директорию {output_dir}/")

    success_count = 0
    for length in lengths:
        filename = os.path.join(output_dir, f"{length}.txt")
        if generate_file_with_length(filename, length):
            success_count += 1

    print(f"\nГотово! Успешно сгенерировано {success_count} из {len(lengths)} файлов")
    print(f"Файлы сохранены в директории: {output_dir}/")