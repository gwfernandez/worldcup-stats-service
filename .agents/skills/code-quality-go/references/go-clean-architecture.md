# Clean Architecture en Go — Patrones y Antipatrones

Referencia de estructuras correctas e incorrectas para auditar código Go con Clean Architecture.

---

## Estructura de capas

```
/internal
├── handler/        ← HTTP, gRPC, CLI (entrada)
│   └── user_handler.go
├── service/        ← Lógica de negocio
│   └── user_service.go
├── repository/     ← Acceso a datos
│   └── user_repository.go
├── domain/         ← Entidades, interfaces, errores de dominio
│   ├── user.go
│   └── errors.go
└── infra/          ← Implementaciones concretas (DB, HTTP clients)
    └── postgres/
        └── user_repo.go
```

---

## Regla de dependencias

```
Handler → Service Interface → Repository Interface
                ↑                      ↑
         ServiceImpl            RepositoryImpl (infra)
```

Las capas internas NO conocen las externas. El dominio no importa nada del proyecto.

---

## ✅ Patrones correctos

### Handler correcto
```go
// Solo orquesta: parsea request, llama al service, formatea response
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")

    user, err := h.userService.GetByID(c.Request.Context(), id)
    if err != nil {
        if errors.Is(err, domain.ErrUserNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
        return
    }

    c.JSON(http.StatusOK, toUserResponse(user))
}
```

### Service correcto
```go
// Solo lógica de negocio, depende de interfaces
type UserService struct {
    repo domain.UserRepository  // ← interfaz, no implementación concreta
    mailer domain.Mailer
}

func (s *UserService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
    if err := input.Validate(); err != nil {
        return nil, fmt.Errorf("input inválido: %w", err)
    }

    existing, err := s.repo.GetByEmail(ctx, input.Email)
    if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
        return nil, fmt.Errorf("UserService.Register: %w", err)
    }
    if existing != nil {
        return nil, domain.ErrEmailAlreadyExists
    }

    user := domain.NewUser(input)
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("UserService.Register guardar usuario: %w", err)
    }

    return user, nil
}
```

### Repository correcto (interfaz en domain)
```go
// domain/user.go — interfaz definida en el dominio
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Save(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

// infra/postgres/user_repo.go — implementación concreta en infra
type postgresUserRepo struct {
    db *sql.DB
}

func (r *postgresUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
    // query SQL aquí
}
```

---

## ❌ Antipatrones frecuentes

### Handler con lógica de negocio
```go
// ❌ MAL — validación de negocio en el handler
func (h *UserHandler) Register(c *gin.Context) {
    var input RegisterInput
    c.BindJSON(&input)

    // Esto es lógica de negocio, no pertenece aquí
    existing, _ := h.db.QueryRow("SELECT id FROM users WHERE email = ?", input.Email)
    if existing != nil {
        c.JSON(400, gin.H{"error": "email ya existe"})
        return
    }
    // ...
}
```

### Service con acceso directo a DB
```go
// ❌ MAL — SQL en el service
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := s.db.QueryRowContext(ctx,
        "SELECT id, name, email FROM users WHERE id = ?", id,
    ).Scan(&user.ID, &user.Name, &user.Email)
    return &user, err
}
```

### Dependencia de implementación concreta en lugar de interfaz
```go
// ❌ MAL — depende de la implementación concreta
type UserService struct {
    repo *PostgresUserRepository  // ← concreto, no testeable
}

// ✅ BIEN
type UserService struct {
    repo UserRepository  // ← interfaz
}
```

### Interfaz definida donde se implementa (no en el dominio)
```go
// ❌ MAL — la interfaz vive junto a su implementación
// infra/postgres/user_repo.go
type UserRepository interface { ... }  // no debería estar aquí

// ✅ BIEN — la interfaz vive en domain/
// domain/user.go
type UserRepository interface { ... }
```

---

## Interfaces: cuándo y cómo

### Regla de Go: interfaces donde se consumen, no donde se implementan
```go
// El service define qué necesita, no importa cómo se implementa
package service

type UserRepository interface {
    GetByID(ctx context.Context, id string) (*domain.User, error)
    Save(ctx context.Context, user *domain.User) error
}
```

### Interfaces pequeñas y enfocadas (Interface Segregation)
```go
// ❌ MAL — interfaz demasiado grande
type UserRepository interface {
    GetByID(...)
    GetByEmail(...)
    Save(...)
    Delete(...)
    List(...)
    Count(...)
    Search(...)
    GetWithOrders(...)
}

// ✅ BIEN — dividir por responsabilidad
type UserReader interface {
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserWriter interface {
    Save(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

---

## Context: propagación correcta

```go
// ✅ Context siempre como primer parámetro y propagado
func (h *Handler) GetUser(c *gin.Context) {
    ctx := c.Request.Context()
    user, err := h.service.GetByID(ctx, id)  // ← ctx propagado
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
    return s.repo.GetByID(ctx, id)  // ← ctx propagado al repo
}

// ❌ MAL — context ignorado o creado de la nada
func (s *Service) GetByID(id string) (*User, error) {
    return s.repo.GetByID(context.Background(), id)  // pierde deadline/cancelación
}
```