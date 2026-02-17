# Руководство по использованию библиотеки slog в Go

## Введение

Библиотека `slog` (structured logging) - это современный способ логирования в Go, представленный в версии 1.21. Она предоставляет структурированное логирование с возможностью выбора обработчика (handler) и форматирования вывода.

## Установка

Библиотека `slog` является частью стандартной библиотеки Go, поэтому для ее использования не требуется установка дополнительных пакетов.

## Основные концепции

### 1. Логгер (Logger)

Логгер - это основной объект для записи логов. Он имеет уровень логирования (log level) и обработчик (handler).

```go
import "log/slog"

// Создание нового логгера
logger := slog.New(handler)
```

### 2. Уровни логирования

slog поддерживает следующие уровни логирования (от низкого к высокому):

- `slog.LevelDebug`: Debug сообщения
- `slog.LevelInfo`: Информационные сообщения
- `slog.LevelWarn`: Предупреждающие сообщения
- `slog.LevelError`: Сообщения об ошибках

### 3. Обработчики (Handlers)

Обработчики определяют, куда и как будут записываться логи. Встроенные обработчики:

- `slog.NewTextHandler`: Форматированный вывод в текстовом формате
- `slog.NewJSONHandler`: Вывод в формате JSON

## Примеры использования

### 1. Создание базового логгера

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	// Создание текстового обработчика
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	
	// Создание логгера
	logger := slog.New(handler)
	
	// Использование логгера
	logger.Info("Приложение запущено", "version", "1.0.0")
	logger.Debug("Отладочная информация")
	logger.Warn("Предупреждение")
	logger.Error("Ошибка", "err", "что-то пошло не так")
}
```

### 2. Настройка уровня логирования

```go
// Создание обработчика с уровнем Info
handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo, // Только Info, Warn и Error
})

logger := slog.New(handler)

logger.Debug("Это сообщение не будет выведено") // Не выводится
logger.Info("Информационное сообщение")         // Выводится
logger.Warn("Предупреждение")                  // Выводится
logger.Error("Ошибка")                         // Выводится
```

### 3. Использование JSON обработчика

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	// Создание JSON обработчика
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	
	logger := slog.New(handler)
	
	// Запись логов
	logger.Info("Приложение запущено", "version", "1.0.0", "service", "auth")
	logger.Error("Ошибка аутентификации", "user_id", "123", "error", "invalid password")
}
```

### 4. Группировка атрибутов

```go
logger = logger.With("service", "user-service")

logger.Info("Пользователь создан", 
    "user", slog.Group("user_data",
        "id", "123",
        "email", "user@example.com",
    ),
)
```

### 5. Настройка через переменные окружения

```go
package main

import (
	"log/slog"
	"os"
	"strconv"
)

func main() {
	// Чтение уровня логирования из переменной окружения
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "info" // Значение по умолчанию
	}

	var level slog.Level
	switch logLevelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)
	logger.Info("Логгер настроен", "level", logLevelStr)
}
```

### 6. Использование контекста

```go
package main

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	// Добавление контекста
	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	
	// Использование логгера с контекстом
	logger.Info("Обработка запроса", 
		"request_id", ctx.Value("request_id"),
		"start_time", time.Now().Format(time.RFC3339),
	)
}
```

## Продвинутые возможности

### 1. Кастомный обработчик

```go
package main

import (
	"log/slog"
	"os"
)

type CustomHandler struct {
	slog.Handler
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	// Добавляем префикс к сообщению
	r.AddAttrs(slog.String("prefix", "[CUSTOM]"))
	return h.Handler.Handle(ctx, r)
}

func main() {
	baseHandler := slog.NewTextHandler(os.Stdout, nil)
	customHandler := &CustomHandler{Handler: baseHandler}
	
	logger := slog.New(customHandler)
	logger.Info("Тестовое сообщение")
}
```

### 2. Форматирование времени

```go
handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true, // Добавление источника лога
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			return slog.String("time", a.Value.Time().Format("2006-01-02 15:04:05"))
		}
		return a
	},
})
```

### 3. Логирование ошибок

```go
err := fmt.Errorf("ошибка выполнения операции")
logger.Error("Опция завершилась с ошибкой", 
    "operation", "create_user",
    "error", err,
)
```

## Best Practices

1. **Используйте осмысленные уровни логирования**:
   - Debug: детальная информация для отладки
   - Info: общая информация о работе приложения
   - Warn: предупреждения о потенциальных проблемах
   - Error: ошибки, которые требуют внимания

2. **Добавляйте контекст**: используйте атрибуты для добавления контекстной информации

3. **Используйте структурированные данные**: вместо форматирования строк используйте атрибуты

4. **Настраивайте уровень логирования**: используйте разные уровни для разных окружений (development, production)

5. **Не логируйте чувствительную информацию**: пароли, токены и другая конфиденциальная информация не должна попадать в логи

## Пример интеграции в приложение

```go
package main

import (
	"log/slog"
	"os"
)

type LoggerConfig struct {
	Env     string
	Service string
	Level   slog.Level
}

func NewLogger(config LoggerConfig) *slog.Logger {
	var handler slog.Handler
	
	if config.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: config.Level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     config.Level,
			AddSource: true,
		})
	}
	
	logger := slog.New(handler)
	logger = logger.With("service", config.Service)
	
	return logger
}

func main() {
	config := LoggerConfig{
		Env:     "development",
		Service: "auth-service",
		Level:   slog.LevelDebug,
	}
	
	logger := NewLogger(config)
	
	logger.Info("Приложение запущено", "version", "1.0.0")
	logger.Debug("Отладочная информация")
}
```

## Заключение

Библиотека `slog` предоставляет мощные и гибкие возможности для логирования в Go. Она заменяет старый пакет `log` и предлагает структурированное логирование с поддержкой различных форматов вывода. Использование `slog` поможет сделать ваше приложение более наблюдаемым (observable) и упростит отладку и мониторинг.