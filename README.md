# ⏱️ Subtle Race

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Concurrency](https://img.shields.io/badge/concurrency-timer%20loop-success?style=flat)](https://go.dev/doc/effective_go#concurrency)
[![Race-safe](https://img.shields.io/badge/race-safe-brightgreen?style=flat)](https://go.dev/doc/articles/race_detector)

Реализация cron-подобного планировщика, который многократно
вызывает `action()`, ожидая перед каждым запуском `next()` времени

## 🎯 Интерфейс

```go
type Cron interface {
	Run(ctx context.Context, action func(), next func() time.Duration)
}
```

`Run`:
- ждёт `next()` перед каждым вызовом `action`;
- прекращает вызовы `action` после `ctx.Done()`.

## ❗ Проблемная реализация

Мейнтейнеры Go тоже пытались реализовать такой интерфейс. Сходу у них получилась следующая реализация:

```go
func (c *cronImpl) Run(ctx context.Context, action func(), next func() time.Duration) {
    var t *time.Timer
    
    t = time.AfterFunc(next(), func() {
        select {
        case <-ctx.Done():
            return
        default:
            action()
            t.Reset(next())
        }
    })

    <-ctx.Done()
}
```

Проблема в том, что такая реализация приводит к data race, который довольно сложно отловить

## 🧪 Зачем нужны `lite_test` и `hard_test`

Две версии тестов в данном проекте нужны, чтобы показать, как трудно найти эту ошибку 

### `lite_test`

Да же с включённым race detector не отловят data race

### `hard_test`

Делают очень большое количество запусков и там, с включённым race detector, становится видна гонка

## ▶️ Как запускать тесты для двух реализаций

В коде есть две реализации:
- `New()` — основная, исправленная (`time.NewTimer`);
- `NewAfterFunc()` — та самая версия на `time.AfterFunc`, добавленная для демонстрации проблемы.

Тесты выбирают реализацию через `CRON_IMPL`:
- `CRON_IMPL=fixed` (или пусто) → `New()`;
- `CRON_IMPL=afterfunc` → `NewAfterFunc()`.

Для fixed версии:

```bash
CRON_IMPL=fixed go test -v -count=1 ./internal/cron -run '^TestLite'
CRON_IMPL=fixed go test -race -v -count=1 ./internal/cron -run '^TestLite'
CRON_IMPL=fixed go test -v -count=1 ./internal/cron -run '^TestHard'
CRON_IMPL=fixed go test -race -v -count=1 ./internal/cron -run '^TestHard'
```

Для afterfunc версии:

```bash
CRON_IMPL=afterfunc go test -v -count=1 ./internal/cron -run '^TestLite'
CRON_IMPL=afterfunc go test -race -v -count=1 ./internal/cron -run '^TestLite'
CRON_IMPL=afterfunc go test -v -count=1 ./internal/cron -run '^TestHard'
CRON_IMPL=afterfunc go test -race -v -count=1 ./internal/cron -run '^TestHard'
```

## 🚀 Пример использования

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Koval_Dmitrii/subtlerace/internal/cron"
)

func main() {
	c := cron.New()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	count := 0
	c.Run(ctx, func() {
		count++
		fmt.Println("tick:", count)
	}, func() time.Duration {
		if count < 5 {
			return 200 * time.Millisecond
		}

		return 500 * time.Millisecond
	})
}
```

## 🏗️ Структура проекта

```text
subtlerace/
├── internal/
│   └── cron/
│       ├── cron.go
│       ├── lite_test.go
│       └── hard_test.go
├── go.mod
└── README.md
```

## 👨‍💻 Автор

**Коваль Дмитрий**

- GitHub: [@Koval-Dmitrii](https://github.com/Koval-Dmitrii)
