# Estrategia de Versionado de API

Este documento detalla la estrategia de versionado adoptada para la API de World Cups Stats Service.

## Enfoque Elegido: Versionado por Header (Custom Header)

Se ha decidido implementar un esquema de versionado dinámico basado en un header HTTP personalizado: `X-API-Version`.

### ¿Por qué esta decisión?
1. **URLs Limpias**: Mantiene las rutas de la API limpias y semánticas (`/api/confederations` en lugar de `/api/v1/confederations`), enfocándose en el recurso y no en su versión.
2. **Flexibilidad**: Facilita la transición entre versiones sin requerir actualizaciones en los endpoints expuestos si no hay "breaking changes".
3. **Escalabilidad**: Evita la proliferación de rutas redundantes en el router.

## Cómo Usarlo

### Consumidores de la API

Para acceder a una versión específica de la API, los clientes deben incluir el header `X-API-Version` en sus solicitudes HTTP.

**Ejemplo de Request:**
```http
GET /api/confederations HTTP/1.1
Host: api.example.com
X-API-Version: 1
```

### Comportamiento por Defecto

Si no se provee el header `X-API-Version`, el sistema asumirá automáticamente la versión por defecto de la API. 
* **Versión por defecto actual**: `1`

### Respuestas de la API

El sistema responderá indicando qué versión fue efectivamente procesada a través del header `X-API-Version-Used`.

**Ejemplo de Response:**
```http
HTTP/1.1 200 OK
X-API-Version-Used: 1
Content-Type: application/json
```

Si el cliente solicita una versión que no existe o no está soportada para ese endpoint, la API responderá con un error `400 Bad Request`.

## Implementación Técnica (Para Desarrolladores)

### Middlewares

La estrategia se apoya en dos middlewares de Gin:
1. `Versioning()`: Middleware global. Extrae la versión del header, establece el fallback y expone el valor en el contexto (gin.Context).
2. `RequireVersion(v int)`: Middleware a nivel de grupo de rutas. Asegura que la versión en el contexto coincida con la requerida por el handler.

### Organización de Código

Los handlers deben agruparse en subpaquetes según la versión a la que corresponden.
Ejemplo: `internal/handler/v1/confederation_handler.go`
