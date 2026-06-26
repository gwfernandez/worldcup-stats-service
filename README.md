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
| `CORS_ALLOWED_ORIGINS` | Origenes permitidos para CORS, separados por coma | `http://localhost:5173,https://app.example.com` |

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

### Internacionalización

La API permite solicitar datos localizados mediante el header HTTP `Accept-Language`.
Si no se envía el header, o si el idioma solicitado no está soportado, se utiliza `es` por defecto.

Idiomas soportados inicialmente:

- `es`: español, idioma por defecto.
- `en`: inglés.

Cuando falta una traducción para el idioma solicitado, la API responde el valor base almacenado en la tabla principal. Para confederaciones, el campo `confederations.name` funciona como fallback en español; para selecciones, el fallback es `teams.name`.

Actualmente el header `Accept-Language` localiza nombres de confederaciones y nombres de selecciones en los endpoints que los exponen o filtran: `/api/confederations`, `/api/teams`, `/api/champions`, `/api/standings`, `/api/scorers`, `/api/players/:playerId/goals`, `/api/championships`, `/api/championships/:year/teams`, `/api/championships/:year/stadiums` y `/api/championships/:year/stadiums/:stadiumId`.

**Ejemplo:** `Accept-Language: en`

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

### Convención de query params

Todos los query params de la API usan `camelCase`.

Los parámetros de una sola palabra se mantienen sin cambios, por ejemplo `name`, `host`, `year`, `page` y `size`.

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

- `name`: búsqueda por contiene, case-insensitive, sobre el nombre localizado de la selección según `Accept-Language`.
- `confederationCode`: filtro por igualdad exacta, case-insensitive.
- `federationName`: búsqueda por contiene, case-insensitive.
- `federationCode`: filtro por igualdad exacta, case-insensitive.
- `includeDissolved`: `true|false` (por defecto `false`).
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- `dissolutionDate` se expone en formato `YYYY-MM-DD` cuando aplica.
- `code` y `federationCode` se normalizan a mayúsculas.
- Se incluye el campo calculado `isDissolved`.
- `name` se resuelve según `Accept-Language` con fallback a `teams.name`.

### Campeones Mundiales

 | Método | Ruta | Descripción |
 |--------|------|-------------|
 | `GET` | `/api/champions` | Listar tabla histórica de campeones mundiales con paginación |
 | `GET` | `/api/champions/:teamCode` | Listar finales ganadas por una selección con paginación |

Parámetros soportados para `/api/champions`:

- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- No soporta filtros.
- Los resultados se ordenan por `wins` descendente y, ante empates, por `team.code` ascendente.
- `team.name` se resuelve según `Accept-Language` con fallback a `teams.name`.
- `years` se expone como array de números ordenado ascendentemente.

Parámetros soportados para `/api/champions/:teamCode`:

- `teamCode`: código unificado de la selección, normalizado a mayúsculas.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Devuelve exclusivamente finales ganadas, ordenadas cronológicamente por `year`.
- Incluye el partido decisivo de 1950 como final histórica.
- `hostCodes` expone los anfitriones de cada edición como un array de `SimpleTeam`.
- `homeTeam` y `awayTeam` exponen `code` y `name`.
- Los nombres de `hostCodes`, `homeTeam` y `awayTeam` se resuelven según `Accept-Language`.
- Si la selección no tiene finales ganadas o el código no existe, responde `200 OK` con `data: []`.
- `matchDate`, `matchTime` y los resultados pueden ser `null` cuando el dato no está disponible.

Ejemplo de respuesta:

```json
{
  "data": [
    {
      "team": {
        "code": "BRA",
        "name": "Brasil"
      },
      "wins": 5,
      "years": [1958, 1962, 1970, 1994, 2002]
    }
  ],
  "pagination": {
    "page": 1,
    "size": 20,
    "totalElements": 8,
    "totalPages": 1,
    "hasNext": false,
    "hasPrevious": false
  }
}
```

Ejemplo de finales ganadas:

