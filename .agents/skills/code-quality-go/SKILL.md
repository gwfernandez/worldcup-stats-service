---
name: code-quality-go
description: >
  Auditoría y refactorización de código Go post-desarrollo para asegurar Clean Architecture,
  idiomas de Go, documentación, seguridad y robustez. Soporta múltiples modos de alcance:
  porción de código, archivos individuales, carpetas, combinaciones, proyecto completo o
  archivos pendientes de push. Analiza el código, clasifica los hallazgos por severidad,
  presenta una propuesta con justificación y aplica los cambios SOLO tras confirmación
  explícita del usuario. Ejecuta los tests para verificar que el refactor no rompió nada.

  Usar este skill SIEMPRE que el usuario diga: "revisá el código que escribí", "hacé una
  revisión de calidad", "refactorizá esto", "auditá el código", "mejorar el código",
  "code review", "terminé la implementación, revisalo", o cuando pegue código Go y pida
  feedback o mejoras. Activar también cuando el usuario mencione "clean architecture",
  "Go idioms", "deuda técnica", "mejorar la estructura del código", "revisá lo que no
  pusheé", "auditá mis cambios pendientes", "audit completo" o "revisá todo el proyecto".
---

# Code Quality — Go (Post-Desarrollo)

Skill de auditoría y refactorización para proyectos Go. Transforma código funcional en
código de calidad senior siguiendo el flujo **Alcance → Contexto → Auditoría → Propuesta → Aprobación → Ejecución → Verificación**.

---

## Flujo principal

```
0. ALCANCE    → Determinar qué se va a auditar (porción, archivos, carpeta, proyecto, git diff)
1. CONTEXTO   → Leer AGENTS.md y README.md para conocer las reglas del proyecto
2. AUDITAR    → Analizar en las 5 áreas de revisión
3. PROPONER   → Presentar hallazgos clasificados por severidad
4. APROBAR    → Pedir confirmación explícita antes de tocar nada
5. EJECUTAR   → Aplicar solo los cambios aprobados
6. VERIFICAR  → Correr tests y reportar resultado
```

> **Regla de oro**: nunca modificar un archivo sin confirmación del usuario en el Paso 4.

---

## Paso 0 — Determinar el alcance

Antes de auditar, resolver **qué** se va a revisar. La skill soporta 6 modos de alcance.

> **Todos los paths son relativos a la raíz del proyecto.** El usuario indica archivos y carpetas
> usando rutas como `internal/handler/team_handler.go` o `internal/service/`, nunca paths absolutos
> del filesystem. El agente resuelve la ruta absoluta internamente a partir del workspace.

### Archivos excluidos (siempre)

Independientemente del modo, **excluir siempre** de la auditoría:

- Archivos generados por sqlc: `db.go`, `models.go`, `querier.go`, `*.sql.go` dentro de `internal/repository/`
- Directorio `vendor/`
- Directorio `.agents/`
- Archivos `*_test.go` (salvo que el usuario los pida explícitamente)

### Algoritmo de inferencia del modo

```
¿El usuario pegó código en el chat?
  └── SÍ → Modo Porción

¿El usuario mencionó paths de archivos específicos?
  └── SÍ → ¿Son solo archivos? → Modo Archivo
           ¿Incluye carpetas?  → ¿Incluye archivos también? → Modo Mixto
                                                             → Modo Carpeta

¿El usuario dijo "todo", "completo", "todo el proyecto"?
  └── SÍ → Modo Proyecto (con confirmación obligatoria)

¿El usuario mencionó "sin pushear", "pendientes", "lo que cambié", "git diff"?
  └── SÍ → Modo Git Diff

¿No se puede determinar?
  └── PREGUNTAR: "¿Qué archivos o carpetas querés que revise?"
```

### Modo Porción

**Cuándo**: El usuario pega código directamente en el chat.

