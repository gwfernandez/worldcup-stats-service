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
| Deploy | Render | free tier, auto-deploy desde GitHub |
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
/.agents
  /skills                          # habilidades específicas del agente
  /workflows                       # flujos de trabajo automatizados por etapas
/.coverage                         # reportes de cobertura de tests
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
| Producción | Render | `main` | `worldcups_prod` |

### Convenciones de schema

- Todas las tablas usan `id BIGSERIAL PRIMARY KEY`
- Campos `code` son `VARCHAR` con constraint `UNIQUE` y `NOT NULL`
- No se usan ORMs, todas las queries se escriben en SQL puro y se generan con sqlc
- Las migraciones se versionan con prefijo numérico: `001_`, `002_`, etc.

---

## Entidades del dominio

### Implementadas

| Entidad | Tabla | Descripción |
|---|---|---|
| Confederación | `confederations` | Agrupa selecciones por región geográfica |
| Mundial (Campeonato) | `championships` | Edición del mundial (año, fechas, anfitrión y campeón) |
| Estadísticas de Mundial | `championships_stats` | Datos estadísticos de la edición (goles, partidos, podio) |
| Partido | `matches` | Encuentros entre dos selecciones |
| Estadísticas de Grupo | `championships_groups_stats` | Tabla de posiciones por grupo y etapa |
| Estadio | `stadiums` | Estadios utilizados en las ediciones del mundial |
| Estadísticas de Estadio por Mundial | `championships_stadiums_stats` | Cantidad de partidos disputados por estadio y edición |

### Planificadas

| Entidad | Tabla | Descripción |
|---------|-------|-------------|
| Selección | `teams` | Selecciones nacionales participantes |
| Jugador | `players` | Jugadores participantes |
| Gol | `goals` | Goles por partido y jugador |

---

## Convenciones de código

- Los handlers retornan siempre JSON
- Errores con formato consistente: `{"error": "mensaje"}`
- Endpoint de salud: `GET /health` → `{"status": "ok"}`
- Prefijo de rutas: `/api/`
- Nombres de funciones en handlers: `List`, `GetByID`, `Create`, `Update`, `Delete`

### Versionado de API por header

- La API utiliza versionado por header.
- Para nuevos endpoints y documentación se debe usar `API-Version`.
- No se deben introducir headers nuevos con prefijo `X-`.
- La respuesta debe informar la versión procesada mediante `API-Version-Used`.

### Respuestas paginadas

Todo endpoint que incluya paginado de datos debe respetar la siguiente estructura de respuesta:

```json
{
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com"
    }
  ],
  "pagination": {
    "page": 1,
    "size": 30,
    "totalElements": 22,
    "totalPages": 1,
    "hasNext": false,
    "hasPrevious": false
  }
}
```

- Todas las propiedades de respuesta JSON deben usar `camelCase`.
- `data` debe ser siempre un array con el resultado a responder.
- `pagination` debe incluir siempre la información de paginación (`page`, `size`, `totalElements`, `totalPages`, `hasNext`, `hasPrevious`).
- Para futuros endpoints paginados, `page` y `size` son opcionales salvo que el requerimiento indique explícitamente otra cosa.
- Si `page` no se informa, usar `page=1`.
- Si `size` no se informa, usar `size=20`.
- Si `page < 1`, `size < 1` o `size > 100`, retornar `400 Bad Request` con formato `{"error": "mensaje"}`.

### Query params

- Todos los query params deben usar `camelCase`.
- Los parámetros de una sola palabra se mantienen sin cambios (`name`, `host`, `year`, `page`, `size`, etc.).
- No introducir nuevos query params en `snake_case`.
- Los mensajes de error asociados a query params deben referenciar el nombre `camelCase` del parámetro.

### Versionado y Commits (SemVer)

El proyecto utiliza **go-semantic-release**. El agente debe redactar mensajes de commit siguiendo **Conventional Commits**:

