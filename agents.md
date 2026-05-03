# Contexto del Proyecto — World Cups API

## Descripción general

API REST que expone datos históricos de los Mundiales de Fútbol desde 1930 hasta la actualidad.
El proyecto está orientado a ser una fuente de consulta sobre confederaciones, selecciones, jugadores, partidos, goles y estadísticas de todos los mundiales.

---

## Stack tecnológico

| Componente | Tecnología | Detalle |
|---|---|---|
| Lenguaje | Go 1.23 | versión estable más reciente |
| Framework HTTP | Gin | routing, validación, middleware |
| Arquitectura | Clean Architecture | handler → service → repository |
| Base de datos | PostgreSQL | Neon free tier, 3 GiB |
| Driver DB | pgx v5 | driver nativo para PostgreSQL |
| Queries | sqlc | generación de código Go tipado desde SQL |
| Migraciones | golang-migrate | control de versiones del schema |
| Deploy | Railway | free tier, auto-deploy desde GitHub |
| Repositorio | GitHub | rama `main` → PROD / rama `dev` → DEV |

---

## Arquitectura

El proyecto sigue **Clean Architecture** con tres capas bien definidas:

```
handler      →  recibe el request HTTP, valida input, retorna JSON
  ↓
service      →  contiene la lógica de negocio, no conoce HTTP ni SQL
  ↓
repository   →  acceso a datos, usa código generado por sqlc con pgx
```

Cada capa depende de interfaces, no de implementaciones concretas.
Esto permite testear el service de forma aislada sin levantar HTTP ni base de datos.

---

## Estructura de carpetas

```
/cmd
  main.go                          # entrypoint de la aplicación
/internal
  /domain                          # structs de entidades
  /handler                         # controllers HTTP (Gin)
  /service                         # lógica de negocio + interfaces
  /repository                      # acceso a datos + interfaces
/db
  /migrations                      # archivos SQL de migraciones
  /queries                         # archivos SQL para sqlc
/config
  config.go                        # configuración desde variables de entorno
/.ai
  /prompts                         # prompts para el IDE
  /instructions                    # documentación técnica del proyecto
  context.md                       # este archivo
sqlc.yaml                          # configuración de sqlc
.env.example                       # variables de entorno requeridas
go.mod
```

---

## Variables de entorno

```bash
DATABASE_URL=postgresql://user:password@host/worldcups_dev
PORT=8080
GIN_MODE=debug   # usar "release" en producción
```

---

## Base de datos

### Estrategia de entornos

| Entorno | Plataforma | Branch Neon | Base de datos |
|---------|-----------|-------------|---------------|
| Desarrollo | local | `dev` | `worldcups_dev` |
| Producción | Railway | `main` | `worldcups_prod` |

### Convenciones de schema

- Todas las tablas usan `id BIGSERIAL PRIMARY KEY`
- Campos `code` son `VARCHAR` con constraint `UNIQUE` y `NOT NULL`
- No se usan ORMs, todas las queries se escriben en SQL puro y se generan con sqlc
- Las migraciones se versionan con prefijo numérico: `001_`, `002_`, etc.

---

## Entidades del dominio

### Implementadas

| Entidad | Tabla | Descripción |
|---------|-------|-------------|
| Confederación | `confederations` | Agrupa selecciones por región geográfica |

### Planificadas

| Entidad | Tabla | Descripción |
|---------|-------|-------------|
| Selección | `teams` | Selecciones nacionales participantes |
| Mundial | `world_cups` | Ediciones del mundial (año, sede, campeón) |
| Fase | `stages` | Grupos, octavos, cuartos, semifinal, final |
| Partido | `matches` | Encuentros entre dos selecciones |
| Jugador | `players` | Jugadores participantes |
| Gol | `goals` | Goles por partido y jugador |

---

## Convenciones de código

- Los handlers retornan siempre JSON
- Errores con formato consistente: `{"error": "mensaje"}`
- Endpoint de salud: `GET /health` → `{"status": "ok"}`
- Prefijo de rutas: `/api/v1/`
- Nombres de funciones en handlers: `List`, `GetByID`, `Create`, `Update`, `Delete`

---

## Endpoints implementados

### Confederaciones `/api/v1/confederations`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/v1/confederations` | Listar todas las confederaciones |
| `GET` | `/api/v1/confederations/:id` | Obtener confederación por id |
| `POST` | `/api/v1/confederations` | Crear nueva confederación |
| `PUT` | `/api/v1/confederations/:id` | Actualizar confederación |
| `DELETE` | `/api/v1/confederations/:id` | Eliminar confederación |

---

## Equipo

- **Lead Developer:** @gwfernandez (GitHub)

---

*Proyecto: World Cups API — Historia de los Mundiales de Fútbol*
