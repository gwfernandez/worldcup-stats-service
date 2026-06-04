---
description: Stage 1 — Análisis, plan y desarrollo
---

Stage 1 — Análisis, plan y desarrollo

## Propósito

Leer el issue, demostrar entendimiento explícito del problema, planificar la implementación respetando la Clean Architecture, y desarrollar el código junto con sus tests. Esta etapa tiene **dos checkpoints humanos**: uno antes de codear y uno al terminar.

---

## Input requerido

Invocar este workflow indicando el número de issue:

```
/stage-1-analysis.md issue:#<número>
```

---

## Pasos

### 1. Leer el issue y preparar el contexto

```bash
gh issue view <número> --json title,body,labels,assignees,comments
```

Del body del issue, extraer y registrar el campo **"Rama sugerida"**:
- Si existe → usar ese nombre exacto en el Paso 3
- Si no existe → usar la convención `<tipo>/<número-issue>-<descripcion-corta-kebab-case>`

A continuación, ejecutar las siguientes acciones de trazabilidad en GitHub:

1. **Asignarme el issue** usando la integración de GitHub
2. **Publicar el comentario de inicio** en el issue:

```
🔄 Issue tomado — comenzando análisis y planificación.
```

---

### 2. Reformular el issue — CHECKPOINT A

Antes de escribir código, demostrar entendimiento explícito. Presentar al usuario:

**Título del issue:** (repetir)

**Entendimiento:** Reformular con palabras propias qué hay que implementar y por qué.

**Capas afectadas y archivos a crear o modificar:**

| Capa | Archivo | Acción |
|---|---|---|
| `domain` | `internal/domain/...` | crear / modificar |
| `repository` | `internal/repository/...` | crear / modificar |
| `service` | `internal/service/...` | crear / modificar |
| `handler` | `internal/handler/...` | crear / modificar |
| SQL queries | `db/queries/...` | crear / modificar |
| Migración | `db/migrations/...` | crear (si aplica) |

Recordar: las dependencias siempre apuntan hacia adentro (`handler → service → repository`). El `service` nunca importa `handler`. El `repository` nunca importa `service` ni `handler`.

**Nuevos endpoints (si aplica):**
- Método, ruta y descripción
- Confirmar que son solo `GET` — este proyecto es de **solo lectura**. No implementar `POST`, `PUT`, `PATCH` ni `DELETE` salvo excepción justificada explícitamente en el issue.

**Queries SQL necesarias (si aplica):**
- Describir las queries que hay que escribir en `db/queries/`
- Indicar si requieren migración nueva

**Edge cases identificados:**
- Recurso no encontrado (404)
- Parámetros inválidos o malformados (400)
- Errores de base de datos (500)
- Casos límite del dominio específico del issue

**Nuevas dependencias (si aplica):**
- Si el issue requiere una nueva dependencia externa, justificar técnicamente por qué no puede resolverse con la stdlib o las dependencias existentes. Requiere aprobación explícita.

**Preguntas o ambigüedades** (si las hay)

> ⏸️ **Esperar aprobación del usuario antes de continuar.**
> Si el usuario corrige el entendimiento, actualizar el plan y volver a presentar.

Una vez aprobado, **publicar el Implementation Plan completo como comentario en el issue** usando la integración de GitHub. Usar este formato:

```
📋 Implementation Plan — Issue #<número>

**Entendimiento:** <reformulación>

**Capas y archivos:**
<tabla de archivos>

**Nuevos endpoints (si aplica):**
<lista>

**Queries SQL (si aplica):**
<descripción>

**Edge cases identificados:**
<lista>
```

---

### 3. Crear la rama

Usar el nombre obtenido del campo **"Rama sugerida"** del issue (leído en el Paso 1):
- Si existe → usar ese nombre exacto
- Si no existe → usar la convención: `<tipo>/<número-issue>-<descripcion-corta-kebab-case>`

```bash
git checkout main
git pull origin main
git checkout -b <nombre-de-rama>
```

Ejemplos:
- Con "Rama sugerida": `feat/teams-endpoint`
- Con convención: `feat/42-endpoint-goleadores`

---

