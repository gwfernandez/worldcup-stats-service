# GitHub Labels — Convención Recomendada

Sistema de etiquetas para categorizar issues de forma consistente.

---

## Tipo (obligatorio — exactamente 1)

| Label | Color | Descripción |
|-------|-------|-------------|
| `type:feature` | `#0075ca` | Nueva funcionalidad |
| `type:bug` | `#d73a4a` | Error o comportamiento incorrecto |
| `type:improvement` | `#a2eeef` | Mejora de algo existente |
| `type:task` | `#e4e669` | Trabajo técnico / infra / refactor |
| `type:epic` | `#7057ff` | Agrupador de múltiples issues |
| `type:spike` | `#f9d0c4` | Investigación o PoC |
| `type:docs` | `#0052cc` | Solo documentación |

---

## Prioridad (obligatorio — exactamente 1)

| Label | Color | Descripción |
|-------|-------|-------------|
| `priority:critical` | `#b60205` | Bloquea operación, urgente |
| `priority:high` | `#e11d48` | Sprint actual |
| `priority:medium` | `#f97316` | Próximos sprints |
| `priority:low` | `#84cc16` | Backlog / deseable |

---

## Área (1 o más)

| Label | Color | Descripción |
|-------|-------|-------------|
| `area:backend` | `#5319e7` | Lógica del servidor |
| `area:frontend` | `#006b75` | UI / UX |
| `area:database` | `#fbca04` | Cambios de datos o esquema |
| `area:api` | `#0e8a16` | Endpoints / integraciones |
| `area:auth` | `#c5def5` | Autenticación / autorización |
| `area:infra` | `#bfd4f2` | Infraestructura / DevOps |
| `area:mobile` | `#d4c5f9` | Aplicación móvil |
| `area:testing` | `#e6e6e6` | Pruebas |

---

## Estado (manejado por el equipo)

| Label | Color | Descripción |
|-------|-------|-------------|
| `status:in-progress` | `#0075ca` | En desarrollo |
| `status:blocked` | `#d73a4a` | Bloqueado por dependencia |
| `status:needs-review` | `#fbca04` | Esperando revisión |
| `status:ready` | `#0e8a16` | Listo para comenzar |

---

## Complejidad (opcional)

| Label | Color | Descripción |
|-------|-------|-------------|
| `size:xs` | `#ededed` | < 2 horas |
| `size:s` | `#c2e0c6` | < 1 día |
| `size:m` | `#fef2c0` | 2–3 días |
| `size:l` | `#fad8c7` | ~1 semana |
| `size:xl` | `#f9d0c4` | > 1 semana |

---

## Especiales

| Label | Color | Descripción |
|-------|-------|-------------|
| `good first issue` | `#7057ff` | Apto para nuevos colaboradores |
| `help wanted` | `#008672` | Se busca colaboración |
| `wontfix` | `#ffffff` | No se implementará |
| `duplicate` | `#cfd3d7` | Duplicado de otro issue |
| `needs-clarification` | `#e4e669` | Requiere más información |

---

## Ejemplos de combinaciones

```
# Feature de prioridad alta en backend y frontend
type:feature  priority:high  area:backend  area:frontend  size:m

# Bug crítico en autenticación
type:bug  priority:critical  area:auth  size:s

# Tarea de infraestructura sin urgencia
type:task  priority:low  area:infra  size:l
```