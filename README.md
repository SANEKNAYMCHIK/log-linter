# log-linter

Линтер для проверки лог-записей в Go, совместимый с `golangci-lint`.
\
Поддерживает `log/slog` и `go.uber.org/zap`.

## Правила

- **Строчная буква в начале** – сообщение не должно начинаться с заглавной буквы.
- **Английский язык** – сообщение должно быть только на английском языке.
- **Эмодзи и спецсимволы** – запрещены.
- **Чувствительные данные** – проверка по ключевым словам (по умолчанию: password, secret, token, key) и регулярным выражениям.

## Установка

### Как отдельный бинарник

```bash
go install github.com/SANEKNAYMCHIK/log-linter/cmd/log-linter@latest
```

Затем запуск:

```bash
log-linter ./...
```

### Интеграция с golangci-lint (custom linter)
Соберите бинарник:

```bash
go build -o loglint ./cmd/log-linter
```
В .golangci.yml добавьте секцию:

```yaml
linters-settings:
  custom:
    loglint:
      path: ./loglint
      description: Линтер для логов
      original-url: github.com/yourname/log-linter
linters:
  enable:
    - loglint
```

Запустите golangci-lint run.

## Конфигурация
Линтер поддерживает флаги командной строки (передаются через golangci-lint в секции settings):

**sensitive-words** – список слов через запятую (по умолчанию password,secret,token,key)

**sensitive-patterns** – список регулярных выражений для поиска чувствительных данных (например, \b\d{16}\b для номеров карт)


#### Автоисправление
Правило lowercase поддерживает автоматическое исправление. 

Запустите 
```bash
golangci-lint run --fix
```
чтобы линтер сам исправил сообщения с заглавной буквы.

#### Разработка
Запуск тестов
```bash
go test ./...
```
Добавление новых правил

Поместите новое правило в пакет internal/rules и вызовите его в analyzer.go:run.

#### CI/CD
В репозитории есть GitHub Actions workflow (.github/workflows/ci.yml), который запускает тесты и сборку.