- Auditar solo el fragmento recibido
- No requiere confirmación (el alcance está implícito)
- Si el fragmento es menor a 20 líneas, ofrecer feedback inline compacto en lugar del formato completo de auditoría

**Ejemplo**: `"Revisá este handler que escribí"` + bloque de código

### Modo Archivo(s)

**Cuándo**: El usuario indica uno o más archivos por path relativo al proyecto.

- Leer cada archivo con `view_file` y auditar
- No requiere confirmación
- Si un archivo indicado es `*_test.go`, auditar el test

**Ejemplos**:
- `"revisá internal/handler/confederation_handler.go"`
- `"auditá estos archivos: internal/service/confederation_service.go, internal/domain/confederation.go"`

### Modo Carpeta

**Cuándo**: El usuario indica un directorio relativo al proyecto.

- Escanear recursivamente con `list_dir`, recolectar todos los `*.go`
- Aplicar los filtros de exclusión
- Si la carpeta contiene **más de 10 archivos .go**, mostrar la lista y pedir confirmación:

```
📁 internal/service/ contiene 12 archivos Go. ¿Audito todos o preferís seleccionar algunos?
```

**Ejemplos**:
- `"revisá internal/service/"`
- `"auditá toda la capa de handlers"` → se resuelve a `internal/handler/`

### Modo Mixto

**Cuándo**: El usuario combina archivos y carpetas (siempre relativos al proyecto).

- Expandir carpetas + agregar archivos individuales, deduplicar
- Aplicar las reglas de Archivo y Carpeta
- Si el total supera 10 archivos, mostrar lista y pedir confirmación

**Ejemplo**: `"revisá internal/handler/ y también internal/domain/confederation.go"`

### Modo Proyecto completo

**Cuándo**: El usuario pide auditar "todo", "el proyecto completo" o "audit completo".

- Escanear archivos `*.go` en: `cmd/`, `internal/`, `config/`
- Aplicar todos los filtros de exclusión
- **⚠️ Confirmación OBLIGATORIA** antes de empezar

Presentar el resumen antes de iniciar:

```
⚠️ Auditoría de proyecto completo

Voy a auditar todos los archivos Go del proyecto:
  📁 cmd/           → X archivo(s)
  📁 config/        → X archivo(s)
  📁 internal/
     ├── domain/    → X archivo(s)
     ├── handler/   → X archivo(s)
     ├── service/   → X archivo(s)
     └── repository/→ X archivo(s) (excluyendo generados por sqlc)

Total: ~XX archivos Go

¿Confirmo la auditoría completa? [Sí / Solo una capa / Cancelar]
```

> Este modo se alinea con la regla de autonomía de `AGENTS.md`: *"debe pedir confirmación explícita
> antes de realizar cualquier acción de alto impacto"*.

### Modo Git Diff (archivos sin pushear)

**Cuándo**: El usuario quiere revisar solo lo que cambió y no se subió al remote.

**Detección de archivos:**

```bash
# Opción 1: archivos distintos al remote (sin pushear)
git diff --name-only origin/$(git branch --show-current) -- '*.go'

# Fallback si no hay remote tracking branch
git diff --name-only HEAD -- '*.go'

# Archivos staged pero no committeados
git diff --cached --name-only -- '*.go'
```

**Presentación al usuario:**

```
🔀 Archivos Go modificados sin pushear (branch: feature/teams-crud):

  M internal/handler/team_handler.go
  M internal/service/team_service.go
  A internal/domain/team.go

Total: 3 archivos

¿Audito estos 3 archivos? [Sí / Agregar más / Cancelar]
```

**Ejemplo**: `"revisá lo que no pusheé"`, `"auditá mis cambios pendientes"`

---

## Paso 1 — Contexto del proyecto

**Antes de analizar código, el agente DEBE leer los siguientes archivos** (una sola vez por sesión de auditoría, no por cada archivo):

- **`AGENTS.md`** — Reglas de arquitectura, convenciones de código, principios de comportamiento, seguridad y restricciones
- **`README.md`** — Stack tecnológico, estructura del proyecto, formato de errores, comandos de tests

