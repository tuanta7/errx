# errx

⛔ An ergonomic Go error package for layered architectures

## Features

- Immutable error wrapping with builder-style methods
- HTTP and gRPC status code mapping
- User-defined internal codes and localized messages
- Designed for clean architecture (repo → usecase → handler)

## Installation

```bash
go get github.com/tuanta7/errx
```

## Architecture

`errx` is designed around a layered architecture pattern:

```
┌─────────────────────────────────────────────────────────────┐
│  Handler Layer                                              │
│  - Maps error to HTTP/gRPC status code                      │
│  - Returns localized message to client                      │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │ *errx.Error with code
                            │
┌─────────────────────────────────────────────────────────────┐
│  Usecase Layer                                              │
│  - Receives predefined error from repo                      │
│  - Attaches internal code for business context              │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │ predefined *errx.Error
                            │
┌─────────────────────────────────────────────────────────────┐
│  Repository Layer                                           │
│  - Returns predefined errors (ErrRecordNotFound, etc.)      │
└─────────────────────────────────────────────────────────────┘
```

## Usage

### 1. Register Codes and Messages (on startup)

Define your internal codes with status codes and localized messages:

```go
func init() {
    // Register status codes for your internal codes
    errx.RegisterHTTPErrorCode("USER_NOT_FOUND", http.StatusNotFound)
    errx.RegisterHTTPErrorCode("EMAIL_TAKEN", http.StatusConflict)
    
    // Or with both HTTP and gRPC
    errx.RegisterErrorCode("PAYMENT_FAILED", errx.StatusCode{
        HTTPCode: http.StatusPaymentRequired,
        GRPCCode: uint32(codes.FailedPrecondition),
    })

    // Register localized messages for your internal codes
    errx.RegisterMessage("USER_NOT_FOUND", "en", "User not found")
    errx.RegisterMessage("USER_NOT_FOUND", "vi", "Không tìm thấy người dùng")
    
    errx.RegisterMessage("EMAIL_TAKEN", "en", "Email already registered")
    errx.RegisterMessage("EMAIL_TAKEN", "vi", "Email đã được đăng ký")
}
```

### 2. Repository Layer

Return predefined errors:

```go
func (r *UserRepo) FindByID(id int) (*User, error) {
    user, err := r.db.Query("SELECT ...")
    if errors.Is(err, sql.ErrNoRows) {
        return nil, errx.ErrRecordNotFound // predefined error
    }
    return user, err
}
```

### 3. Usecase Layer

Attach your internal code:

```go
func (u *UserUsecase) GetUser(id int) (*User, error) {
    user, err := u.repo.FindByID(id)
    if err != nil {
        // Attach business-specific code
        return nil, err.(*errx.Error).WithCode("USER_NOT_FOUND")
    }
    return user, nil
}
```

### 4. Handler Layer

Map to HTTP/gRPC response:

```go
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.usecase.GetUser(id)
    if err != nil {
        lang := r.Header.Get("Accept-Language")
        statusCode, message := errx.HTTPResponse(err.(*errx.Error), lang)
        // statusCode = 404, message = "User not found"
        http.Error(w, message, statusCode)
        return
    }
    json.NewEncoder(w).Encode(user)
}
```

### Predefined Errors

| Error                          | HTTP | gRPC             |
| ------------------------------ | ---- | ---------------- |
| `ErrInternal`                  | 500  | Internal         |
| `ErrInvalidParameter`          | 400  | InvalidArgument  |
| `ErrRecordNotFound`            | 404  | NotFound         |
| `ErrConnectionTimeout`         | 503  | Unavailable      |
| `ErrOperationTimeout`          | 504  | DeadlineExceeded |
| `ErrForeignKeyViolation`       | 400  | InvalidArgument  |
| `ErrUniqueConstraintViolation` | 400  | AlreadyExists    |


## Thread Safety

All methods return new `*Error` instances, leaving the original unchanged. Safe for concurrent use.
