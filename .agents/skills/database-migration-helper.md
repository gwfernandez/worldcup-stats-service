# Skill — Gestor de Migraciones de Base de Datos

Este skill proporciona un flujo seguro para la creación y gestión de cambios en el esquema de la base de datos usando `golang-migrate`.

## Convenciones de Migración
- **Ubicación:** `db/migrations/`
- **Formato de nombre:** `00X_descripcion_breve.up.sql` y `00X_descripcion_breve.down.sql`
- **Numeración:** Incremental de tres dígitos (001, 002, 003...).

---

## Procedimiento

### 1. Creación de una nueva migración
- Identificar el siguiente número disponible.
- Crear el archivo `.up.sql` con las sentencias `CREATE TABLE`, `ALTER TABLE`, etc.
- Crear SIEMPRE el archivo `.down.sql` que revierte exactamente lo hecho en el `up`.

### 2. Estándares SQL (PostgreSQL)
- Usar `BIGSERIAL PRIMARY KEY` para IDs.
- Usar `TIMESTAMPTZ` para fechas con zona horaria.
- Añadir constraints de `NOT NULL` y `UNIQUE` donde corresponda.
- Nombres de tablas en plural (ej: `teams`, `matches`).

### 3. Verificación
- Ejecutar las migraciones localmente para validar la sintaxis:
  ```bash
  migrate -path db/migrations -database "$DATABASE_URL" up
  ```
- Probar un "Rollback" para asegurar que el archivo `.down.sql` funciona:
  ```bash
  migrate -path db/migrations -database "$DATABASE_URL" down 1
  migrate -path db/migrations -database "$DATABASE_URL" up
  ```
- Ejecutar `sqlc generate` inmediatamente después para asegurar que el código Go se actualiza con el nuevo esquema.

---

## Seguridad
- **Prohibido:** Modificar archivos de migración ya existentes (committeados). Si algo está mal, se crea una nueva migración correcciva.
- **Validación:** Antes de hacer push, asegurar que no hay conflictos de numeración con la rama `main`.
