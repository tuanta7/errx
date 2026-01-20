# errx

⛔ An ergonomic Go error package for layered architectures

```bash
go get github.com/tuanta7/errx
```

## Features

- Immutable error wrapping with builder-style methods
- HTTP and gRPC status code mapping
- User-defined internal codes and localized messages
- Load messages from JSON files
- Designed for layered architecture

### TODO

- [ ] Add support for loading configuration from YAML files
- [ ] Add support for a hierarchical error format within the file
- [ ] Add more common errors

## Architecture

`errx` is designed around a typical layered architecture pattern:

```
┌─────────────────────────────────────────────────────────────┐
│  Transport Layer                                            │
│  - Maps internal code to HTTP/gRPC status code              │
│  - Returns localized message to client                      │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │ *errx.Error with internal code
                            │
┌─────────────────────────────────────────────────────────────┐
│  Business Layer                                             │
│  - Receives predefined error from repo                      │
│  - Attaches internal code for business context              │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │ predefined *errx.Error
                            │
┌─────────────────────────────────────────────────────────────┐
│  Adapter Layer                                              │
│  - Returns predefined errors (ErrRecordNotFound, etc.)      │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

### 1. Define Internal Error Codes

Define your project-specific error codes as constants:

```go
// domain/errors.go
package domain

const (
    ErrCounterNotFound = "COUNTER_NOT_FOUND"
    ErrUserNotFound    = "USER_NOT_FOUND"
    ErrEmailTaken      = "EMAIL_TAKEN"
)
```

### 2. Register Codes and Messages (on startup)

Register status codes and localized messages for your internal codes:

```go
package main

import (
    "net/http"

    "github.com/tuanta7/errx"
    "golang.org/x/text/language"
)

func main() {
    // Register HTTP status code for your internal codes
    errx.Global.RegisterHTTPErrorCode(ErrCounterNotFound, http.StatusNotFound)
    errx.Global.RegisterHTTPErrorCode(ErrUserNotFound, http.StatusNotFound)
    errx.Global.RegisterHTTPErrorCode(ErrEmailTaken, http.StatusConflict)

    // Register localized messages
    errx.Global.RegisterMessage(ErrCounterNotFound, language.English.String(), "Counter not found")
    errx.Global.RegisterMessage(ErrCounterNotFound, language.Vietnamese.String(), "Không tìm thấy bộ đếm")

    // ... start your server
}
```

### 3. Adapter Layer (Repository)

Return predefined errors from your data layer:

```go
package repository

import (
    "encoding/json"

    "github.com/tuanta7/errx/errors"
)

func (r *Repository) GetCounter(key string) (*Counter, error) {
    value, err := r.cache.Get(key)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, errors.ErrRecordNotFound // Return predefined error
        }
        return nil, err
    }

    counter := &Counter{}
    err = json.Unmarshal(value, counter)
    if err != nil {
        return nil, errors.New("failed to unmarshal counter", err)
    }

    return counter, nil
}
```

### 4. Business Layer (Usecase)

Attach your internal code to give business context:

```go
package usecase

import (
    "github.com/tuanta7/errx/errors"
)

func (uc *UseCase) GetCounter(key string) (*Counter, error) {
    counter, err := uc.repo.GetCounter(key)
    if err != nil {
        if errors.Is(err, errors.ErrRecordNotFound) {
            // Attach business-specific code
            return nil, errors.ErrRecordNotFound.WithCode(ErrCounterNotFound)
        }
        return nil, err
    }

    return counter, nil
}
```

### 5. Transport Layer (Handler)

Map to HTTP/gRPC response using registered status codes and messages:

```go
package handler

import (
    "encoding/json"
    "net/http"

    "github.com/tuanta7/errx"
    lang "golang.org/x/text/language"
)

