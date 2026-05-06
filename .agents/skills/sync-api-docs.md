# Skill — Sincronizar Documentación de API

Este skill asegura que los documentos de referencia (`AGENTS.md` y `README.md`) reflejen siempre el estado actual del código, especialmente en lo que respecta a endpoints y entidades.

## Cuándo ejecutar este skill
- Después de añadir una nueva tabla en una migración.
- Después de crear o modificar un handler con nuevos endpoints.
- Antes de realizar un Pull Request.

---

## Procedimiento

### 1. Auditoría de Entidades
- Revisar `/internal/domain` para nuevas entidades.
- Actualizar la tabla de **"Entidades del dominio"** en `AGENTS.md`.
- Clasificar como "Implementada" o "Planificada".

### 2. Auditoría de Endpoints
- Revisar los archivos en `/internal/handler` (especialmente el método `RegisterRoutes`).
- Actualizar la sección de **"Endpoints implementados"** en `AGENTS.md` y `README.md`.
- Formato de tabla: `Método | Ruta | Descripción`.

### 3. Auditoría del Stack
- Si se ha añadido una dependencia importante en `go.mod`, actualizar la tabla del **"Stack tecnológico"**.

---

## Reglas de Formato
- Mantener el estilo de tablas Markdown existente.
- Asegurar que los links a archivos funcionen correctamente.
- Mantener las descripciones breves y en español.
- **Consistencia:** Si un cambio se refleja en `AGENTS.md`, debe evaluarse si también corresponde al `README.md` (pensado para humanos) para que no haya discrepancias.
