---
description: Workflow que guía al agente en la resolución completa de un issue de GitHub, desde la asignación y creación del branch hasta la apertura del PR, asegurando cobertura de tests del 90% y trazabilidad de cada paso como comentario en el issue
---

# Workflow — Resolver Issue

## Contexto del proyecto
- Repositorio: worldcups-api
- Lenguaje: Go 1.23
- Framework: Gin
- Arquitectura: Clean Architecture (handler → service → repository)
- Rama principal: main
- Convención de branches: indicada en el campo "branch name:" del issue

---

## Fase 1 — Preparación

1. Leer el issue completo y entender el requerimiento
2. Asignarme el issue
3. Leer el campo "branch name:" del issue
   - Si existe → usar ese nombre
   - Si no existe → solicitarme el nombre antes de continuar
4. Crear el branch con ese nombre desde `main`
5. Hacer checkout al branch creado
6. Cambiar el estado del issue a **En curso**

---

## Fase 2 — Planificación

1. Analizar el issue y elaborar un plan de acción detallado que incluya:
   - Archivos a crear o modificar
   - Capas involucradas (handler / service / repository)
   - Queries SQL nuevas si aplica
   - Tests unitarios necesarios
2. Presentarme el plan y **esperár mi confirmación antes de continuar**
3. Una vez confirmado, agregar el plan como comentario en el issue

---

## Fase 3 — Desarrollo

1. Ejecutar el plan de acción confirmado
2. Seguir las convenciones del proyecto:
   - Nombres de archivos: `{entidad}_{capa}.go`
   - Tests: `{entidad}_{capa}_test.go` en el mismo directorio
   - Errores HTTP con formato: `{"error": "mensaje"}`
   - Prefijo de rutas: `/api/v1/`
3. No modificar archivos fuera del alcance del issue

---

## Fase 4 — Testing

1. Crear los tests unitarios necesarios para alcanzar un **coverage mínimo del 90%**
   - Handler: black-box testing (`package handler_test`)
   - Service: white-box testing (`package service`)
   - Repository: white-box testing (`package repository`)
2. Ejecutar todos los tests:
```bash
   go test ./... -coverprofile=coverage.out
   go tool cover -func=coverage.out
```
3. Si algún test falla → corregir el código o el test hasta que todos pasen
4. Si el coverage es menor al 90% → agregar los tests faltantes
5. Una vez que todos los tests pasen y el coverage sea ≥ 90%, agregar como comentario en el issue el reporte de ejecución como evidencia

---

## Fase 5 — Pull Request

1. Generar un PR desde el branch actual hacia `main` con:
   - **Título:** `[#numero_issue] descripcion_breve`
   - **Descripción:**
     - Resumen del issue resuelto
     - Principales cambios realizados
     - Endpoints nuevos o modificados (si aplica)
     - Link al issue: `Closes #numero_issue`
2. Asignarme como reviewer