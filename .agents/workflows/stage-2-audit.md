---
description: Stage 2 — Auditoría y corrección
---

# Stage 2 — Auditoría y corrección

## Propósito

Revisar el código generado en el Stage 1 contra los estándares del proyecto usando la skill
`code-quality-go`, y corregir todos los hallazgos críticos e importantes antes de avanzar al testing.

Este stage **no duplica lógica de auditoría** — delega completamente en la skill, que es la
fuente de verdad para los estándares de calidad del proyecto (Clean Architecture, Go idioms,
manejo de errores, documentación, seguridad).

---

## Prerrequisito

La skill `code-quality-go` debe estar disponible en `.agents/skills/code-quality-go/SKILL.md`.
Si no está presente → avisar al usuario y detener.

---

## Pasos

### 1. Verificación estática previa (condición de entrada)

Antes de invocar la skill, correr los checks que son condición de entrada:

```bash
go build ./...
go vet ./...
```

Corregir **todos** los errores reportados. La skill audita código que compila, no errores de
compilación básicos. No continuar si alguno falla.

---

### 2. Invocar la skill `code-quality-go` en Modo Git Diff

Leer el archivo [SKILL.md](../skills/code-quality-go/SKILL.md) y activar la skill en
**Modo Git Diff**, pasando el contexto del issue como referencia:

```
Activar skill code-quality-go.

Modo: Git Diff
Alcance: archivos .go modificados en esta rama respecto a main (cambios del issue #<número>)
Contexto: issue #<número> — <título del issue>
```

La skill se encarga de:
- Detectar automáticamente los archivos `.go` modificados sin pushear
- Leer `AGENTS.md` y `README.md` para conocer las reglas del proyecto
- Auditar en 5 áreas: Arquitectura, Manejo de errores, Calidad Go, Documentación, Seguridad
- Presentar hallazgos clasificados por severidad (🔴 Crítico / 🟠 Importante / 🟡 Mejora / 🔵 Info)

---

### 3. Resolver la propuesta de la skill

Cuando la skill presente los hallazgos y ofrezca opciones de aplicación, responder con la
siguiente política **no negociable en el contexto de este workflow**:

| Severidad | Política |
|-----------|----------|
| 🔴 Crítico | **Corregir obligatoriamente** antes de avanzar al Stage 3 |
| 🟠 Importante | Corregir salvo decisión explícita del usuario de postergarlos |
| 🟡 Mejora | Aplicar si no representa riesgo de regresión; postergar si hay dudas |
| 🔵 Info | No requieren acción; documentar si son relevantes para el PR |

Opción recomendada a seleccionar: **"Aplicar todo"** o **"Elegir por severidad"** priorizando
🔴 y 🟠.

> Si el usuario rechaza explícitamente algún hallazgo 🔴 Crítico → **detener y documentar
> la decisión** antes de continuar. No avanzar con código que tiene hallazgos críticos sin
> justificación explícita.

---

### 4. Verificación post-corrección

La skill ejecuta `go test ./...` automáticamente al finalizar (Paso 6 de la skill). Esperar
el resultado antes de continuar:

| Resultado | Acción |
|-----------|--------|
| ✅ Todos los tests pasan | Continuar al checkpoint |
| ⚠️ Fallos pre-existentes | Documentar y continuar (no son responsabilidad de este stage) |
| ❌ Nuevos fallos por el refactor | La skill revierte el cambio causante — analizar y decidir con el usuario |

---

### 5. Registrar las correcciones (si aplica)

Si la skill aplicó correcciones, registrarlas con la skill `semantic-commit`:

```
Activar skill semantic-commit.

Contexto: correcciones de auditoría de calidad sobre los archivos del issue #<número>
Tipo esperado: refactor o fix según la naturaleza de las correcciones
```

> Si no hubo hallazgos o ninguno requirió corrección → no commitear; continuar directamente
> al Stage 3.

---

### 6. Presentar resultado — CHECKPOINT

```
✅ Stage 2 — Auditoría completada

Herramienta:     skill code-quality-go (Modo Git Diff)
Archivos auditados: X archivos .go

Hallazgos:
  🔴 Críticos:    0 (o N — todos corregidos)
  🟠 Importantes: 0 (o N — corregidos / postergados con justificación)
  🟡 Mejoras:     N aplicadas / N postergadas
  🔵 Info:        N (sin acción)

Correcciones aplicadas:
  - [archivo]: [descripción breve del hallazgo corregido]

Tests post-corrección: X passing, 0 failing

Commit de correcciones: [hash] (si aplica) / No requerido

Listo para continuar con /stage-3-testing.md
```

Si quedaron hallazgos 🔴 sin resolver → listarlos explícitamente con la justificación del
usuario antes de continuar.

> ⏸️ **Esperar aprobación del usuario antes de continuar.**
