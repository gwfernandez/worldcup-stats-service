---
name: semantic-commit
description: >
  Redactar mensajes de commit que cumplan con Conventional Commits y sean compatibles con
  go-semantic-release para versionado semántico automático (SemVer). Analiza los cambios del
  usuario, determina el tipo correcto (feat/fix/chore/etc.), el impacto en versión (MAJOR/MINOR/PATCH),
  genera el mensaje formateado y ejecuta el commit con confirmación.

  Usar este skill SIEMPRE que el usuario mencione: "hacer un commit", "commitear", "escribir un
  commit message", "preparar un commit", "git commit", "quiero versionar", "cambios para commitear",
  o cuando describa cambios en el código y quiera registrarlos en git. Activar también cuando
  el usuario diga "rompí algo", "agregué una feature", "arreglé un bug", "refactoricé X",
  "actualicé dependencias", o cualquier descripción de cambio en el código que implique
  querer hacer commit.
---

# Semantic Commit Skill

Skill para generar mensajes de commit que cumplan con la especificación **Conventional Commits**
compatible con **go-semantic-release**, determinando automáticamente el impacto en versión SemVer
y ejecutando el commit con confirmación del usuario.

---

## Flujo principal

```
1. RECIBIR   → Descripción libre de los cambios realizados
2. ANALIZAR  → Determinar tipo, scope, breaking change e impacto en versión
3. COMPLETAR → Pedir datos faltantes si es necesario (máx. 1 ronda)
4. GENERAR   → Mensaje de commit formateado según la spec
5. MOSTRAR   → Preview con impacto en versión explicado
6. CONFIRMAR → Pedir aprobación y ejecutar git commit
```

---

## Paso 1 — Recibir la descripción

Aceptar cualquier descripción libre. El usuario puede escribir:

> "agregué validación de email en el formulario de registro"
> "rompí la API de pagos, cambiaron los parámetros del endpoint"
> "moví funciones helper a un archivo utils.js"

No pedir información extra hasta haber intentado inferirla.

---

## Paso 2 — Analizar y determinar campos

### Estructura del mensaje (formato Conventional Commits)

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Tipos y su impacto en versión

| Tipo | Impacto SemVer | Cuándo usarlo |
|------|---------------|---------------|
| `feat` | **MINOR** (0.X.0) | Nueva funcionalidad visible al usuario |
| `fix` | **PATCH** (0.0.X) | Corrección de bug |
| `perf` | **PATCH** | Mejora de performance sin cambio de API |
| `refactor` | Sin release* | Refactoring sin cambio de comportamiento |
| `docs` | Sin release* | Solo documentación |
| `style` | Sin release* | Formato, espacios, punto y coma |
| `test` | Sin release* | Agregar o corregir tests |
| `build` | Sin release* | Sistema de build, dependencias |
| `ci` | Sin release* | Configuración de CI/CD |
| `chore` | Sin release* | Tareas de mantenimiento |
| `revert` | **PATCH** | Revertir un commit anterior |

> *Sin release por defecto en go-semantic-release. Se puede configurar en `.semrelrc`.

### Breaking Change → MAJOR (X.0.0)

Un commit de **cualquier tipo** se convierte en MAJOR si incluye:

```
feat(api)!: cambiar formato de respuesta

BREAKING CHANGE: el campo `userId` fue renombrado a `user_id`
```

Indicadores de breaking change en la descripción del usuario:
- "rompí", "cambié la firma", "eliminé", "renombré", "cambiaron los parámetros"
- "ya no es compatible", "hay que migrar", "cambió la API", "breaking"

### Inferencia de scope

El scope es opcional pero recomendado. Inferirlo del contexto:
- Si menciona un módulo, servicio o archivo específico → usarlo como scope
- Si el cambio es transversal → omitir scope
- Si hay duda → preguntar solo si agrega valor real

---

## Paso 3 — Completar información faltante

Hacer preguntas **solo si son necesarias** para escribir el commit correctamente.
Máximo **una ronda** de preguntas. Agrupar todo en un solo mensaje.

**Preguntar si:**
- No está claro si es breaking change (consecuencias importantes)
- No está claro si es `feat` o `fix` (impacto diferente en versión)
- El scope agregaría claridad significativa

**No preguntar si:**
- El tipo es obvio por contexto
- El scope no agrega valor
- La descripción es suficientemente clara

---

## Paso 4 — Generar el mensaje

