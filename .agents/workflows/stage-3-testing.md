---
description: Stage 3 — Testing y coverage
---

# Stage 3 — Testing y coverage

## Propósito

Ejecutar la suite completa de tests, verificar que el coverage mínimo del 90% se mantiene en los paquetes modificados y validar que los tests son significativos: cubren los edge cases del issue, no solo el happy path.

---

## Pasos

### 1. Ejecutar la suite completa

```bash
go test ./... 2>&1 | tee /tmp/stage3_tests.txt
```

Comparar con el baseline de `/tmp/baseline_tests.txt` (generado en Stage 0):

- Los tests que pasaban antes deben seguir pasando
- No deben aparecer regresiones en paquetes no relacionados con el issue

Si hay regresiones en código preexistente → **investigar y corregir antes de continuar**.

---

### 2. Verificar coverage por paquete

```bash
mkdir -p .coverage
go test -coverprofile=.coverage/coverage.out ./...
go tool cover -func=.coverage/coverage.out | tee /tmp/stage3_coverage.txt
```

**Coverage mínimo requerido: 90% por paquete modificado**

Revisar específicamente los paquetes tocados en este issue:

```bash
go tool cover -func=.coverage/coverage.out | grep -E "internal/(handler|service|repository)|config"
```

Paquetes excluidos del cálculo (código generado o entrypoint):
- `cmd/` — entrypoint, sin lógica testeable
- `db/sqlc/` o el paquete generado por sqlc — código auto-generado

Si algún paquete modificado tiene coverage < 90% → identificar las funciones sin cobertura y agregar tests antes de continuar:

```bash
go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
# Revisar visualmente qué ramas no están cubiertas
```

---

### 3. Revisar calidad de los tests

#### 3.1 Cobertura de edge cases del issue

Confrontar la lista de edge cases identificados en el **Checkpoint A del Stage 1** con los tests escritos.

Por cada edge case, debe existir al menos un test que lo cubra:

- [ ] Happy path — respuesta exitosa con datos
- [ ] Recurso no encontrado — debe retornar HTTP 404 con `{"error": "..."}`
- [ ] Parámetros inválidos — debe retornar HTTP 400 con `{"error": "..."}`
- [ ] Error de base de datos — debe retornar HTTP 500 con `{"error": "..."}`
- [ ] Edge cases específicos del issue (lista vacía, filtros, paginación, etc.)

#### 3.2 Tests de handlers — verificar que usan httptest correctamente

```go
// ✅ Correcto — usa httptest, mocka el service, verifica status y body JSON
func TestListScorers_OK(t *testing.T) {
    mockService := new(MockScorerService)
    mockService.On("ListScorers", mock.Anything, 2022).Return(mockScorers, nil)

    router := gin.New()
    h := NewScorerHandler(mockService)
    router.GET("/api/scorers", h.List)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/scorers?year=2022", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    // verificar body JSON
}

// ✅ Debe existir también el caso de error
func TestListScorers_DBError(t *testing.T) {
    mockService.On("ListScorers", mock.Anything, 2022).Return(nil, errors.New("db error"))
    // debe retornar 500
}
```

#### 3.3 Tests de services — verificar aislamiento del repository

```go
// ✅ Correcto — mockea el repository, testea solo la lógica del service
func TestScorerService_ListScorers_NotFound(t *testing.T) {
    mockRepo := new(MockScorerRepository)
    mockRepo.On("GetScorersByYear", mock.Anything, 2022).Return(nil, ErrNotFound)

    svc := NewScorerService(mockRepo)
    result, err := svc.ListScorers(context.Background(), 2022)

    assert.Nil(t, result)
    assert.ErrorIs(t, err, ErrNotFound)
}
```

#### 3.4 Tests de repositories — verificar uso de pgxmock

```go
// ✅ Correcto — usa pgxmock para simular la DB sin levantar PostgreSQL
func TestScorerRepository_GetScorersByYear(t *testing.T) {
    db, mock, _ := pgxmock.NewPool()
    mock.ExpectQuery("SELECT").
        WithArgs(2022).
        WillReturnRows(pgxmock.NewRows(cols).AddRow(...))

    repo := NewScorerRepository(db)
    result, err := repo.GetScorersByYear(context.Background(), 2022)

    assert.NoError(t, err)
    assert.Len(t, result, 1)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

#### 3.5 Detectar tests triviales

```go
// ❌ Test trivial — solo verifica que el mock devuelve lo que fue programado
func TestListScorers(t *testing.T) {
    mockService.On("ListScorers", ...).Return(mockData, nil)
    result, _ := mockService.ListScorers(ctx, 2022)
    assert.Equal(t, mockData, result)  // siempre pasa, no testea nada
}

// ✅ Test útil — testea comportamiento real del service
func TestScorerService_FiltraGoleadoresPorAnio(t *testing.T) {
    // testea que el service llama al repo con el año correcto
    // y transforma el resultado como se espera
}
```

Si se detectan tests triviales → reemplazarlos por tests con comportamiento real.

---

### 4. Verificar que los tests no dependen de orden de ejecución

```bash
go test -shuffle=on ./... 2>&1 | tail -20
```

Si hay tests que fallan con orden aleatorio → corregir el aislamiento (setup/teardown, estado global compartido).

---

### 5. Verificar que los tests no requieren base de datos real

```bash
# Ningún test debe necesitar DATABASE_URL real — todo debe ir por mocks
grep -rn "DATABASE_URL\|os.Getenv.*DATABASE" internal/ | grep "_test.go"
```

Si hay tests que conectan a una base de datos real → reemplazar con `pgxmock`.

---

### 6. Reporte final de coverage

```bash
go tool cover -func=.coverage/coverage.out | grep "total:"
```

Comparar con el baseline del Stage 0.

---

### 7. Publicar reporte de coverage en el issue de GitHub

Una vez completado el análisis de coverage, publicar el siguiente comentario en el issue usando
la integración de GitHub:

```
✅ Testing completado — Issue #<número>

Tests totales:     X passing, 0 failing (+Y nuevos)
Regresiones:       ninguna

Coverage por paquete:
  internal/handler:     X%  (antes: X%)
  internal/service:     X%  (antes: X%)
  internal/repository:  X%  (antes: X%)
  config:               X%  (antes: X%)
  TOTAL:                X%  (antes: X%)

Cobertura de edge cases:
  ✅ <edge case 1>
  ✅ <edge case 2>
  ...
```

> Este comentario sirve como evidencia de testing para el reviewer del PR.

---

### 8. Presentar resultado — CHECKPOINT

```
✅ Testing completado

Tests totales:     105 passing, 0 failing  (+18 nuevos)
Regresiones:       ninguna

Coverage por paquete:
  internal/handler:     92%  (baseline: 91%)
  internal/service:     96%  (baseline: 97%)
  internal/repository:  91%  (baseline: 93%)
  config:               100% (baseline: 100%)
  TOTAL:                94%  (baseline: 94%)

Cobertura de edge cases del issue #42:
  ✅ Happy path: lista de goleadores con datos
  ✅ Lista vacía: mundial sin goleadores cargados
  ✅ Año inválido: retorna 400
  ✅ Año sin mundial: retorna 404
  ✅ Error de base de datos: retorna 500

Calidad de tests:
  ✅ Handlers usan httptest correctamente
  ✅ Services mockean el repository
  ✅ Repositories usan pgxmock (sin DB real)
  ✅ Sin tests triviales detectados
  ✅ Sin dependencias de orden entre tests

Listo para continuar con /stage-4-pr.md
```

> ⏸️ **Esperar aprobación del usuario antes de continuar.**
