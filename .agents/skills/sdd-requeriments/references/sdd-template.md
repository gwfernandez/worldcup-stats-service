# SDD Template — Plantilla Markdown para Issues

Esta es la plantilla base. Completar todos los campos; usar `[POR DEFINIR]` solo si es inevitable.

---

```markdown
# [TIPO] Título del requerimiento

## 📋 Información General

| Campo | Valor |
|-------|-------|
| **Tipo** | Feature / Bug / Improvement / Task / Epic |
| **Prioridad** | Critical / High / Medium / Low |
| **Complejidad estimada** | XS / S / M / L / XL |
| **Componentes afectados** | Módulo A, Servicio B, Pantalla C |
| **Rama sugerida** | `type/descripcion-corta` |
| **Sprint / Milestone** | [POR DEFINIR] |

---

## 👤 Historia de Usuario

> Como **[rol del usuario]**,  
> quiero **[acción o funcionalidad]**,  
> para **[beneficio o valor obtenido]**.

---

## 📝 Descripción

### Contexto
[Situación actual que origina este requerimiento. Qué existe hoy y por qué no es suficiente.]

### Problema / Necesidad
[Descripción clara del problema a resolver o la necesidad a cubrir.]

### Solución Propuesta
[Descripción de la solución esperada a alto nivel. No es la especificación técnica, sino la solución desde la perspectiva del usuario/negocio.]

### Impacto Esperado
[Qué cambia o mejora una vez implementado este requerimiento.]

---

## ✅ Criterios de Aceptación

<!-- Usar formato Given/When/Then o lista de condiciones verificables -->

- [ ] DADO que [contexto], CUANDO [acción], ENTONCES [resultado esperado]
- [ ] DADO que [contexto], CUANDO [acción], ENTONCES [resultado esperado]
- [ ] [Criterio adicional verificable]

---

## 📋 Tareas Técnicas

<!-- Desglose en subtareas. Adaptar categorías según el requerimiento. -->

### Backend
- [ ] [BACKEND] Descripción tarea 1
- [ ] [BACKEND] Descripción tarea 2

### Frontend
- [ ] [FRONTEND] Descripción tarea 1
- [ ] [FRONTEND] Descripción tarea 2

### Base de Datos
- [ ] [DB] Descripción del cambio de esquema o migración (si aplica)

### Testing
- [ ] [TEST] Pruebas unitarias (Garantizar cobertura ≥ 90% según AGENTS.md)
- [ ] [TEST] Pruebas de integración para [flujo]
- [ ] [TEST] Casos edge: [lista de casos borde]

### Documentación
- [ ] [DOCS] Actualizar README / Wiki con [sección]
- [ ] [DOCS] Documentar nueva API / endpoint (si aplica)

### Revisión
- [ ] [REVIEW] Code review por [persona/equipo]
- [ ] [QA] Sign-off de QA

---

## 🔗 Dependencias

<!-- Issues, PRs, servicios externos o condiciones previas -->

- Requiere: #[número de issue] — [descripción breve]
- Bloqueado por: [condición o issue]
- Relacionado con: #[número de issue]

---

## 🚫 Fuera de Alcance

<!-- Lo que explícitamente NO está incluido en este issue -->

- No incluye: [descripción]
- Será cubierto en: #[número o título de issue futuro]

---

## 💡 Supuestos

<!-- Supuestos realizados al redactar este requerimiento -->

- Se asume que [supuesto 1]
- Se asume que [supuesto 2]

---

## 🛠️ Notas Técnicas

<!-- Consideraciones técnicas para el desarrollador: librerías, restricciones, patrones a usar -->

[Notas técnicas relevantes o `N/A`]

---

## 📎 Referencias

<!-- Links a diseños, documentación, tickets relacionados, ejemplos -->

- Diseño en Figma: [link]
- Documentación relacionada: [link]
- Ticket de origen: [link]

---

**Labels**: `type:feature` `area:backend` `priority:high`  
**Creado por**: @[autor]  
**Fecha**: YYYY-MM-DD
```