### 4. Escribir las queries SQL (si aplica)

Si el issue requiere nuevas queries, escribirlas primero en `db/queries/`:

```sql
-- db/queries/scorers.sql

-- name: ListScorers :many
SELECT ...
```

Luego regenerar el código tipado con sqlc:

```bash
sqlc generate
```

Verificar que el código generado en `/internal/repository` compila:

```bash
go build ./...
```

**Si el issue requiere una nueva migración**, presentar el SQL al usuario antes de ejecutarla. Las migraciones son acciones de alto impacto — **requieren confirmación explícita**.

Formato del nombre: `NNN_descripcion_corta.up.sql` / `NNN_descripcion_corta.down.sql`

---

### 5. Implementar el código siguiendo Clean Architecture

Orden de implementación recomendado:

**5.1 — `internal/domain/`**
- Definir o actualizar structs de entidades
- Sin dependencias externas, solo tipos Go nativos

**5.2 — `internal/repository/`**
- Implementar la interfaz del repository
- Usar el código generado por sqlc con pgx v5
- Nunca escribir SQL directamente en Go — siempre vía sqlc

**5.3 — `internal/service/`**
- Implementar la lógica de negocio
- Depender de la **interfaz** del repository, nunca de la implementación concreta
- Sin conocimiento de HTTP ni SQL

**5.4 — `internal/handler/`**
- Implementar el handler de Gin
- Depender de la **interfaz** del service
- Responsabilidades: parsear input, llamar al service, retornar JSON

**Checklist de implementación:**

- [ ] Todas las funciones públicas tienen comentario godoc
- [ ] Errores siempre propagados explícitamente, nunca ignorados con `_`
- [ ] Formato de error consistente: `{"error": "mensaje"}`
- [ ] Nombres de funciones de handlers: `List`, `GetByID`, `Create`, `Update`, `Delete`
- [ ] Prefijo de rutas: `/api/`
- [ ] Sin ORMs — toda interacción con la DB vía sqlc
- [ ] Sin `fmt.Println` ni `log.Printf` con datos sensibles (passwords, tokens, DATABASE_URL)
- [ ] Sin variables de entorno hardcodeadas — usar siempre el paquete `config`
- [ ] Seguir idioms de Effective Go y Google Go Style

---

### 6. Implementar tests

Por cada archivo de código creado o modificado, crear o actualizar su test:

**Estructura por capa:**

| Capa | Archivo de test | Mock |
|---|---|---|
| `handler` | `internal/handler/nombre_test.go` | mock del service con `testify/mock` |
| `service` | `internal/service/nombre_test.go` | mock del repository con `testify/mock` |
| `repository` | `internal/repository/nombre_test.go` | `pgxmock` |
| `config` | `config/config_test.go` | variables de entorno de test |

**Reglas:**
- Coverage objetivo en paquetes modificados: **≥ 90%**
- Testear: happy path, recurso no encontrado, error de base de datos, parámetros inválidos
- Un test por comportamiento observable
- No testear el código generado por sqlc en `db/sqlc/` (excluido del coverage)

---

### 7. Verificación rápida pre-checkpoint

```bash
go build ./...
go vet ./...
go test ./... 2>&1 | tail -20
```

Corregir errores de compilación, vet o tests rotos antes de presentar el trabajo.

---

### 8. Presentar resultado — CHECKPOINT B

```
✅ Desarrollo completado

Rama:    feat/42-endpoint-goleadores
Issue:   #42 — Endpoint de goleadores por mundial

Archivos creados:
  + internal/domain/scorer.go
  + internal/repository/scorer_repository.go
  + internal/service/scorer_service.go
  + internal/handler/scorer_handler.go
  + internal/handler/scorer_handler_test.go
  + internal/service/scorer_service_test.go
  + internal/repository/scorer_repository_test.go
  + db/queries/scorers.sql

Archivos modificados:
  ~ cmd/main.go  (registro de rutas)

Tests nuevos:   18
Tests totales:  105 passing, 0 failing
Build:          OK

Listo para continuar con /stage-2-audit.md
```

> ⏸️ **Esperar aprobación del usuario antes de continuar.**
