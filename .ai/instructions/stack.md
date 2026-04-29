# Stack Técnico — World Cups API

## Diagrama de Arquitectura

```
┌─────────────────────────────────┐
│         Cliente                 │
│   Postman / browser / app       │
└────────────────┬────────────────┘
                 │ HTTP
                 ▼
┌─────────────────────────────────┐
│       API — Go 1.23 + Gin       │
│  Handlers · Middleware · Routing│
│         · Validación            │
└────────────────┬────────────────┘
                 │
                 ▼
┌─────────────────────────────────┐
│           Railway               │
│  Deploy automático desde GitHub │
└────────────────┬────────────────┘
                 │ pgx
                 ▼
┌─────────────────────────────────┐
│        Neon PostgreSQL          │
│  Mundiales · Equipos · Jugadores│
│          · Partidos             │
└─────────────────────────────────┘
         ▲              ▲
         │              │
┌────────┴──────┐ ┌─────┴──────────┐
│     sqlc      │ │golang-migrate  │
│ Genera código │ │ Migraciones SQL│
│      Go       │ │                │
└───────────────┘ └────────────────┘
```

---

## Componentes del Stack

### 🔵 Go 1.23

**Versión:** 1.23

**Beneficios para este proyecto:**
- Compilación estática que genera un binario único, ideal para deploy en Railway sin dependencias externas
- Concurrencia nativa con goroutines, eficiente para manejar múltiples requests simultáneos
- Tipado estático que reduce errores en tiempo de ejecución al mapear datos de la DB
- Arranque extremadamente rápido, perfecto para entornos cloud con cold starts

---

### ⚙️ Gin (Framework HTTP)

**Repositorio:** github.com/gin-gonic/gin

**Beneficios para este proyecto:**
- El framework HTTP más usado en Go, con enorme comunidad y documentación
- Validación de input con binding declarativo (`binding:"required"`) sin código extra
- Manejo de rutas con parámetros (`/api/v1/confederations/:id`) simple y legible
- Middleware integrado para logging, recovery y CORS
- Rendimiento muy alto, ideal para una API que expone datos históricos de mundiales

---

### 🏗️ Clean Architecture (handler → service → repository)

**Beneficios para este proyecto:**
- Separación clara de responsabilidades: cada capa tiene una única función
- El service no conoce ni HTTP ni SQL, lo que facilita pruebas unitarias sin levantar servidor ni DB
- Si en el futuro se cambia PostgreSQL por otra DB, solo se modifica el repository sin tocar el resto
- Estructura escalable: agregar nuevas entidades (equipos, jugadores, partidos) sigue el mismo patrón
- Código más legible y mantenible a largo plazo

---

### 🐘 PostgreSQL en Neon

**Plan:** Free tier — 3 GiB de almacenamiento
**Branches:** `main` (PROD) / `dev` (DEV)

**Beneficios para este proyecto:**
- Motor relacional ideal para datos altamente relacionados: mundiales → grupos → partidos → goles → jugadores
- JOINs eficientes para estadísticas y rankings sin duplicar datos
- Full-text search nativo para búsquedas por nombre de jugador o equipo
- Window functions y CTEs para consultas analíticas complejas (goleadores históricos, etc.)
- Branching de Neon permite tener DEV y PROD aislados desde el día uno, sin costo adicional
- 3 GiB es más que suficiente para almacenar todos los mundiales desde 1930 hasta la actualidad

---

### 🚀 Railway (Deploy)

**Plan:** Free tier — $5 crédito mensual

**Beneficios para este proyecto:**
- Deploy automático con cada `git push` a la rama `main` de GitHub, sin configuración manual
- Variables de entorno (`DATABASE_URL`) gestionadas desde el dashboard de forma segura
- Logs en tiempo real para debugging desde el navegador
- Soporte nativo para binarios Go sin necesidad de configurar Dockerfile
- El crédito de $5/mes cubre sobrado una API pequeña con tráfico bajo

---

### 🔌 pgx (Driver PostgreSQL)

**Repositorio:** github.com/jackc/pgx/v5

**Beneficios para este proyecto:**
- Driver más moderno y performante para Go + PostgreSQL
- Soporte nativo para tipos de PostgreSQL (BIGSERIAL, VARCHAR, etc.) sin conversiones manuales
- Connection pooling integrado con `pgxpool`, eficiente para una API con múltiples requests concurrentes
- Compatible con sqlc de forma nativa
- Manejo de errores más expresivo que el driver estándar `database/sql`

---

### 📝 sqlc (Generador de código)

**Repositorio:** github.com/sqlc-dev/sqlc

**Beneficios para este proyecto:**
- Escribís SQL puro y sqlc genera automáticamente el código Go tipado correspondiente
- Cero reflection en runtime, a diferencia de los ORMs
- Los errores de queries se detectan en tiempo de compilación, no en producción
- El SQL generado es predecible y auditable, ideal para un proyecto con datos históricos donde la precisión importa
- Configuración simple con `sqlc.yaml` apuntando a `/db/queries`

---

### 🗄️ golang-migrate (Migraciones)

**Repositorio:** github.com/golang-migrate/migrate

**Beneficios para este proyecto:**
- Control de versiones del schema de la DB, igual que git pero para tablas y columnas
- Cada migración tiene un archivo `up` (aplicar) y `down` (revertir), permitiendo rollbacks seguros
- Se puede correr desde el código Go al arrancar la aplicación o desde CLI
- Permite replicar el schema exacto tanto en `dev` como en `prod` con un solo comando
- Fundamental para un proyecto que va a crecer en entidades (confederaciones → equipos → jugadores → partidos)

---

### 🐙 GitHub (Repositorio + CI/CD)

**Beneficios para este proyecto:**
- Control de versiones del código con historial completo de cambios
- Integración directa con Railway para auto-deploy en cada push a `main`
- Rama `main` → PROD / rama `dev` → desarrollo, espejando la estrategia de branches de Neon
- Posibilidad de agregar GitHub Actions en el futuro para correr tests antes del deploy

---

## Resumen del Stack

| Componente | Tecnología | Versión |
|---|---|---|
| Lenguaje | Go | 1.23 |
| Framework HTTP | Gin | latest |
| Arquitectura | Clean Architecture | handler → service → repository |
| Base de datos | PostgreSQL (Neon) | free tier, 3 GiB |
| Deploy | Railway | free tier |
| Driver DB | pgx | v5 |
| Queries | sqlc | latest |
| Migraciones | golang-migrate | latest |
| Repositorio | GitHub | auto-deploy |

---

## Flujo de trabajo

```
Desarrollo local  →  git push (branch dev)  →  Neon dev branch
                  →  git push (branch main) →  Railway deploy → Neon prod branch
```

---

*Documento generado para el proyecto World Cups API*
