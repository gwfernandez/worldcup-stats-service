---
description: Workflow que guía al agente en la resolución completa de un issue de GitHub, usando herramientas de integración para la asignación, comentarios y PRs, asegurando cobertura de tests del 90% y trazabilidad total en GitHub.
---

# Workflow — Resolver Issue de GitHub

## Contexto del proyecto
- Repositorio: worldcups-api
- Lenguaje: Go 1.23
- Framework: Gin
- Arquitectura: Clean Architecture (handler → service → repository)
- Rama principal: main
- Convención de branches: indicada en el campo "Rama sugerida" del issue de GitHub

---

## Fase 0 — Validación del entorno

Verificar que el entorno está en condiciones antes de arrancar. Si algún paso falla → **detener y reportar**. No continuar con un entorno roto.

### 1. Verificar estado del repositorio

```bash
git fetch origin
git status
git diff main origin/main --stat
```

Condiciones requeridas:
- Working tree limpio (sin cambios sin commitear ni archivos sin trackear)
- `main` local sincronizado con `origin/main`

Si no se cumplen → **detener y reportar**.

### 2. Verificar dependencias

```bash
go mod verify
go mod tidy
git diff go.mod go.sum
```

Si `go mod tidy` genera cambios en `go.mod` o `go.sum` → reportar advertencia antes de continuar. No commitear cambios de dependencias como parte del issue.

### 3. Verificar que el proyecto compila

```bash
go build ./...
```

Si hay errores de compilación → **detener y reportar**. No continuar con código que no compila.

### 4. Ejecutar tests y guardar baseline

```bash
mkdir -p .coverage
go test -coverprofile=.coverage/baseline.out ./... 2>&1 | tee /tmp/baseline_tests.txt
go tool cover -func=.coverage/baseline.out | tee /tmp/baseline_coverage.txt
```

Capturar y registrar en memoria:
- Cantidad de tests que pasan / fallan
- Coverage por paquete (`handler`, `service`, `repository`, `config`)
- Coverage total

Si hay tests fallando en `main` antes del issue → **detener y reportar**. Los tests rotos preexistentes no son responsabilidad del issue.

Reportar el resumen antes de continuar:

```
✅ Fase 0 — Entorno validado

Repo:           main limpio, sincronizado con origin
Dependencias:   OK (go.mod/go.sum sin cambios)
Build:          OK
Tests:          X passing, 0 failing
Coverage total: X%
  handler:      X%
  service:      X%
  repository:   X%
  config:       X%

Listo para continuar con la Fase 1.
```

---

## Fase 1 — Preparación

1. Leer el issue de GitHub completo usando las herramientas de búsqueda/lectura y entender el requerimiento
2. Asignarme el issue usando la herramienta de actualización de issues de GitHub
3. Leer el campo "Rama sugerida" del issue
   - Si existe → usar ese nombre
   - Si no existe → solicitarme el nombre antes de continuar
4. Crear el branch localmente con ese nombre desde `main`
5. Hacer checkout al branch creado
6. Cambiar el estado del issue a **En curso** publicando un comentario en el issue de GitHub mediante la integración

---

## Fase 2 — Planificación

1. Analizar el issue de GitHub, tomando como base la sección `📋 Tareas Técnicas` del SDD si existe, y elaborar un plan de acción detallado que incluya:
   - Archivos a crear o modificar
   - Capas involucradas (handler / service / repository)
   - Queries SQL nuevas si aplica
   - Tests unitarios necesarios que cubran los `✅ Criterios de Aceptación`
2. Presentarme el plan (Implementation Plan) y **esperár mi confirmación antes de continuar**
3. Una vez que apruebe el plan, **DEBES usar la herramienta de GitHub para publicar el "Implementation Plan" completo como un comentario en el issue original**.

---

## Fase 3 — Desarrollo

1. Ejecutar el plan de acción confirmado
2. Seguir las convenciones del proyecto:
   - Nombres de archivos: `{entidad}_{capa}.go`
   - Tests: `{entidad}_{capa}_test.go` en el mismo directorio
   - Errores HTTP con formato: `{"error": "mensaje"}`
   - Prefijo de rutas: `/api/`
3. No modificar archivos fuera del alcance del issue de GitHub
4. Si el issue contiene una lista de tareas (tasks), realizar un commit independiente utilizando la skill `semantic-commit` por cada tarea completada. El mensaje del commit debe reflejar fielmente la tarea realizada. En caso de no haber una lista, realizar commits por cada hito lógico finalizado.
5. Respetar estrictamente la sección `🚫 Fuera de Alcance` del issue para evitar cambios innecesarios y mantener el foco.

