---
name: sdd-requirements
description: >
  Redacción y estructuración de requerimientos de software con formato SDD (Software Design Document),
  incluyendo análisis iterativo del requerimiento, generación de documento Markdown completo, propuesta
  de tareas/subtareas, y creación del issue en GitHub con confirmación del usuario.
  
  Usar este skill SIEMPRE que el usuario mencione: requerimientos, issues, historias de usuario,
  SDD, tareas de desarrollo, epics, feature requests, bug reports estructurados, o cuando pida
  "crear un issue", "documentar un requerimiento", "redactar una historia de usuario", o cuando
  quiera registrar algo en GitHub como trabajo de desarrollo. Activar también cuando el usuario
  diga "necesito implementar X", "quiero agregar Y al sistema" o "hay un problema con Z".
---

# SDD Requirements Skill

Skill para analizar requerimientos en lenguaje natural, completar información faltante de forma iterativa, generar documentación SDD estructurada en Markdown y crear el issue en GitHub con confirmación del usuario.

---

## Flujo principal

```
1. RECIBIR   → Texto libre del usuario con el requerimiento
2. ANALIZAR  → Extraer campos SDD, detectar vacíos
3. ITERAR    → Solicitar información faltante (máx. 3 rondas)
4. GENERAR   → Documento Markdown con estructura SDD completa
5. CONFIRMAR → Mostrar preview y pedir aprobación
6. CREAR     → Issue en GitHub vía MCP o instrucciones manuales
```

---

## Paso 1 — Recibir el requerimiento

Aceptar texto libre. No pedir formularios al inicio. El usuario puede escribir tan poco como:

> "Necesito un botón para exportar reportes en PDF"

Eso es suficiente para comenzar el análisis.

---

## Paso 2 — Analizar y detectar vacíos

Extraer del texto todos los campos posibles del esquema SDD (ver `references/sdd-schema.md`).

**Campos OBLIGATORIOS** — sin ellos no se puede generar el issue:
- `title` — Título conciso del requerimiento
- `type` — Feature / Bug / Improvement / Task / Epic
- `description` — Qué se necesita y por qué
- `acceptance_criteria` — Al menos 1 criterio verificable

**Campos IMPORTANTES** — pedir si no están implícitos:
- `user_story` — Como [rol], quiero [acción], para [beneficio]
- `affected_components` — Módulos, servicios o pantallas afectadas
- `priority` — Critical / High / Medium / Low

**Campos OPCIONALES** — inferir o dejar vacíos:
- `assumptions`, `out_of_scope`, `dependencies`, `mockups_refs`

### Reglas de inferencia
- Si el texto menciona un rol de usuario → construir user story
- Si hay una acción clara → proponer acceptance criteria básicos
- Si el componente es obvio del contexto → no preguntar
- Si el tipo es ambiguo → asumir Feature y aclarar al usuario

---

## Paso 3 — Iteración para completar información

Si faltan campos obligatorios o importantes, hacer preguntas **agrupadas y específicas**.

### Reglas de iteración
- Máximo **3 rondas** de preguntas
- Agrupar todas las preguntas faltantes en UNA sola ronda, no una por vez
- Ser concreto: ofrecer opciones cuando sea posible
- Si tras 3 rondas aún falta algo no-crítico, documentarlo como `[POR DEFINIR]` y continuar
- Si falta algo crítico tras 3 rondas, generar el documento con los mejores supuestos y marcarlo explícitamente

### Ejemplo de pregunta agrupada bien formulada
```
Para completar el requerimiento necesito algunos datos:

1. **Tipo**: ¿Es una nueva funcionalidad (Feature), una mejora (Improvement) o un bug?
2. **Usuario objetivo**: ¿Quién usará esta función? (ej: admin, cliente final, operador)
3. **Criterio de éxito**: ¿Cómo sabremos que está correctamente implementado?
   - ¿El PDF debe generarse en menos de X segundos?
   - ¿Debe incluir algún formato o logo específico?
```

---

## Paso 4 — Generar documento Markdown SDD

Una vez completa la información mínima, generar el documento usando la plantilla en `references/sdd-template.md`.

El documento debe:
- Estar en el idioma en que el usuario escribió el requerimiento
- Incluir propuesta de tareas técnicas (ver sección Tasks)
- Usar etiquetas/labels de GitHub apropiadas
- Tener estimación de complejidad (XS / S / M / L / XL)

### Propuesta de Tareas

Generar una lista de subtareas técnicas desglosadas:

```markdown
## 📋 Tasks

- [ ] [BACKEND] Descripción de la tarea de backend
- [ ] [FRONTEND] Descripción de la tarea de frontend  
- [ ] [DB] Cambios de base de datos si aplica
- [ ] [TEST] Casos de prueba unitarios
- [ ] [TEST] Casos de prueba de integración
- [ ] [DOCS] Actualizar documentación técnica
- [ ] [REVIEW] Code review y QA sign-off
```

Adaptar las tareas al contexto real del requerimiento. No agregar tareas genéricas irrelevantes.

---

## Paso 5 — Confirmación del usuario

Mostrar el documento completo y preguntar:

```
📄 Documento SDD generado. ¿Confirmás que está correcto para crear el issue en GitHub?

[SÍ, crear el issue] | [Quiero hacer cambios] | [Cancelar]
```

Si el usuario pide cambios → volver al punto relevante del flujo y regenerar.
Si confirma → proceder al Paso 6.

---

## Paso 6 — Crear el issue en GitHub

### Si hay MCP de GitHub conectado
Usar la herramienta disponible para crear el issue con:
- `title`: título del issue
- `body`: el Markdown completo generado
- `labels`: las etiquetas determinadas
- `milestone` (si el usuario lo especificó)
- `assignees` (si el usuario lo especificó)

### Si NO hay MCP de GitHub conectado
Informar al usuario y ofrecer alternativas:

```
No tengo conexión a GitHub en este momento. Podés:

1. **Copiar el Markdown** generado y pegarlo directamente en GitHub
2. **Conectar el MCP de GitHub** desde el menú de herramientas para que lo cree automáticamente
3. **Usar GitHub CLI**: 
   gh issue create --title "TÍTULO" --body "$(cat issue.md)" --label "LABELS"
```

Presentar el archivo Markdown para descarga.

---

## Referencias

- `references/sdd-schema.md` — Esquema completo de campos SDD con descripciones
- `references/sdd-template.md` — Plantilla Markdown lista para usar
- `references/github-labels.md` — Convención de labels de GitHub recomendados

---

## Notas de comportamiento

- **Idioma**: Responder siempre en el mismo idioma que el usuario
- **Tono**: Técnico pero claro. Evitar jerga innecesaria
- **Longitud**: Las preguntas deben ser concisas. El documento generado puede ser extenso
- **Supuestos**: Siempre documentar explícitamente los supuestos realizados
- **Validación**: Los acceptance criteria deben ser verificables y medibles (no subjetivos)