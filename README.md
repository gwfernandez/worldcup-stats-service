# World Cups API

API REST que expone datos históricos de los Mundiales de Fútbol desde 1930 hasta la actualidad. Permite consultar información sobre confederaciones, selecciones, jugadores, partidos, goles y estadísticas de todas las ediciones del mundial.

---

## Stack tecnológico

| Componente | Tecnología |
|---|---|
| Lenguaje | Go 1.23 |
| Framework HTTP | Gin |
| Arquitectura | Clean Architecture (handler → service → repository) |
| Base de datos | PostgreSQL (Neon) |
| Driver DB | pgx v5 |
| Queries | sqlc |
| Migraciones | golang-migrate |
| Deploy | Railway |

---

## Requisitos previos

Asegurate de tener instalado lo siguiente antes de correr el proyecto:

- [Go 1.23](https://go.dev/dl/)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- Una base de datos PostgreSQL (se recomienda [Neon](https://neon.tech) free tier)

---

## Variables de entorno

Copiá el archivo de ejemplo y completá los valores:

```bash
cp .env.example .env
```

| Variable | Descripción | Ejemplo |
|---|---|---|
| `DATABASE_URL` | Connection string de PostgreSQL | `postgresql://user:pass@host/worldcups_dev` |
| `PORT` | Puerto donde corre la API | `8080` |
| `GIN_MODE` | Modo de Gin | `debug` / `release` |

---

## Cómo correr el proyecto localmente

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/worldcups-api.git
cd worldcups-api
```

### 2. Instalar dependencias

```bash
go mod tidy
```

### 3. Configurar variables de entorno

```bash
cp .env.example .env
# Editá .env con tu DATABASE_URL de Neon
```

### 4. Correr las migraciones

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

### 5. Iniciar la API

```bash
go run cmd/main.go
```

La API va a estar disponible en `http://localhost:8080`

### 6. Verificar que funciona

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

---

## Cómo correr las migraciones

### Aplicar todas las migraciones

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

### Revertir la última migración

```bash
migrate -path db/migrations -database "$DATABASE_URL" down 1
```

### Revertir todas las migraciones

```bash
migrate -path db/migrations -database "$DATABASE_URL" down
```

---

## Cómo regenerar sqlc

Cada vez que modifiques un archivo en `/db/queries`, regenerá el código Go con:

```bash
sqlc generate
```

El código generado se ubica en `/internal/repository`.

---

## Cómo correr los tests

El proyecto incluye pruebas unitarias para todas las capas (config, repository, service y handler) alcanzando un 100% de cobertura lógica, utilizando `testify` y `pgxmock`.

### Ejecutar todos los tests

```bash
go test ./...
```

### Ver la cobertura (Coverage)

Para generar el reporte de cobertura y ver el porcentaje total:

```bash
go test -coverprofile=.coverage/coverage.out ./...
go tool cover -func=.coverage/coverage.out
```

*Nota: Es habitual ignorar los paquetes `cmd` y `db/sqlc` (auto-generado) en el cálculo final de cobertura de la lógica de negocio.*

---

## Endpoints disponibles

### Health

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/health` | Estado de la API |

### Confederaciones

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/v1/confederations` | Listar todas las confederaciones |
| `GET` | `/api/v1/confederations/:id` | Obtener confederación por id |
| `POST` | `/api/v1/confederations` | Crear nueva confederación |
| `PUT` | `/api/v1/confederations/:id` | Actualizar confederación |
| `DELETE` | `/api/v1/confederations/:id` | Eliminar confederación |

### Ejemplos de request

**Crear confederación**
```bash
curl -X POST http://localhost:8080/api/v1/confederations \
  -H "Content-Type: application/json" \
  -d '{"code": "UEFA", "name": "Union of European Football Associations"}'
```

**Listar confederaciones**
```bash
curl http://localhost:8080/api/v1/confederations
```

---

## Estructura del proyecto

```
/cmd
  main.go                       # entrypoint de la aplicación
/internal
  /domain                       # structs de las entidades del dominio
  /handler                      # controllers HTTP, validación de input
  /service                      # lógica de negocio e interfaces
  /repository                   # acceso a datos e interfaces
/db
  /migrations                   # archivos SQL versionados de migraciones
  /queries                      # archivos SQL para generación con sqlc
/config
  config.go                     # lectura de variables de entorno
/.ai
  /prompts                      # prompts para el IDE
  /instructions                 # documentación técnica interna
  context.md                    # contexto general del proyecto para el IDE
sqlc.yaml                       # configuración de sqlc
.env.example                    # variables de entorno requeridas
go.mod
```

---

## Formato de errores

Todos los errores retornan JSON con el siguiente formato:

```json
{
  "error": "descripción del error"
}
```