func (h *Handler) GetCounter(w http.ResponseWriter, r *http.Request) {
    counterName := r.PathValue("id")
    language := r.URL.Query().Get("language")
    if language == "" {
        language = lang.English.String()
    }

    counter, err := h.uc.GetCounter(counterName)
    if err != nil {
        // Get HTTP status code and localized message
        httpCode, message := errx.Global.GetHTTPResponse(err, language)
        // httpCode = 404, message = "Counter not found" (or "Không tìm thấy bộ đếm" for Vietnamese)
        http.Error(w, message, httpCode)
        return
    }

    jsonCounter, _ := json.Marshal(counter)
    _, _ = w.Write(jsonCounter)
}
```

## Loading Messages from Files

Instead of registering messages one by one, you can load them from JSON or YAML files. Each file contains messages for a single language.

### JSON Format

Create one file per language with a flat code-to-message mapping:

**static/en.json**
```json
{
  "ERR_RESOURCE_NOT_FOUND": "Resource not found",
  "ERR_INVALID_INPUT": "Invalid input provided",
  "ERR_UNAUTHORIZED_ACCESS": "Unauthorized access"
}
```

**static/vi.json**
```json
{
  "ERR_RESOURCE_NOT_FOUND": "Không tìm thấy tài nguyên",
  "ERR_INVALID_INPUT": "Dữ liệu đầu vào không hợp lệ",
  "ERR_UNAUTHORIZED_ACCESS": "Truy cập không được phép"
}
```

**static/es.json**
```json
{
  "ERR_RESOURCE_NOT_FOUND": "Recurso no encontrado",
  "ERR_INVALID_INPUT": "Entrada no válida proporcionada",
  "ERR_UNAUTHORIZED_ACCESS": "Acceso no autorizado"
}
```

### Load Messages

```go
import (
    "github.com/tuanta7/errx"
    "github.com/tuanta7/errx/errors"
    "github.com/tuanta7/errx/parsers/json"
    lang "golang.org/x/text/language"
)

func main() {
    errx.SetGlobal(errx.New())

    // Load messages for each language
    err := errx.Global.LoadMessages(lang.English.String(), "./static/en.json", json.Parser())
    if err != nil {
        panic(err)
    }

    err = errx.Global.LoadMessages(lang.Spanish.String(), "./static/es.json", json.Parser())
    if err != nil {
        panic(err)
    }

    err = errx.Global.LoadMessages(lang.Vietnamese.String(), "./static/vi.json", json.Parser())
    if err != nil {
        panic(err)
    }

    // Now you can use the messages
    ErrNotFound := errors.New("default not found message").WithCode("ERR_RESOURCE_NOT_FOUND")

    fmt.Println(errx.Global.GetMessage(ErrNotFound, lang.English.String()))    // "Resource not found"
    fmt.Println(errx.Global.GetMessage(ErrNotFound, lang.Spanish.String()))    // "Recurso no encontrado"
    fmt.Println(errx.Global.GetMessage(ErrNotFound, lang.Vietnamese.String())) // "Không tìm thấy tài nguyên"
}
```

## Predefined Errors

These errors are designed to be returned from the adapter/repository layer:

| Error                          | HTTP | gRPC             | Use Case                        |
| ------------------------------ | ---- | ---------------- | ------------------------------- |
| `ErrInternal`                  | 500  | Internal         | Unexpected server errors        |
| `ErrServiceUnavailable`        | 503  | Unavailable      | Service temporarily unavailable |
| `ErrInvalidParameter`          | 400  | InvalidArgument  | Validation failures             |
| `ErrRecordNotFound`            | 404  | NotFound         | Entity not found in database    |
| `ErrConnectionTimeout`         | 503  | Unavailable      | Database/service connection     |
| `ErrOperationTimeout`          | 504  | DeadlineExceeded | Query/operation timeout         |
| `ErrForeignKeyViolation`       | 400  | InvalidArgument  | Invalid reference               |
| `ErrUniqueConstraintViolation` | 400  | AlreadyExists    | Duplicate entry                 |

## Thread Safety

All methods return new `*Error` instances, leaving the original unchanged. Safe for concurrent use.

## License

MIT
