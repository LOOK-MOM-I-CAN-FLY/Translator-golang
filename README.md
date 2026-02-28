# Simple Go Interpreter 

## Описание

**Simple Go Interpreter** — это интерпретатор для упрощенного подмножества языка программирования Go, разработанный с использованием ANTLR4.

Интерпретатор поддерживает следующие конструкции:
- Объявление переменных (`var x int = 5`)
- Короткое объявление (`x := 5`)
- Базовые типы данных: `int`, `string`, `bool`
- Математические операции: `+`, `-`, `*`, `/`
- Операции сравнения: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Логические операции: `&&`, `||`, `!`
- Условные конструкции: `if-else`
- Циклы: `for`
- Вывод на экран: `fmt.Println()`


## Требования к установке

### Системные требования
- Go 1.21 или выше
- Git

### Зависимости
Проект использует ANTLR4 Go runtime:
```
github.com/antlr4-go/antlr/v4 v4.13.0
```

## Инструкции по сборке

### 1. Клонирование/получение проекта
```bash
cd translator
```

### 2. Загрузка зависимостей
```bash
go mod download
go mod tidy
```

### 3. Компиляция (сборка)
```bash
go build -o translator main.go
```

Исполняемый файл `translator` (или `translator.exe` на Windows) будет создан в текущем каталоге.

### 4. Проверка сборки
```bash
# Вывести справку
./translator

# Должно появиться:
# Simple Go Interpreter
# 
# Usage:
#   translator                    - Start REPL
#   translator FILE               - Run file
#   translator -c CODE            - Run code
```


## Инструкции по использованию

### Запуск кода из командной строки
```bash
./translator -c "var x int = 5; fmt.Println(x);"
```

### Запуск из файла
```bash
./translator examples/example1_simple_var.go
```

### Интерактивный режим (REPL)
```bash
./translator
```

В REPL-режиме доступны команды:
- `exit` — выход из интерпретатора
- `run FILE` — запустить файл

Пример использования REPL:
```
> var x int = 10;
> fmt.Println(x);
10
> x = x + 5;
> fmt.Println(x);
15
> exit
```

## Примеры использования

### Пример 1: Простое объявление и вывод
```go
var x int = 42;
fmt.Println(x);
```
Вывод: `42`

### Пример 2: Математические операции
```go
var a int = 10;
var b int = 5;
var sum int = a + b;
fmt.Println(sum);
```
Вывод: `15`