---

## Fase 3.5 — Auditoría de calidad

1. Ejecutar la skill `code-quality-go` en **modo Git Diff** para auditar todos los archivos `.go` modificados en el branch actual
2. Revisar los hallazgos reportados y aplicar las correcciones aprobadas
3. Si hay hallazgos 🔴 Críticos → **corregir obligatoriamente** antes de continuar
4. Si hay hallazgos 🟠 Importantes → corregir salvo decisión explícita del usuario de postergarlos
5. Registrar las correcciones con la skill `semantic-commit` si corresponde

> Esta fase es obligatoria. No avanzar a Testing sin haber ejecutado la auditoría de calidad.

---

## Fase 4 — Testing

1. Crear los tests unitarios necesarios para alcanzar un **coverage mínimo del 90%**
   - Handler: black-box testing (`package handler_test`)
   - Service: white-box testing (`package service`)
   - Repository: white-box testing (`package repository`)
2. Ejecutar todos los tests:
```bash
   go test ./... -coverprofile=.coverage/coverage.out
   go tool cover -func=.coverage/coverage.out
```
3. Si algún test falla → corregir el código o el test hasta que todos pasen
4. Si el coverage es menor al 90% → agregar los tests faltantes
5. Validar que los tests cubren todos los `✅ Criterios de Aceptación` definidos en el issue original.
6. Una vez que todos los tests pasen y el coverage sea ≥ 90%, **usar la integración de GitHub para agregar como comentario en el issue original** el reporte de ejecución como evidencia

---

## Fase 5 — Pull Request

1. Usar la herramienta de creación de Pull Requests de GitHub para generar un PR desde el branch actual hacia `main`:

```bash
gh pr create \
  --base main \
  --title "<tipo>(<scope>): <descripción en español> (#<numero_issue>)" \
  --body "$(cat <<'EOF'
## Closes #<numero-de-issue>

## ¿Qué hace este PR?

<descripción de 2-3 oraciones en español explicando qué se implementó y por qué>

## Cambios realizados

### Archivos creados
- `internal/domain/...` — descripción
- `internal/repository/...` — descripción
- `internal/service/...` — descripción
- `internal/handler/...` — descripción
- `db/queries/...` — descripción (si aplica)
- `db/migrations/...` — descripción (si aplica)

### Archivos modificados
- `cmd/main.go` — descripción del cambio (si aplica)
- `README.md` — descripción del cambio (si aplica)

## Nuevos endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/...` | ... |

(Omitir sección si no hay nuevos endpoints)

## Decisiones de implementación

<Si se tomó alguna decisión no obvia de arquitectura o diseño, documentarla aquí. Si todo sigue los patrones establecidos, escribir "Sigue los patrones establecidos del proyecto.">

## Testing

| Métrica | Antes | Después |
|---|---|---|
| Tests passing | <baseline> | <actual> |
| Tests nuevos | — | <cantidad> |
| Coverage handler | <baseline> | <actual> |
| Coverage service | <baseline> | <actual> |
| Coverage repository | <baseline> | <actual> |
| Coverage total | <baseline> | <actual> |

Edge cases cubiertos:
- <listar los edge cases del issue que tienen test>

## Puntos de atención para el reviewer

<Si hay algo que merece atención especial: decisiones de diseño no obvias, cambios de comportamiento, impacto en otros módulos. Si no hay nada, escribir "Ninguno".>

## Checklist

- [ ] El código compila sin errores ni warnings
- [ ] Coverage ≥ 90% en paquetes modificados
- [ ] Todas las funciones públicas tienen godoc
- [ ] Sin datos sensibles en logs
- [ ] Documentación actualizada (README, AGENTS.md si aplica)
- [ ] API de solo lectura respetada (solo GET)
- [ ] Criterios de aceptación del issue cumplidos

EOF
)"
```

2. Asignarme como reviewer (si la API lo permite, si no, dejar documentado)
3. **Usar la herramienta de comentarios de GitHub** para publicar el documento de resumen ("Walkthrough") como un comentario final en el issue original, indicando que el trabajo ha concluido:

```
✅ Trabajo completado — PR listo para review

## Resumen

<descripción de 2-3 oraciones de qué se implementó y qué problema resuelve>

## Cambios principales

- `<archivo>`: <descripción del cambio>
- `<archivo>`: <descripción del cambio>

## Testing

- Tests nuevos: X
- Coverage total: X% (sin regresiones)
- Edge cases cubiertos: <lista resumida>

## Pull Request

<URL del PR> — listo para review de @gwfernandez
```