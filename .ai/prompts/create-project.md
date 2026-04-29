# Prompt — Crear Proyecto World Cups API

Crear un proyecto de API REST en Go 1.23 utilizando el siguiente stack y especificaciones.

## Stack

- **Lenguaje:** Go 1.23
- **Framework:** Gin
- **Arquitectura:** Clean Architecture con tres capas: handler → service → repository
- **Base de datos:** PostgreSQL (conexión via Neon free tier)
- **Driver DB:** pgx
- **Generación de queries:** sqlc
- **Deploy:** Railway (free tier)

---

## Estructura del proyecto

```
/cmd
  main.go
/internal
  /domain
    confederation.go              # struct de la entidad
  /handler
    confederation_handler.go
  /service
    confederation_service.go
    confederation_service_interface.go
  /repository
    confederation_repository.go
    confederation_repository_interface.go
/db
  /migrations
    001_create_confederations.sql
  /queries
    confederation.sql
/config
  config.go
/.ai
  /prompts
    create-project.md
  /instructions
    stack.md
sqlc.yaml
.env.example
go.mod
```

---

## Primer caso de uso — ABM de confederaciones de fútbol

### Tabla en base de datos

```sql
CREATE TABLE confederations (
    id      BIGSERIAL    PRIMARY KEY,
    code    VARCHAR(20)  NOT NULL UNIQUE,
    name    VARCHAR(100) NOT NULL
);
```

### Datos iniciales

```sql
INSERT INTO confederations (code, name) VALUES
    ('UEFA',     'Union of European Football Associations'),
    ('CONMEBOL', 'South American Football Confederation'),
    ('CONCACAF', 'Confederation of North, Central America and Caribbean Association Football'),
    ('CAF',      'Confederation of African Football'),
    ('AFC',      'Asian Football Confederation'),
    ('OFC',      'Oceania Football Confederation');
```

### Endpoints REST

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/v1/confederations` | Listar todas |
| `GET` | `/api/v1/confederations/:id` | Obtener por id |
| `POST` | `/api/v1/confederations` | Crear |
| `PUT` | `/api/v1/confederations/:id` | Actualizar |
| `DELETE` | `/api/v1/confederations/:id` | Eliminar |

---

## Requerimientos por capa

### Handler
- Recibir el request HTTP
- Validar el input usando Gin binding (`binding:"required"`)
- Llamar al service correspondiente
- Retornar respuesta JSON con el HTTP status code apropiado

### Service
- Contener la lógica de negocio
- Depender de la interfaz del repository, no de la implementación concreta
- No conocer nada de HTTP ni de SQL

### Repository
- Todas las operaciones de base de datos usando código generado por sqlc con driver pgx

---

## Requerimientos adicionales

- La cadena de conexión a la base de datos debe leerse desde la variable de entorno `DATABASE_URL`
- Incluir un endpoint `/health` que retorne `{"status": "ok"}`
- Retornar respuestas de error JSON consistentes: `{"error": "mensaje"}`
- Generar `sqlc.yaml` configurado para pgx/v5 apuntando a `/db/queries`
- Incluir `go.mod` con todas las dependencias necesarias
- Incluir `.env.example` con las variables de entorno requeridas
- No utilizar ningún ORM

---

## Variables de entorno requeridas

```bash
# .env.example
DATABASE_URL=postgresql://user:password@host/worldcups_dev
PORT=8080
GIN_MODE=debug
```

---

*Proyecto: World Cups API — Historia de los Mundiales de Fútbol*
