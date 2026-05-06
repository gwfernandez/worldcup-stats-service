# Skill — Generar Tests Unitarios con Mocks

Este skill automatiza la creación de pruebas unitarias para las capas de `service` y `repository`, asegurando el cumplimiento del objetivo de **90% de cobertura** utilizando `testify` y `pgxmock`.

## Contexto técnico
- **Librería de aserciones:** `github.com/stretchr/testify/assert` y `require`
- **Mock de base de datos:** `github.com/pashagolub/pgxmock/v3`
- **Ubicación:** Los tests deben estar en el mismo paquete que el código fuente (`package service` o `package repository`).

---

## Procedimiento para Services

1. **Analizar la interfaz:** Identificar el método a testear y sus dependencias (normalmente una interfaz de repository).
2. **Setup del Mock:**
   - Crear un mock del repository usando la interfaz correspondiente.
   - Instanciar el service inyectando el mock.
3. **Definir casos de prueba (Table-Driven Tests):**
   - **Success**: El mock retorna los datos esperados, el service retorna `nil` error.
   - **Error del Repository**: El mock retorna un error, el service debe propagarlo.
   - **Validación de lógica**: Si el service tiene lógica extra (ej: check de existencia), añadir casos para esos caminos.
4. **Ejecución y Verificación:**
   - Usar `mock.AssertExpectations(t)` para asegurar que todas las llamadas esperadas ocurrieron.

---

## Procedimiento para Repositories

1. **Setup de pgxmock:**
   - Iniciar un `pgxmock.NewPool()`.
   - Asegurar el `defer mock.Close()`.
2. **Simular Consultas SQL:**
   - Usar `mock.ExpectQuery` para SELECTs o `mock.ExpectExec` para INSERT/UPDATE/DELETE.
   - Definir los parámetros exactos y las filas retornadas (`NewRows`).
3. **Aserciones:**
   - Verificar que el objeto retornado por el repository mapea correctamente los datos del mock.
   - Validar el manejo de errores (ej: `pgx.ErrNoRows`).

---

## Criterios de Calidad
- Cada test debe ser independiente.
- Usar nombres descriptivos: `TestServiceName_MethodName_Success`, `TestServiceName_MethodName_NotFound`.
- No testear implementaciones internas, sino el comportamiento de la interfaz.
- **Importante:** Si el coverage no llega al 90%, identificar los bloques "rojos" y añadir casos específicos para ellos.
