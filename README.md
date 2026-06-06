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
| Deploy | Render |

---

## Versionado y Convenciones

Este proyecto utiliza **[go-semantic-release](https://github.com/go-semantic-release/semantic-release)** para la gestión automática de versiones y changelogs.

### Convención de Commits
Para que el versionado automático funcione, todos los commits en `main` deben seguir la especificación de **[Conventional Commits](https://www.conventionalcommits.org/)**:

- `feat:` Nuevas funcionalidades (genera versión MINOR)
- `fix:` Corrección de errores (genera versión PATCH)
- `perf:` Mejoras de rendimiento (genera versión PATCH)
- `docs:`, `style:`, `refactor:`, `test:`, `chore:`, `ci:` (no generan release por defecto)

Cualquier commit con `!` después del tipo (ej: `feat!:`) o con el texto `BREAKING CHANGE:` en el footer generará una versión **MAJOR**.

---

## Requisitos previos

Asegurate de tener instalado lo siguiente antes de correr el proyecto:

- [Go 1.23](https://go.dev/dl/)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- Base de datos PostgreSQL ([Neon](https://neon.tech) free tier)

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
git clone https://github.com/gwfernandez/worldcup-stats-service
cd worldcup-stats-service
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

### 4. Correr las migraciones (si la base de datos esta vacía)

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

### Versionado de API

La API utiliza versionado por Header en lugar de versionado en la URL. 
Para acceder a una versión específica, se debe incluir el header HTTP `API-Version`.
Si no se envía el header, la API utilizará la versión `1` por defecto.

**Ejemplo:** `API-Version: 1`

Para más detalles, consultar [Estrategia de Versionado](docs/API_VERSIONING.md).

### Convención de respuestas paginadas

Todo endpoint que incluya paginado de datos debe responder con un objeto JSON que contenga:

- `data`: array con los elementos resultantes de la consulta.
- `pagination`: objeto con la información de paginación de la respuesta.
- Todas las propiedades de respuesta JSON usan `camelCase`.
- `page` y `size` son opcionales para futuros endpoints paginados, salvo que el requerimiento indique explícitamente otra cosa.
- Si `page` no se informa, se usa `page=1`.
- Si `size` no se informa, se usa `size=20`.
- Si `page < 1`, `size < 1` o `size > 100`, la API responde `400 Bad Request` con formato `{"error": "mensaje"}`.

La estructura esperada es:

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

### Health

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/health` | Estado de la API |

### Confederaciones
 
 | Método | Ruta | Descripción |
 |--------|------|-------------|
 | `GET` | `/api/confederations` | Listar todas las confederaciones |
 | `GET` | `/api/confederations/:code` | Obtener confederación por code |

### Selecciones Nacionales

 | Método | Ruta | Descripción |
 |--------|------|-------------|
 | `GET` | `/api/teams` | Listar selecciones nacionales con filtros y paginación |
 | `GET` | `/api/teams/:code` | Obtener selección por código FIFA |

Parámetros soportados para `/api/teams`:

- `name`: búsqueda por contiene, case-insensitive.
- `confederation_code`: filtro por igualdad exacta, case-insensitive.
- `federation_name`: búsqueda por contiene, case-insensitive.
- `federation_code`: filtro por igualdad exacta, case-insensitive.
- `include_dissolved`: `true|false` (por defecto `false`).
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- `dissolutionDate` se expone en formato `YYYY-MM-DD` cuando aplica.
- `code` y `federationCode` se normalizan a mayúsculas.
- Se incluye el campo calculado `isDissolved`.

### Campeonatos Mundiales

 | Método | Ruta | Descripción |
 |--------|------|-------------|
 | `GET` | `/api/championships` | Listar ediciones de campeonatos mundiales con filtros y paginación |
 | `GET` | `/api/championships/:year` | Obtener detalle de una edición por año con estadísticas |
 | `GET` | `/api/championships/:year/fixture` | Obtener fixture completo de una edición por año |
 | `GET` | `/api/championships/:year/teams` | Listar selecciones participantes de una edición con filtros y paginación |

Parámetros soportados para `/api/championships`:

- `year`: filtro exacto por año del campeonato.
- `host`: búsqueda por nombre del país anfitrión (contiene, case-insensitive, sobre el nombre de la selección en `teams`).
- `confederation_code`: filtro por código de la confederación de los países anfitriones.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- `hostCodes` y `championCode` se normalizan a mayúsculas.
- Si no hay estadísticas cargadas para una edición, `stats` devuelve valores predeterminados (enteros en `0`, strings vacíos `""` y arrays vacíos `[]`).
- El fixture agrupa stages de tipo `group` con `groups[].matches` y `groups[].standings`; los stages `knockout` exponen `matches` directamente.

Parámetros soportados para `/api/championships/:year/teams`:

- `name`: búsqueda por nombre de selección (contiene, case-insensitive, sobre `teams.name`).
- `confederation_code`: filtro por igualdad sobre `teams.confederation_code`, normalizado a mayúsculas.
- `group_code`: filtro por igualdad sobre `championships_teams_stats.group_code`, normalizado a mayúsculas.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Si `:year` no es numérico, responde `400 Bad Request` con `{"error":"invalid year parameter"}`.
- Si `:year` es numérico pero no tiene selecciones asociadas, responde `200 OK` con `data: []` y metadata de paginación.
- `managers` devuelve string vacío `""` cuando no hay DTs asociados.
- Los resultados se ordenan por posición ascendente e instancia alcanzada descendente.

### Ejemplos de request
 
 **Listar confederaciones**
```bash
curl -H "API-Version: 1" http://localhost:8080/api/confederations
```

**Listar selecciones activas (paginado por defecto)**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/teams"
```

**Filtrar selecciones por nombre y confederación**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/teams?name=argen&confederation_code=CONMEBOL&page=1&size=20"
```

**Listar campeonatos mundiales (orden cronológico ascendente)**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships"
```

**Filtrar mundiales por confederación del anfitrión**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships?confederation_code=CONMEBOL"
```

**Obtener fixture completo de un mundial**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1978/fixture"
```

**Listar selecciones participantes de un mundial**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1930/teams?page=1&size=10"
```

**Filtrar selecciones participantes por nombre, confederación o grupo**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1930/teams?name=argentina&confederation_code=CONMEBOL&group_code=1"
```

**Obtener detalle completo de un mundial por año**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1986"
```

**Obtener selección por código FIFA**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/teams/code/urs"
```

---

## Deploy en Render

El proyecto incluye un archivo [`render.yaml`](render.yaml) para despliegue automático como **Infrastructure as Code**.

### Pasos para el primer deploy

1. **Conectar el repositorio** en [render.com](https://render.com) → _New_ → _Blueprint_
2. **Configurar la variable de entorno** `DATABASE_URL` manualmente en el dashboard de Render (apuntando al branch `main` de Neon)
3. **Aplicar las migraciones** antes del primer deploy (una sola vez):
   ```bash
   migrate -path db/migrations -database "$DATABASE_URL" up
   ```
4. Render detecta el `render.yaml` y despliega automáticamente en cada push a `main`

### Variables de entorno en Render

| Variable | Cómo se configura |
|---|---|
| `DATABASE_URL` | Manual en el dashboard de Render (valor secreto) |
| `GIN_MODE` | Definida en `render.yaml` como `release` |
| `PORT` | Inyectada automáticamente por Render |

> **Nota:** El archivo `.env` nunca debe subirse al repositorio. En producción, todas las variables se gestionan desde el dashboard o el `render.yaml`.

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
/.agents
  /instructions                 # documentación técnica del proyecto
  /prompts                      # prompts reutilizables para el IDE
  /rules                        # reglas y restricciones del agente
  /skills                       # habilidades específicas del agente
  /workflows                    # flujos de trabajo automatizados
/.coverage                      # reportes de cobertura de tests
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