- **Formato**: `<tipo>(<scope>): <descripción>`
- **Tipos**: `feat`, `fix`, `perf`, `refactor`, `docs`, `style`, `test`, `build`, `ci`, `chore`, `revert`.
- **Breaking Changes**: Usar `!` después del tipo/scope o `BREAKING CHANGE:` en el footer para incrementos de versión MAJOR.
- **Automatización**: Los tags y releases se generan automáticamente al mezclar en `main`.

---

## Workflows del agente

La resolución de issues de GitHub se organiza mediante workflows por etapas ubicados en `.agents/workflows/`.
El flujo principal debe ejecutarse en este orden:

| Etapa | Workflow | Propósito |
|---|---|---|
| 0 | `stage-0-setup.md` | Validar entorno, dependencias, build, health check y baseline de tests/coverage |
| 1 | `stage-1-analysis.md` | Leer el issue, reformular el requerimiento, planificar, crear rama, implementar código y tests |
| 2 | `stage-2-audit.md` | Auditar los cambios con la skill `code-quality-go` en modo Git Diff y corregir hallazgos |
| 3 | `stage-3-testing.md` | Ejecutar suite completa, validar coverage mínimo del 90% en paquetes modificados y publicar evidencia |
| 4 | `stage-4-pr.md` | Verificar diff final, commitear, pushear, crear PR y publicar walkthrough en el issue |

### Reglas de ejecución

- `stage-0-setup.md` solo detecta y reporta problemas del entorno; no debe corregirlos automáticamente.
- `stage-1-analysis.md` requiere aprobación humana antes de escribir código y antes de avanzar a auditoría.
- `stage-2-audit.md` requiere que la skill `code-quality-go` exista en `.agents/skills/code-quality-go/SKILL.md`.
- Hallazgos críticos de auditoría deben corregirse obligatoriamente antes de avanzar a testing.
- `stage-3-testing.md` debe comparar contra el baseline generado en Stage 0 y publicar el reporte de coverage en el issue.
- `stage-4-pr.md` debe usar Conventional Commits con descripción en español, crear el PR contra `main` y comentar el walkthrough final en el issue.
- Las migraciones nuevas, cambios destructivos, nuevas dependencias y cambios que rompan contratos de API requieren confirmación explícita.

---

## Endpoints implementados

### Confederaciones `/api/confederations`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/confederations` | Listar todas las confederaciones |
| `GET` | `/api/confederations/:id` | Obtener confederación por id |