### Qué extraer de AGENTS.md

| Sección | Uso en la auditoría |
|---------|---------------------|
| Arquitectura (Clean Architecture) | Validar flujo handler → service → repository |
| Convenciones de código | Verificar formato de errores `{"error": "..."}`, prefijo `/api/v1/`, nombres de funciones |
| Criterios de "Done" | Validar que el código cumple los requisitos de completitud |
| Principios de comportamiento | Verificar claridad, godoc, manejo explícito de errores |
| Restricciones | No ORMs, no deps sin justificación, interfaces consistentes |
| Seguridad | No loguear datos sensibles, no hardcodear vars de entorno |

### Qué extraer de README.md

| Sección | Uso en la auditoría |
|---------|---------------------|
| Stack tecnológico | Validar que se usan las tecnologías correctas (Gin, pgx, sqlc) |
| Estructura del proyecto | Verificar que archivos nuevos están en el directorio correcto |
| Cómo correr los tests | Usar el comando correcto para verificación |
| Formato de errores | Validar el formato JSON de respuestas de error |

---

## Paso 2 — Auditoría en 5 áreas

Analizar el código en este orden. Ver `references/go-clean-architecture.md` para
patrones y antipatrones concretos de cada área.

### Área 1 — Arquitectura y Dependencias

Verificar el flujo de dependencias: `Handler → Service → Repository`

**Detectar:**
- Lógica de base de datos (queries, ORM) en el Service layer
- Lógica de negocio (validaciones, cálculos) en el Handler
- Imports directos de infraestructura (DB, HTTP clients) desde capas internas
- Structs concretos donde deberían ir interfaces

**Preguntas clave:**
- ¿El Handler solo orquesta la request/response?
- ¿El Service solo contiene reglas de negocio?
- ¿El Repository solo habla con la base de datos?

### Área 2 — Manejo de Errores

**Detectar:**
- Errores ignorados: `result, _ := algo()`
- Errores genéricos sin contexto: `return err` sin wrapping
- Panic innecesario donde debería haber return de error
- Errores de negocio mezclados con errores de infraestructura
- Respuestas HTTP sin el formato `{"error": "mensaje"}`

**Patrón correcto:**
```go
// ✅ Error con contexto
if err != nil {
    return fmt.Errorf("userService.GetByID: %w", err)
}

// ✅ Respuesta de error en handler
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
```

### Área 3 — Calidad de Código Go

**Detectar:**
- Complejidad ciclomática alta: funciones con más de **3 niveles de anidamiento** o más de **5 branches** (if/switch/for combinados) → candidatas a dividir
- Nombres no idiomáticos: variables de un solo carácter fuera de loops, nombres con tipos redundantes (`userList` → `users`), abreviaciones poco claras
- Interfaces demasiado grandes (más de 5 métodos) → candidatas a dividir
- Goroutines sin manejo de cierre o cancelación
- Context no propagado a través de las capas

**Umbrales de complejidad:**
| Complejidad | Acción |
|-------------|--------|
| 1–5 | OK, sin acción |
| 6–10 | Sugerencia de división |
| > 10 | Hallazgo crítico, división necesaria |

### Área 4 — Documentación

**Detectar:**
- Funciones, structs o interfaces **públicas** sin comentario Godoc
- Comentarios que repiten el nombre en lugar de explicar el propósito
- Falta de ejemplo en funciones complejas o públicas de librerías

**Formato correcto:**
```go
// UserService maneja la lógica de negocio relacionada con usuarios.
// Depende de UserRepository para el acceso a datos.
type UserService struct { ... }

// GetByID busca un usuario por su identificador único.
// Retorna ErrUserNotFound si el usuario no existe.
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) { ... }
```

### Área 5 — Seguridad

Basado en las reglas de seguridad definidas en `AGENTS.md`.

