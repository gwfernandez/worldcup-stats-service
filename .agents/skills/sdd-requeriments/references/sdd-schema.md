# SDD Schema — Campos del Requerimiento

Referencia completa de todos los campos del documento SDD con descripciones, ejemplos y reglas de validación.

---

## Campos Obligatorios

### `title`
**Qué es**: Título corto y descriptivo del issue (máx. 80 caracteres)  
**Formato**: `[TIPO] Verbo + objeto + contexto`  
**Ejemplos**:
- `[FEATURE] Exportar reportes de ventas en formato PDF`
- `[BUG] Error 500 al guardar usuario sin email`
- `[IMPROVEMENT] Optimizar tiempo de carga del dashboard`

**Reglas**:
- Comenzar con verbo en infinitivo
- Incluir el tipo entre corchetes
- No usar jerga interna sin explicación

---

### `type`
**Valores válidos**:
| Valor | Cuándo usarlo |
|-------|--------------|
| `Feature` | Nueva funcionalidad que no existe |
| `Bug` | Comportamiento incorrecto o inesperado |
| `Improvement` | Mejora de algo que ya funciona |
| `Task` | Trabajo técnico sin impacto directo en usuario (refactor, infra) |
| `Epic` | Agrupador de múltiples features o stories relacionadas |
| `Spike` | Investigación o prueba de concepto |

---

### `description`
**Qué es**: Explicación clara de qué se necesita y por qué  
**Estructura sugerida**:
```
**Contexto**: [Situación actual que origina el requerimiento]
**Problema / Necesidad**: [Qué falla o qué falta]
**Solución propuesta**: [Qué se debe implementar]
**Impacto esperado**: [Beneficio para el usuario o el sistema]
```

---

### `acceptance_criteria`
**Qué es**: Lista de condiciones verificables que deben cumplirse para cerrar el issue  
**Formato recomendado** (Given/When/Then):
```
- DADO que [contexto], CUANDO [acción], ENTONCES [resultado esperado]
```

**También válido** (lista simple):
```
- El PDF se genera en menos de 3 segundos
- El archivo incluye el logo de la empresa en el header
- Si el reporte está vacío, se muestra un mensaje "Sin datos"
```

**Reglas**:
- Mínimo 2 criterios, máximo 10
- Deben ser verificables por QA sin ambigüedad
- No usar términos subjetivos ("rápido", "bonito", "fácil")

---

## Campos Importantes

### `user_story`
**Formato**: `Como [rol], quiero [acción], para [beneficio]`  
**Ejemplo**: `Como administrador, quiero exportar el reporte mensual en PDF, para enviárselo al cliente sin necesidad de acceso al sistema`

---

### `affected_components`
**Qué es**: Lista de módulos, servicios, pantallas o APIs que se modificarán  
**Ejemplos**: `ReportService`, `Dashboard UI`, `PDF Generator`, `users table`, `/api/reports endpoint`

---

---
+
+### `suggested_branch`
+**Qué es**: Nombre de la rama de Git recomendada para este trabajo.  
+**Formato**: `tipo/descripción-corta-con-guiones`  
+**Ejemplo**: `feature/exportar-pdf-reportes`
+
+---
+
+### `priority`
**Valores**:
| Valor | Descripción |
|-------|-------------|
| `Critical` | Bloquea operación, requiere atención inmediata |
| `High` | Importante para el sprint actual |
| `Medium` | Planificable para próximos sprints |
| `Low` | Mejora deseable sin urgencia |

---

## Campos Opcionales

### `assumptions`
Lista de supuestos asumidos al escribir el requerimiento.  
Ejemplo: *"Se asume que el servicio de generación de PDF ya está integrado"*

### `out_of_scope`
Qué NO está incluido en este issue para evitar scope creep.  
Ejemplo: *"No incluye envío automático por email (issue separado)"*

### `dependencies`
Otros issues, servicios externos o condiciones previas necesarias.  
Ejemplo: *"Requiere completar #123 (integración con S3)"*

### `mockups_refs`
Links a diseños en Figma, screenshots o referencias visuales.

### `technical_notes`
Consideraciones técnicas relevantes para el desarrollador.  
Ejemplo: *"Usar librería WeasyPrint existente en el proyecto"*

### `estimated_complexity`
**Valores**: `XS` (< 2h) / `S` (< 1 día) / `M` (2-3 días) / `L` (1 semana) / `XL` (> 1 semana)

---

## Campos de GitHub

### `labels`
Ver `github-labels.md` para la convención completa.  
Labels mínimos: tipo + área + prioridad

### `milestone`
Sprint o versión objetivo (ej: `v2.4`, `Sprint 12`)

### `assignees`
Username(s) de GitHub del/los responsable(s)