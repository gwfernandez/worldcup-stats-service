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

1. Analizar el issue de GitHub y elaborar un plan de acción detallado que incluya:
   - Archivos a crear o modificar
   - Capas involucradas (handler / service / repository)
   - Queries SQL nuevas si aplica
   - Tests unitarios necesarios
2. Presentarme el plan (Implementation Plan) y **esperár mi confirmación antes de continuar**
3. Una vez que apruebe el plan, **DEBES usar la herramienta de GitHub para publicar el "Implementation Plan" completo como un comentario en el issue original**.

---

## Fase 3 — Desarrollo

1. Ejecutar el plan de acción confirmado
2. Seguir las convenciones del proyecto:
   - Nombres de archivos: `{entidad}_{capa}.go`
   - Tests: `{entidad}_{capa}_test.go` en el mismo directorio
   - Errores HTTP con formato: `{"error": "mensaje"}`
   - Prefijo de rutas: `/api/v1/`
3. No modificar archivos fuera del alcance del issue de GitHub
4. Si el issue contiene una lista de tareas (tasks), realizar un commit independiente utilizando la skill `semantic-commit` por cada tarea completada. El mensaje del commit debe reflejar fielmente la tarea realizada. En caso de no haber una lista, realizar commits por cada hito lógico finalizado.

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
5. Una vez que todos los tests pasen y el coverage sea ≥ 90%, **usar la integración de GitHub para agregar como comentario en el issue original** el reporte de ejecución como evidencia

---

## Fase 5 — Pull Request

1. Usar la herramienta de creación de Pull Requests de GitHub para generar un PR desde el branch actual hacia `main` con:
   - **Título:** `tipo(scope): descripción breve (#numero_issue)`
   - **Descripción:**
     - Resumen del issue de GitHub resuelto
     - Principales cambios realizados
      - Endpoints nuevos o modificados (si aplica)
      - Impacto SemVer estimado (MAJOR/MINOR/PATCH)
      - Link al issue: `Closes #numero_issue`
2. Asignarme como reviewer (si la API lo permite, si no, dejar documentado)
3. **Usar la herramienta de comentarios de GitHub** para publicar el documento de resumen ("Walkthrough") como un comentario final en el issue original, indicando que el trabajo ha concluido.