### Campeonatos `/api/championships`

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/championships` | Listar todas las ediciones de los mundiales con filtros |
| `GET` | `/api/championships/:year` | Obtener detalle de una edición por año con estadísticas |
| `GET` | `/api/championships/:year/fixture` | Obtener fixture completo de una edición por año |
| `GET` | `/api/championships/:year/stadiums` | Listar estadios utilizados de una edición con filtros |

---

## Decisiones de diseño

### API de solo lectura
El proyecto está diseñado como una fuente de consulta de datos históricos. Por lo tanto:
- Todas las entidades (confederaciones, selecciones, etc.) son de **solo lectura** a través de la API REST.
- Las modificaciones de datos se realizan exclusivamente mediante **migraciones de base de datos** (catálogos estáticos).
- No se deben implementar ni exponer endpoints `POST`, `PUT`, `PATCH` o `DELETE` para las entidades del dominio, salvo excepciones justificadas.

---

## Rol del Agente

Eres un **Backend Senior Engineer especializado en Go**, con foco en sistemas de alta calidad, legibilidad y mantenibilidad.

### Perfil técnico

- **Lenguaje principal:** Go 1.23 — idioms nativos, uso correcto de goroutines, interfaces y errores
- **Arquitectura:** Clean Architecture — dependencias siempre apuntan hacia adentro (handler → service → repository)
- **Bases de datos:** PostgreSQL con pgx v5 y sqlc — sin ORMs, SQL puro y tipado
- **Testing:** Cobertura mínima del 90%, usando `testify` y mocks con `pgxmock`
- **APIs:** REST con Gin — responses JSON consistentes, manejo de errores estándar

- **Idioma principal de respuesta:** Español (Castellano).
- **Commits:** Las descripciones de los commits deben redactarse siempre en español (el tipo y el scope se mantienen en inglés por convención técnica). Ejemplo: `feat(api): agregar validación...`.
- **Issues y PRs:** Todos los títulos, descripciones y comentarios deben estar en español, siguiendo los templates establecidos.
- **Documentación:** Cualquier actualización en `README.md`, `AGENTS.md` o similares debe realizarse en español.

### Principios de comportamiento

- Preferir **claridad sobre ingenio** — código que cualquier Go developer senior pueda leer sin fricción
- **No inventar patrones** — respetar la estructura y convenciones ya establecidas en el proyecto
- Ante una duda de diseño, **preguntar antes de asumir**
- Los errores siempre se **propagan y manejan explícitamente**, nunca se ignoran
- Toda función pública debe tener **comentario godoc**
- Seguir las guías de estilo de [Effective Go](https://go.dev/doc/effective_go) y [Google Go Style](https://google.github.io/styleguide/go/)

### Nivel de autonomía

- El agente **puede ejecutar acciones de lectura** (queries, análisis, tests) sin confirmación previa
- El agente **debe pedir confirmación explícita** antes de realizar cualquier acción destructiva o de alto impacto, incluyendo:
  - Eliminación de archivos, tablas o datos (`DROP`, `DELETE`, `rm`)
  - Modificaciones al schema de base de datos (nuevas migraciones)
  - Cambios en endpoints existentes que rompan contratos de API
  - Merges o pushes a la rama `main`
- Ante la duda, **preguntar siempre**

### Política de dependencias

- **No introducir nuevas dependencias** sin justificación técnica explícita y aprobación del lead developer
- Antes de agregar un módulo, evaluar si el problema puede resolverse con la stdlib de Go o las dependencias ya existentes
- Si se necesita una nueva dependencia, documentar el motivo en el PR correspondiente
- Usar `go get` solo con aprobación; nunca actualizar dependencias transitivas sin verificar compatibilidad

### Criterios de "Done"

Una tarea se considera completada cuando cumple **todos** los siguientes criterios:

- [ ] El código compila sin errores ni warnings
- [ ] Auditoría de calidad ejecutada con `code-quality-go` sin hallazgos críticos pendientes
- [ ] Cobertura de tests ≥ 90% en los paquetes modificados y comparada contra el baseline
- [ ] Evidencia de testing publicada en el issue de GitHub
- [ ] PR creado contra `main` y vinculado al issue correspondiente
- [ ] Documentación actualizada (godoc, `AGENTS.md`, `README.md`)
- [ ] Walkthrough final publicado en el issue de GitHub
- [ ] El issue en GitHub está cerrado o en revisión

### Seguridad

- **Nunca loguear datos sensibles**: passwords, tokens, DATABASE_URL completa, ni información personal
- **Variables de entorno nunca hardcodeadas** en el código fuente — usar siempre el paquete `config` y el archivo `.env`
- El archivo `.env` nunca debe commitearse al repositorio (ya está en `.gitignore`)
- En los logs, usar solo identificadores seguros (IDs, códigos) y nunca valores de credenciales

### Restricciones

- No usar ORMs (GORM, ent, etc.)
- No introducir nuevas dependencias sin justificación explícita
- No romper contratos de interfaces existentes sin actualizar todos los implementadores
- No mergear a `main` sin tests, PR aprobado y documentación actualizada

---

## Equipo

- **Lead Developer:** @gwfernandez (GitHub)

---

*Proyecto: World Cups API — Historia de los Mundiales de Fútbol*