```json
{
  "data": [
    {
      "year": 2022,
      "hostCodes": [
        {
          "code": "QAT",
          "name": "Catar"
        }
      ],
      "matchDate": "2022-12-18",
      "matchTime": "18:00:00",
      "homeTeam": {
        "code": "ARG",
        "name": "Argentina"
      },
      "homeTeamScore": 3,
      "homeTeamScorePenalties": 4,
      "awayTeam": {
        "code": "FRA",
        "name": "Francia"
      },
      "awayTeamScore": 3,
      "awayTeamScorePenalties": 2
    }
  ],
  "pagination": {
    "page": 1,
    "size": 20,
    "totalElements": 3,
    "totalPages": 1,
    "hasNext": false,
    "hasPrevious": false
  }
}
```

### Tabla Histórica de Posiciones

 | Método | Ruta | Descripción |
 |--------|------|-------------|
 | `GET` | `/api/standings` | Listar tabla histórica de posiciones de todos los mundiales con filtros y paginación |

Parámetros soportados para `/api/standings`:

- `name`: búsqueda por nombre de selección localizado según `Accept-Language`, case-insensitive y por contiene.
- `confederationCode`: filtro por igualdad exacta sobre el código de confederación. La API normaliza el valor a mayúsculas.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Ejemplo de respuesta:

```json
{
  "data": [
    {
      "team": {
        "code": "BRA",
        "name": "Brasil"
      },
      "confederationCode": "CONMEBOL",
      "matchesPlayed": 114,
      "wins": 79,
      "draws": 14,
      "losses": 21,
      "goalsFor": 237,
      "goalsAgainst": 108,
      "goalDifference": 129,
      "points": 193,
      "unifiedPoints": 237,
      "position": 1,
      "unifiedPosition": 1
    }
  ],
  "pagination": {
    "page": 1,
    "size": 20,
    "totalElements": 1,
    "totalPages": 1,
    "hasNext": false,
    "hasPrevious": false
  }
}
```

Notas de respuesta:

- `confederationCode` se expone al nivel superior de cada posición y se normaliza a mayúsculas.
- `team` conserva la estructura `SimpleTeam` con `code` y el nombre localizado según `Accept-Language`.
- Si no hay posiciones asociadas a los filtros, responde `200 OK` con `data: []` y metadata de paginación.

### Tabla Histórica de Goleadores

 | Método | Ruta | Descripción |
 |--------|------|-------------|
 | `GET` | `/api/scorers` | Listar tabla histórica de goleadores de todos los mundiales con filtros y paginación |
 | `GET` | `/api/scorers/:playerId` | Obtener los datos personales y todos los goles válidos de un goleador |

Parámetros soportados para `/api/scorers`:

- `name`: búsqueda por nombre o apellido de jugador, case-insensitive y por contiene.
- `teamCode`: filtro por igualdad exacta sobre `teams.unified_code`, normalizado a mayúsculas.
- `confederationCode`: filtro por igualdad exacta sobre `teams.confederation_code`, normalizado a mayúsculas.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Si no hay goleadores asociados a los filtros, responde `200 OK` con `data: []` y metadata de paginación.
- Los resultados se ordenan por `goals` descendente y `fullName` ascendente.
- `team.code`, `listTeams` y `confederationCode` se normalizan a mayúsculas.
- La respuesta expone `fullName`, `team`, `goals`, `listTeams` y `confederationCode`.
- `team.name` corresponde al nombre localizado de la selección principal según `Accept-Language`, con fallback a `teams.name`.

Para `/api/scorers/:playerId`:

- `playerId` debe ser un identificador positivo.
- La respuesta es un objeto directo, sin paginación.
- `championships` conserva el orden configurado en el jugador.
- `teams` usa un array de `SimpleTeam`, conserva el orden configurado y respeta `Accept-Language`.
- `goals` contiene todos los goles válidos del jugador, excluye autogoles y conserva el mismo contrato y orden que `/api/players/:playerId/goals`.
- `hosts` y `opponentTeam` respetan `Accept-Language`.
- Los arrays sin datos se serializan como `[]`.
- Un jugador existente sin goles responde `200 OK` con `goals: []`.
- Un jugador inexistente responde `404 Not Found`.

Ejemplo de respuesta:

