---
description: Stage 0 — Setup y validación de entorno
---

# Stage 0 — Setup y validación de entorno

## Propósito

Verificar que el entorno está en condiciones para arrancar el issue. Si algo falla, reportar con claridad y detener el workflow. **No intentar corregir ningún problema encontrado.**

---

## Pasos

### 1. Verificar versiones del entorno

Ejecutar y reportar la salida:

```bash
go version
sqlc version
migrate -version
```

Versión mínima requerida:
- Go >= 1.23

Si no se cumple → **detener y reportar**.

---

### 2. Verificar variables de entorno

```bash
test -f .env && echo ".env existe" || echo ".env NO existe"
grep -q "DATABASE_URL" .env && echo "DATABASE_URL presente" || echo "DATABASE_URL FALTANTE"
grep -q "PORT" .env && echo "PORT presente" || echo "PORT FALTANTE"
grep -q "GIN_MODE" .env && echo "GIN_MODE presente" || echo "GIN_MODE FALTANTE"
```

Si falta alguna variable requerida → **detener y reportar**.

Verificar que `.env` **no está trackeado por git**:

```bash
git check-ignore .env && echo "OK — .env ignorado" || echo "ALERTA — .env podría estar en el repo"
```

---

### 3. Verificar estado del repositorio

```bash
git fetch origin
git status
git diff main origin/main --stat
```

Condiciones requeridas:
- Working tree limpio (sin cambios sin commitear ni archivos sin trackear)
- `main` local sincronizado con `origin/main`

Si el repo no está limpio o main está desactualizado → **detener y reportar**.

---

### 4. Verificar dependencias

```bash
go mod verify
go mod tidy
git diff go.mod go.sum
```

Si `go mod tidy` genera cambios en `go.mod` o `go.sum` → **reportar advertencia** antes de continuar. No commitear cambios de dependencias como parte del issue.

---

### 5. Verificar que el proyecto compila

```bash
go build ./...
```

Si hay errores de compilación → **detener y reportar**. No continuar con código que no compila.

---

### 6. Verificar conectividad con la base de datos

```bash
go run cmd/main.go &
sleep 2
curl -s http://localhost:${PORT:-8080}/health
kill %1 2>/dev/null
```

Respuesta esperada: `{"status":"ok"}`

Si la API no responde o no conecta a la base de datos → **reportar advertencia**. Algunos issues pueden trabajarse sin base de datos (ej: lógica pura de service), pero documentar el estado.

---

### 7. Ejecutar tests existentes y guardar baseline

```bash
mkdir -p .coverage
go test -coverprofile=.coverage/baseline.out ./... 2>&1 | tee /tmp/baseline_tests.txt
go tool cover -func=.coverage/baseline.out | tee /tmp/baseline_coverage.txt
```

Capturar y reportar:
- Cantidad de tests que pasan / fallan
- Coverage por paquete (`config`, `handler`, `service`, `repository`)
- Coverage total

Si hay tests fallando en `main` antes del issue → **detener y reportar**. Los tests rotos existentes no son responsabilidad de este issue.

---

### 8. Generar reporte de baseline

Crear `/tmp/baseline.json`:

```json
{
  "baseline": {
    "timestamp": "<ISO timestamp>",
    "branch": "main",
    "commit": "<hash del último commit>",
    "tests_passed": 0,
    "tests_failed": 0,
    "coverage_total": "0%",
    "coverage_by_package": {
      "internal/handler": "0%",
      "internal/service": "0%",
      "internal/repository": "0%",
      "config": "0%"
    }
  },
  "environment": {
    "go_version": "",
    "sqlc_version": "",
    "database_reachable": true
  },
  "repo_state": "clean"
}
```

Este archivo será consumido por `stage-4-pr.md` para documentar el PR.

---

## Resultado esperado

Reportar al usuario un resumen antes de continuar:

```
✅ Etapa 0 completada

Entorno:          Go 1.23.x
Repo:             main limpio, sincronizado con origin
Dependencias:     OK (go.mod/go.sum sin cambios)
Build:            OK
DB health:        OK
Tests:            87 passing, 0 failing
Coverage total:   94%
  handler:        91%
  service:        97%
  repository:     93%
  config:         100%
Commit base:      a3f9c12

Listo para continuar con /stage-1-analysis.md
```

---

## Regla crítica

Esta etapa **solo detecta y reporta**, nunca corrige. Si algo está mal, el trabajo es del desarrollador, no del agente.