### Reglas de formato

**Header** (obligatorio):
- Máximo 72 caracteres
- Tipo en minúscula
- Descripción: verbo en imperativo, minúscula, sin punto final
- `!` antes del `:` si es breaking change

**Body** (incluir cuando el cambio no es obvio):
- Separado del header por línea en blanco
- Explicar el *qué* y el *por qué*, no el *cómo*
- Máximo 100 caracteres por línea
- En el mismo idioma que el proyecto (inferir del contexto)

**Footer** (obligatorio si hay breaking change):
```
BREAKING CHANGE: descripción del cambio incompatible
```

También válido para referencias:
```
Closes #123
Refs #456
```

### Ejemplos de commits bien formados

```bash
# PATCH — fix simple
fix(auth): corregir validación de token expirado

# MINOR — nueva feature con scope
feat(user): agregar endpoint para actualizar avatar

# MINOR — nueva feature con body explicativo
feat(reportes): exportar datos en formato CSV

Permite descargar el historial completo de transacciones.
El archivo incluye encabezados y respeta el encoding UTF-8.

Closes #89

# MAJOR — breaking change con !
feat(api)!: cambiar estructura de respuesta en /users

BREAKING CHANGE: el campo `data.userId` fue renombrado a `data.user_id`
para alinearse con la convención snake_case del resto de la API.
Los clientes deben actualizar sus integraciones.

# Sin release — refactor
refactor(utils): extraer funciones de fecha a módulo helpers

# Sin release — docs
docs(readme): agregar instrucciones de instalación en Docker
```

---

## Paso 5 — Mostrar preview con impacto en versión

Mostrar siempre el impacto antes de ejecutar:

```
📝 Commit generado:

  feat(auth): agregar login con Google OAuth

  Permite a los usuarios autenticarse usando su cuenta de Google.
  Requiere configurar las variables GOOGLE_CLIENT_ID y GOOGLE_SECRET.

  Closes #112

📦 Impacto en versión (go-semantic-release):
  Versión actual: v1.3.2  →  Nueva versión: v1.4.0  (MINOR)
  Motivo: `feat` incrementa el MINOR

¿Confirmás este commit? [Sí / Modificar / Cancelar]
```

Si no se conoce la versión actual, mostrar solo el tipo de bump:
```
📦 Impacto: MINOR bump (feat → incrementa X.Y.0)
```

---

## Paso 6 — Confirmar y ejecutar

### Si hay terminal disponible (`bash_tool`)

```bash
git add -A   # o el staging que corresponda
git commit -m "<header>" -m "<body>" -m "<footer>"
```

Mostrar el output de git después de ejecutar.

Si el usuario quiere stagear archivos específicos, preguntar antes de `git add`.

### Si NO hay terminal disponible

Mostrar el comando listo para copiar:

```bash
git commit -m "feat(auth): agregar login con Google OAuth" \
  -m "Permite a los usuarios autenticarse usando su cuenta de Google.
Requiere configurar las variables GOOGLE_CLIENT_ID y GOOGLE_SECRET." \
  -m "Closes #112"
```

---

## Validaciones antes de generar

Verificar antes de mostrar el commit:

- [ ] Header ≤ 72 caracteres
- [ ] Tipo es uno de los valores válidos
- [ ] Descripción en imperativo y minúscula
- [ ] Sin punto final en el header
- [ ] Footer `BREAKING CHANGE:` presente si hay `!` en el tipo
- [ ] Body separado del header por línea en blanco

---

## Casos especiales

### Múltiples cambios independientes → múltiples commits
Si el usuario describe cambios de tipos distintos (ej: fix + feat), sugerir dividirlos:
```
Detecté dos cambios independientes. Te recomiendo hacer commits separados:
1. fix(login): corregir error en validación de contraseña
2. feat(perfil): agregar foto de perfil
¿Querés commitearlos por separado?
```

### Revert
```bash
revert: revert "feat(user): agregar avatar"

This reverts commit abc1234.
```

### WIP / trabajo incompleto
Sugerir usar `chore` o `wip` con aclaración, recordando que no genera release:
```bash
chore(wip): avance en módulo de pagos [no release]
```

---

## Referencias

- `references/conventional-commits-spec.md` — Spec completa con ejemplos edge cases
- `references/semver-impact-table.md` — Tabla completa de tipos → impacto en versión con configuración go-semantic-release