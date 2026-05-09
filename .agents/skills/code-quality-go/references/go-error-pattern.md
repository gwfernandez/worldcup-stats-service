# Manejo de Errores en Go — Patrones de Referencia

---

## Regla fundamental

> Cada error debe tener contexto suficiente para entender **dónde** ocurrió y **por qué**.

---

## Wrapping de errores

```go
// ✅ BIEN — agregar contexto con %w (permite errors.Is / errors.As)
if err != nil {
    return fmt.Errorf("UserService.GetByID id=%s: %w", id, err)
}

// ❌ MAL — propagar sin contexto
if err != nil {
    return err
}

// ❌ MAL — perder la cadena de error (no usa %w)
if err != nil {
    return fmt.Errorf("error al obtener usuario: %s", err)
}
```

---

## Errores de dominio (Sentinel Errors)

Definir en `domain/errors.go`:

```go
package domain

import "errors"

// Errores de negocio — el caller puede distinguirlos con errors.Is
var (
    ErrUserNotFound      = errors.New("usuario no encontrado")
    ErrEmailAlreadyExists = errors.New("el email ya está registrado")
    ErrInvalidCredentials = errors.New("credenciales inválidas")
    ErrInsufficientFunds  = errors.New("saldo insuficiente")
)
```

### Usar en el service
```go
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrUserNotFound  // traducir error de infra a dominio
        }
        return nil, fmt.Errorf("UserService.GetByID: %w", err)
    }
    return user, nil
}
```

### Usar en el handler
```go
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrUserNotFound):
            c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
        case errors.Is(err, domain.ErrInvalidCredentials):
            c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
        default:
            // Error inesperado — no exponer detalles internos
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
        }
        return
    }
    c.JSON(http.StatusOK, user)
}
```

---

## Errores con datos estructurados

Para errores que necesitan cargar información adicional:

```go
// domain/errors.go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validación fallida en campo '%s': %s", e.Field, e.Message)
}

// Uso
func (s *UserService) Register(ctx context.Context, input RegisterInput) error {
    if input.Email == "" {
        return &ValidationError{Field: "email", Message: "no puede estar vacío"}
    }
}

// En el handler
var valErr *domain.ValidationError
if errors.As(err, &valErr) {
    c.JSON(http.StatusBadRequest, gin.H{"error": valErr.Error(), "field": valErr.Field})
    return
}
```

---

## Antipatrones de errores

```go
// ❌ Error ignorado con _
result, _ := repo.GetByID(id)

// ❌ Panic donde debería haber return de error
func GetUser(id string) *User {
    user, err := repo.GetByID(id)
    if err != nil {
        panic(err)  // solo para errores de programación, no de runtime
    }
    return user
}

// ❌ log.Fatal en producción (termina el proceso)
if err != nil {
    log.Fatal(err)
}

// ❌ Errores de infraestructura expuestos al cliente
c.JSON(500, gin.H{"error": err.Error()})  // puede exponer detalles de DB/SQL

// ❌ Error wrapping que pierde la cadena
return fmt.Errorf("algo falló: %s", err)  // usar %w para mantener la cadena
```

---

## Errores en goroutines

```go
// ✅ Canal de errores para goroutines
func processItems(items []Item) error {
    errCh := make(chan error, len(items))

    for _, item := range items {
        go func(i Item) {
            if err := process(i); err != nil {
                errCh <- fmt.Errorf("procesar item %s: %w", i.ID, err)
                return
            }
            errCh <- nil
        }(item)
    }

    for range items {
        if err := <-errCh; err != nil {
            return err
        }
    }
    return nil
}

// ✅ errgroup para goroutines con contexto
import "golang.org/x/sync/errgroup"

func processItems(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)

    for _, item := range items {
        item := item
        g.Go(func() error {
            return process(ctx, item)
        })
    }

    return g.Wait()
}
```

---

## Formato de respuesta HTTP — Estándar

Siempre usar este formato en handlers:

```go
// Error
c.JSON(statusCode, gin.H{"error": "mensaje legible para el cliente"})

// Éxito con datos
c.JSON(http.StatusOK, gin.H{"data": result})

// Éxito sin contenido
c.Status(http.StatusNoContent)

// Creación exitosa
c.JSON(http.StatusCreated, gin.H{"data": newResource, "id": newResource.ID})
```