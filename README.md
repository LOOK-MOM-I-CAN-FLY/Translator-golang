# Simple Go Interpreter (ANTLR4)

Интерпретатор для упрощенного подмножества Go на базе ANTLR4. Проект включает грамматику, генерацию парсера и исполняемый интерпретатор.

Поддерживаемые конструкции:
- Объявление переменных `var x int = 5`
- Короткое объявление `x := 5`
- Типы `int`, `string`, `bool`
- Математические операции `+`, `-`, `*`, `/`
- Сравнения `==`, `!=`, `<`, `>`, `<=`, `>=`
- Логические операции `&&`, `||`, `!`
- `if-else`
- `for` (условный и `init; cond; post`)
- `fmt.Println()`
- `type Name struct { ... }`

Важно: в этом подмножестве **точка с запятой обязательна** после операторов/объявлений.

---

## Требования
- Go 1.21+
- Java 17+ (для генерации парсера ANTLR)
- Git

Зависимость Go:
```
github.com/antlr4-go/antlr/v4 v4.13.0
```

---

## Сборка (локально, через ANTLR4)

### 1. Клонировать репозиторий
```bash
git clone https://github.com/LOOK-MOM-I-CAN-FLY/Translator-golang.git
cd Translator-golang
```

### 2. Скачать ANTLR4
```bash
curl -L -o antlr-4.13.2-complete.jar https://www.antlr.org/download/antlr-4.13.2-complete.jar
```

### 3. Сгенерировать парсер (Go)
```bash
java -jar antlr-4.13.2-complete.jar -Dlanguage=Go -visitor -no-listener -package parser -o parser SimpleLexer.g4 SimpleParser.g4
```

### 4. Подтянуть зависимости и собрать
```bash
go mod tidy
go build -o translator .
```

Появится бинарник `./translator`.

---

## Использование

### REPL
```bash
./translator
```
Команды внутри REPL:
- `exit` — выход
- `run FILE` — запустить файл

### Запуск строки
```bash
./translator -c "var x int = 5; fmt.Println(x);"
```

### Запуск файла
```bash
./translator examples/example6_factorial.go
```

### Справка
```bash
./translator -h
```

---

## Пример кода
```go
var n int = 5;
var result int = 1;
var i int = 1;
for i <= n {
	result = result * i;
	i = i + 1;
}
fmt.Println("Factorial of 5 is:");
fmt.Println(result);
```

Ожидаемый вывод:
```
Factorial of 5 is:
120
```

---

## Docker

### Сборка образа
```bash
docker build -t translator .
```

### Запуск REPL внутри контейнера
```bash
docker run --rm -it translator
```

### Запуск строки
```bash
docker run --rm translator -c "var x int = 5; fmt.Println(x);"
```

### Запуск файла
```bash
docker run --rm -v "$PWD/examples":/root/examples translator /root/examples/example6_factorial.go
```

Контейнер сам:
- скачает ANTLR4,
- сгенерирует парсер,
- соберёт бинарник,
- запустит `./translator`.

---

## Структура проекта
- `SimpleLexer.g4`, `SimpleParser.g4` — грамматика
- `parser/` — **генерируется** ANTLR (в репозитории не хранится)
- `interpreter.go`, `main.go` — интерпретатор
- `examples/` — тестовые примеры входного языка
- `Dockerfile` — полностью рабочая сборка через Docker

Если меняете грамматику — пересоздайте `parser/` (см. шаг 3).
