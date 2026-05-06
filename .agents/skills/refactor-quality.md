# Skill — Refactorización y Calidad de Código (Post-Desarrollo)

Este skill se activa tras completar la implementación de una funcionalidad para asegurar que el código cumple con los estándares de **Clean Architecture**, legibilidad y robustez del proyecto.

## Objetivo
Transformar código funcional en código de calidad senior, siguiendo el flujo de **Sugerencia + Aplicación**.

---

## Áreas de Revisión

### 1. Arquitectura y Diseño
- **Clean Architecture:** Verificar que las dependencias fluyan hacia adentro (Handler -> Service -> Repository). No debe haber lógica de base de datos en el service ni lógica de negocio en el handler.
- **Interfaces:** Asegurar que el Service y el Repository dependan de interfaces para facilitar el testing.

### 2. Calidad de Código (Go Idioms)
- **Error Handling:** Verificar que los errores no se ignoren y se propaguen correctamente.
- **Naming:** Usar nombres descriptivos pero concisos, siguiendo las convenciones de Go.
- **Complejidad:** Identificar funciones con alta complejidad ciclomática (demasiados IFs o bucles anidados) y sugerir su división.

### 3. Documentación y Estándares
- **Godoc:** Asegurar que toda función, struct o interfaz pública tenga un comentario de documentación descriptivo.
- **Errores:** Validar que las respuestas de error en los handlers usen el formato `{"error": "mensaje"}`.

---

## Procedimiento (Flujo Sugerencia + Aplicación)

1. **Auditoría:** Analizar el código recién escrito buscando oportunidades de mejora en las áreas mencionadas.
2. **Presentación de Propuesta:** 
   - Listar los hallazgos de forma clara (ej: "Bug detectado", "Mejora de arquitectura", "Falta documentación").
   - Explicar brevemente el **por qué** de cada mejora sugerida.
   - Mostrar un ejemplo o descripción del cambio propuesto.
3. **Solicitud de Aprobación:** Preguntar explícitamente: *"¿Deseas que aplique estas mejoras automáticamente?"*.
4. **Ejecución:** Solo tras la confirmación, proceder a editar los archivos y realizar las correcciones.
5. **Verificación:** Ejecutar los tests unitarios (`go test ./...`) para asegurar que el refactor no rompió la funcionalidad existente.