```json
{
  "data": [
    {
      "fullName": "Lionel Messi",
      "team": {
        "code": "ARG",
        "name": "Argentina"
      },
      "goals": 13,
      "listTeams": ["ARG"],
      "confederationCode": "CONMEBOL"
    }
  ],
  "pagination": {
    "page": 1,
    "size": 20,
    "totalElements": 1,
    "totalPages": 1,
    "hasNext": false,
    "hasPrevious": false
  }
}
```

Ejemplo de detalle:

```json
{
  "id": 1524,
  "firstName": "Lionel",
  "lastName": "Messi",
  "position": "FW",
  "championships": [2006, 2010, 2014, 2018, 2022],
  "teams": [
    {
      "code": "ARG",
      "name": "Argentina"
    }
  ],
  "goals": [
    {
      "year": 2022,
      "hosts": [
        {
          "code": "QAT",
          "name": "Catar"
        }
      ],
      "matchDate": "2022-12-18",
      "opponentTeam": {
        "code": "FRA",
        "name": "Francia"
      },
      "minuteRegular": 23,
      "penalty": true,
      "stage": "final"
    }
  ]
}
```

### Goles por Jugador

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/players/:playerId/goals` | Listar goles válidos de un jugador con filtro opcional por Mundial y paginación |

Parámetros soportados:

- `playerId`: identificador positivo del jugador informado en la ruta.
- `year`: año opcional del Mundial; debe ser un entero positivo.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Se excluyen los autogoles.
- Los resultados se ordenan por `matchDate` ascendente, `minuteRegular` ascendente e identificador de gol.
- `hosts` expone los anfitriones de la edición como un array de `SimpleTeam`, conserva el orden configurado y respeta `Accept-Language`.
- `opponentTeam` usa la estructura `SimpleTeam`; su código se normaliza a mayúsculas y su nombre respeta `Accept-Language`.
- `matchDate`, `penalty` y `stage` conservan `null` cuando el dato no está disponible.
- Si no hay resultados, responde `200 OK` con `data: []` y metadata de paginación.

Ejemplo de respuesta:

```json
{
  "data": [
    {
      "year": 2018,
      "hosts": [
        {
          "code": "RUS",
          "name": "Rusia"
        }
      ],
      "matchDate": "2018-06-16",
      "opponentTeam": {
        "code": "ISL",
        "name": "Islandia"
      },
      "minuteRegular": 64,
      "penalty": true,
      "stage": "group_stage"
    }
  ],
  "pagination": {
    "page": 1,
    "size": 20,
    "totalElements": 1,
    "totalPages": 1,
    "hasNext": false,
    "hasPrevious": false
  }
}
```

### Campeonatos Mundiales

 | Método | Ruta | Descripción |
 |--------|------|-------------|
 | `GET` | `/api/championships` | Listar ediciones de campeonatos mundiales con filtros y paginación |
 | `GET` | `/api/championships/:year` | Obtener detalle de una edición por año con estadísticas |
 | `GET` | `/api/championships/:year/fixture` | Obtener fixture completo de una edición por año |
 | `GET` | `/api/championships/:year/teams` | Listar selecciones participantes de una edición con filtros y paginación |
 | `GET` | `/api/championships/:year/scorers` | Listar goleadores de una edición con filtros y paginación |
 | `GET` | `/api/championships/:year/squads/:teamCode` | Listar jugadores de una selección en una edición con paginación |
 | `GET` | `/api/championships/:year/stadiums` | Listar estadios utilizados en una edición con filtros y paginación |
 | `GET` | `/api/championships/:year/stadiums/:stadiumId` | Listar partidos jugados en un estadio durante una edición con paginación |
 | `GET` | `/api/championships/:year/standings` | Listar tabla de posiciones de una edición con paginación |

Parámetros soportados para `/api/championships`:

- `year`: filtro exacto por año del campeonato.
- `host`: búsqueda por nombre del país anfitrión (contiene, case-insensitive, sobre el nombre localizado de la selección según `Accept-Language`).
- `confederationCode`: filtro por código de la confederación de los países anfitriones.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- `hosts[].code` y `champion.code` se normalizan a mayúsculas.
- `confederationCodes` expone los códigos de confederaciones organizadoras normalizados a mayúsculas.
- `hosts[].name` corresponde al nombre localizado del anfitrión según `Accept-Language`, con fallback a español y luego al código.
- `champion.name` corresponde al nombre localizado del campeón según `Accept-Language`, con fallback a español y luego al código.
- `stats.runnerUp`, `stats.thirdPlace` y `stats.fourthPlace` exponen `code` y `name` localizados desde la caché de selecciones.
- Si no hay estadísticas cargadas para una edición, `stats` devuelve valores predeterminados (enteros en `0`, strings vacíos `""` y arrays vacíos `[]`).
- El fixture agrupa stages de tipo `group` con `groups[].matches` y `groups[].standings`; los stages `knockout` exponen `matches` directamente.
- En el fixture, cada match expone `homeTeam` y `awayTeam`; cada standing de grupo expone `team`.
- Los nombres de selecciones del fixture se resuelven según `Accept-Language` con fallback a `teams.name`.

Parámetros soportados para `/api/championships/:year/teams`:

- `name`: búsqueda por nombre de selección localizado según `Accept-Language` (contiene, case-insensitive, con fallback a `teams.name`).
- `confederationCode`: filtro por igualdad sobre `teams.confederation_code`, normalizado a mayúsculas.
- `groupCode`: filtro por igualdad sobre `championships_teams_stats.group_code`, normalizado a mayúsculas.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Si `:year` no es numérico, responde `400 Bad Request` con `{"error":"invalid year parameter"}`.
- Si `:year` es numérico pero no tiene selecciones asociadas, responde `200 OK` con `data: []` y metadata de paginación.
- `managers` devuelve string vacío `""` cuando no hay DTs asociados.
- Los resultados se ordenan por posición ascendente e instancia alcanzada descendente.
- La respuesta expone `year`, `team`, `confederationCode`, `groupCode`, `stageReached` y `managers`.
- `team.name` se resuelve según `Accept-Language` con fallback a `teams.name`.

Parámetros soportados para `/api/championships/:year/scorers`:

- `name`: búsqueda por nombre o apellido del jugador (contiene, case-insensitive, sobre `players.first_name` y `players.last_name`).
- `teamCode`: filtro por igualdad sobre `squads_stats.team_code`, normalizado a mayúsculas.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Si `:year` no es numérico, responde `400 Bad Request` con `{"error":"invalid year parameter"}`.
- Si `:year` es numérico pero no tiene goleadores asociados, responde `200 OK` con `data: []` y metadata de paginación.
- Los resultados se ordenan por `goals` descendente y `fullName` ascendente.
- `fullName` y `team` se exponen en `camelCase`; `team.code` se normaliza a mayúsculas.
- `team.name` corresponde al nombre localizado de la selección según `Accept-Language`, con fallback a `teams.name`.

Parámetros soportados para `/api/championships/:year/squads/:teamCode`:

- `year`: año del Mundial informado en la ruta; debe ser numérico.
- `teamCode`: código unificado de la selección informado en la ruta, normalizado a mayúsculas.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Si `:year` no es numérico, responde `400 Bad Request` con `{"error":"invalid year parameter"}`.
- Si no hay jugadores asociados para `:year` y `:teamCode`, responde `200 OK` con `data: []` y metadata de paginación.
- Los resultados se ordenan por `position` ascendente, `lastName` ascendente y `firstName` ascendente.
- La respuesta expone `playerId`, `firstName`, `lastName`, `position` y `shirtNumber`.
- `position` se obtiene de `squads.position`.
- `position` y `shirtNumber` conservan `null` cuando no existe un valor en la base.

Parámetros soportados para `/api/championships/:year/stadiums`:

- `name`: búsqueda por nombre de estadio (contiene, case-insensitive, sobre `stadiums.name`).
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Si `:year` no es numérico, responde `400 Bad Request` con `{"error":"invalid year parameter"}`.
- Si `:year` es numérico pero no tiene estadios asociados, responde `200 OK` con `data: []` y metadata de paginación.
- Los resultados se ordenan por `matchesPlayed` descendente y `name` ascendente.
- La respuesta expone `id`, `name`, `cityName`, `country`, `capacity` y `matchesPlayed`.
- `country` usa la estructura `SimpleTeam`, se obtiene de `stadiums.country` y respeta `Accept-Language`.
- Si el estadio no tiene país cargado, `country` se serializa como `null`.

Parámetros soportados para `/api/championships/:year/stadiums/:stadiumId`:

- `year`: año del Mundial informado en la ruta; debe ser numérico positivo.
- `stadiumId`: identificador del estadio informado en la ruta; debe ser numérico positivo.
- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Si `:year` no es numérico positivo, responde `400 Bad Request` con `{"error":"invalid year parameter"}`.
- Si `:stadiumId` no es numérico positivo, responde `400 Bad Request` con `{"error":"invalid stadiumId parameter"}`.
- Si no hay partidos asociados para `:year` y `:stadiumId`, responde `200 OK` con `data: []` y metadata de paginación.
- Los resultados se ordenan por `matchDate` ascendente, con desempates por `matchTime` e identificador interno del partido.
- La respuesta expone `year`, `hosts`, `stage`, `groupCode`, `matchDate`, `matchTime`, `homeTeam`, `homeTeamScore`, `homeTeamScorePenalties`, `awayTeam`, `awayTeamScore` y `awayTeamScorePenalties`.
- `hosts`, `homeTeam` y `awayTeam` usan la estructura `SimpleTeam`; sus nombres respetan `Accept-Language`.
- `hosts` se obtiene de `championships.host_codes`, conserva el orden configurado para la edición y normaliza los códigos a mayúsculas.
- `stage`, `groupCode`, `matchDate`, `matchTime`, scores y penales conservan `null` cuando no existe un valor en la base.

Parámetros soportados para `/api/championships/:year/standings`:

- `page`: número de página (base 1, por defecto `1`).
- `size`: tamaño de página (por defecto `20`, máximo `100`).

Notas de respuesta:

- Si `:year` no es numérico, responde `400 Bad Request` con `{"error":"invalid year parameter"}`.
- Si `:year` es numérico pero no tiene posiciones asociadas, responde `200 OK` con `data: []` y metadata de paginación.
- No soporta filtros adicionales.
- Los resultados se ordenan por `position` ascendente e instancia alcanzada.
- `team.code` y `groupCode` se normalizan a mayúsculas.
- La respuesta expone `team`, `groupCode`, `matchesPlayed`, `wins`, `draws`, `losses`, `goalsFor`, `goalsAgainst`, `goalDifference`, `points`, `unifiedPoints`, `position` y `performance`.
- `team.name` se resuelve según `Accept-Language` con fallback a `teams.name`.

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
curl -H "API-Version: 1" "http://localhost:8080/api/teams?name=argen&confederationCode=CONMEBOL&page=1&size=20"
```