**Detectar:**
- Variables de entorno hardcodeadas en el código fuente (connection strings, puertos, API keys)
- Datos sensibles en logs: passwords, tokens, `DATABASE_URL` completa, información personal
- Credenciales o secretos en el código fuente
- Información interna expuesta en responses de error al cliente (stack traces, queries SQL, paths del filesystem)
- Uso de `log.Fatal` o `panic` donde debería haber manejo controlado de errores

**Patrón correcto:**
```go
// ✅ Log seguro — solo IDs, nunca credenciales
log.Printf("error al buscar confederación id=%d: %v", id, err)

// ✅ Error al cliente — mensaje genérico, sin detalles internos
c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})

// ❌ MAL — expone detalles de infraestructura
c.JSON(500, gin.H{"error": err.Error()})  // puede exponer "pq: connection refused to neon.tech..."

// ❌ MAL — variable hardcodeada
db, _ := pgx.Connect(ctx, "postgresql://user:pass@host/db")
```

---

## Paso 3 — Presentar la propuesta

Clasificar cada hallazgo con una severidad y presentarlos agrupados.

### Niveles de severidad

| Nivel | Emoji | Cuándo usarlo |
|-------|-------|---------------|
| 🔴 Crítico | Bug, error ignorado, lógica en capa incorrecta, dato sensible expuesto | Rompe correctitud, arquitectura o seguridad |
| 🟠 Importante | Falta de interfaces, complejidad alta, error sin contexto | Afecta mantenibilidad o testabilidad |
| 🟡 Mejora | Naming, documentación, simplificaciones | Buenas prácticas, no urgente |
| 🔵 Info | Observaciones sin acción requerida | Contexto útil para el desarrollador |

### Formato de presentación (archivo individual)

```
## 🔍 Auditoría de Código — [nombre del archivo/módulo]

### 🔴 Críticos (2)
**1. Error ignorado en GetUser**
- Línea: 34
- Problema: `user, _ := repo.GetByID(id)` — el error se descarta silenciosamente
- Impacto: si el repositorio falla, `user` será nil y causará un panic en la línea 36
- Corrección propuesta: propagar el error y retornar respuesta 500

**2. Lógica de negocio en el Handler**
[...]

### 🟠 Importantes (1)
[...]

### 🟡 Mejoras (3)
[...]

---
Total: 2 críticos · 1 importante · 3 mejoras

¿Querés que aplique todas las correcciones automáticamente, o preferís elegir cuáles?
[Aplicar todo] | [Elegir por severidad] | [Revisar una por una]
```

### Reporte consolidado (múltiples archivos)

Cuando la auditoría involucra más de un archivo (modos Carpeta, Mixto, Proyecto o Git Diff),
agregar un **reporte resumen** al final de todos los hallazgos individuales:

```
## 📊 Resumen de auditoría

| Archivo | 🔴 | 🟠 | 🟡 | 🔵 |
|---------|----|----|----|----|
| team_handler.go | 1 | 0 | 2 | 1 |
| team_service.go | 0 | 1 | 1 | 0 |
| team.go         | 0 | 0 | 0 | 1 |

Total: 1 crítico · 1 importante · 3 mejoras · 2 info

¿Querés que aplique todas las correcciones, o preferís elegir por archivo o severidad?
[Aplicar todo] | [Elegir por archivo] | [Elegir por severidad] | [Revisar una por una]
```

---

## Paso 4 — Solicitar aprobación

**Siempre preguntar antes de modificar cualquier archivo.**

Opciones a ofrecer:
- **Aplicar todo**: aplica todos los hallazgos en una pasada
- **Por severidad**: "¿Aplico solo los críticos e importantes?"
- **Por archivo** (modo multi-archivo): "¿En cuáles archivos aplico?"
- **Uno por uno**: recorrer cada hallazgo y confirmar individualmente

Si el usuario modifica o rechaza algún hallazgo, documentar la decisión y no insistir.

---

## Paso 5 — Ejecutar los cambios

Para cada cambio aprobado:
1. Leer el archivo actual con `view_file`
2. Aplicar el cambio de forma quirúrgica (no reescribir el archivo entero)
3. Confirmar el cambio aplicado mostrando el diff conceptual

Nunca reescribir un archivo completo si el cambio es pequeño y localizado.

---

## Paso 6 — Verificar con tests

Después de aplicar todos los cambios:

```bash
go test ./...
```

### Interpretar el resultado

| Resultado | Acción |
|-----------|--------|
| ✅ Todos pasan | Reportar éxito, listar tests ejecutados |
| ⚠️ Fallos pre-existentes | Indicar que los fallos no son del refactor (mostrar cuáles) |
| ❌ Nuevos fallos | **Revertir el cambio causante**, reportar al usuario y analizar la causa |

Si los tests fallan por el refactor, revertir el cambio específico y reportar qué
fue lo que rompió para que el usuario decida cómo proceder.

### Post-verificación

Una vez que todos los cambios están aplicados y los tests pasan, sugerir al usuario:

> "Los cambios fueron aplicados y los tests pasan. Podés usar la skill `semantic-commit`
> para registrar estas correcciones en un commit."

---

## Escenarios de uso recomendados

### Cuándo activar esta skill

| Escenario | Modo sugerido | Ejemplo de activación |
|-----------|---------------|----------------------|
| Terminé de implementar una feature | Git Diff | `"revisá lo que escribí"` |
| Antes de hacer un commit | Git Diff | `"quiero commitear, primero revisá"` |
| Después de resolver un issue | Archivo/Carpeta | Los archivos tocados en el issue |
| Revisión de una capa completa | Carpeta | `"revisá la capa de service"` |
| Duda sobre un fragmento específico | Porción | Código pegado + pregunta |
| Auditoría periódica de calidad | Proyecto | `"hacé un audit de todo"` |
| Code review de un PR | Git Diff | `"revisá los cambios del PR"` |
| Refactoring post-merge | Archivo/Mixto | `"refactorizá estos archivos"` |

### Formas simples de activar la skill

```
• "Revisá este código"          → pegás el código y listo (Porción)
• "Auditá team_handler.go"     → le das el archivo (Archivo)
• "Revisá internal/service/"   → le das la carpeta (Carpeta)
• "Revisá lo que no pusheé"    → audita cambios pendientes (Git Diff)
• "Audit completo"             → todo el proyecto con confirmación (Proyecto)
```

---

## Integración con workflows y skills

### Con workflow `resolve-issue`

Esta skill está integrada como **Fase 3.5 obligatoria** en el workflow `resolve-issue`.
Se ejecuta en modo Git Diff entre el desarrollo y el testing:

```
Fase 3   — Desarrollo
  ↓
Fase 3.5 — code-quality-go (modo Git Diff — archivos del issue) ← obligatorio
  ↓
Fase 4   — Testing
```

Los hallazgos 🔴 Críticos deben corregirse antes de avanzar a testing.

### Con skill `semantic-commit`

Flujo natural recomendado:

```
Escribir código → code-quality-go → Corregir hallazgos → semantic-commit
```

Una vez aplicados los cambios y verificados los tests, sugerir al usuario que use
`semantic-commit` para registrar las correcciones.

### Con skill `sdd-requirements`

No hay integración directa. Cuando se crea un issue con `sdd-requirements`, el trabajo
resultante de resolverlo (vía `resolve-issue`) debería pasar por `code-quality-go` antes del PR.

---

## Referencias

- `references/go-clean-architecture.md` — Patrones y antipatrones de Clean Architecture en Go con ejemplos de código
- `references/go-error-patterns.md` — Guía de manejo de errores: wrapping, tipos de error, sentinel errors, errores de negocio vs infraestructura
- `AGENTS.md` — Reglas de arquitectura, convenciones, seguridad y restricciones del proyecto
- `README.md` — Stack tecnológico, estructura, formato de errores y comandos de tests