**Listar selecciones en inglés**
```bash
curl -H "API-Version: 1" -H "Accept-Language: en" "http://localhost:8080/api/teams"
```

**Listar campeones mundiales**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/champions?page=1&size=10"
```

**Listar finales ganadas por una selección**
```bash
curl -H "API-Version: 1" -H "Accept-Language: es" "http://localhost:8080/api/champions/ARG?page=1&size=20"
```

**Listar tabla histórica de goleadores**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/scorers?page=1&size=10"
```

**Filtrar tabla histórica de goleadores**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/scorers?name=messi&teamCode=ARG&confederationCode=CONMEBOL"
```

**Obtener detalle de un goleador**
```bash
curl -H "API-Version: 1" -H "Accept-Language: es" "http://localhost:8080/api/scorers/1524"
```

**Listar campeonatos mundiales (orden cronológico ascendente)**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships"
```

**Filtrar mundiales por confederación del anfitrión**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships?confederationCode=CONMEBOL"
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
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1930/teams?name=argentina&confederationCode=CONMEBOL&groupCode=1"
```

**Listar goleadores de un mundial**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1930/scorers?page=1&size=10"
```

**Filtrar goleadores por nombre y selección**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1930/scorers?name=stabile&teamCode=ARG"
```

**Listar estadios utilizados en un mundial**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1930/stadiums?page=1&size=10"
```

**Filtrar estadios por nombre**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1930/stadiums?name=centenario"
```

**Listar tabla de posiciones de un mundial**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1930/standings?page=1&size=10"
```

**Obtener detalle completo de un mundial por año**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/championships/1986"
```

**Obtener selección por código FIFA**
```bash
curl -H "API-Version: 1" "http://localhost:8080/api/teams/urs"
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
