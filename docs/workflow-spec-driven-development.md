# 🧠 Workflow “Spec Driven Development”

## 1. 📌 Definir requisitos (como “spec mode”)

Siempre arrancar acá.

**Prompt:**

```
Actúa como Product Owner.
Quiero construir: [tu idea]

Genera:
- Objetivo del sistema
- Requisitos funcionales
- Requisitos no funcionales
- Casos de uso principales
- Supuestos
- Riesgos

Sé claro y estructurado.
```

📁 Guardar en:

```
/docs/requirements.md
```

---

## 2. 👤 Historias de usuario

**Prompt:**

```
A partir de estos requisitos, genera historias de usuario en formato:

Como [tipo de usuario]
Quiero [acción]
Para [beneficio]

Incluye criterios de aceptación (Gherkin).
Prioriza (Alta/Media/Baja).
```

📁 Guardar en:

```
/docs/user-stories.md
```

---

## 3. 🏗️ Diseño técnico

**Prompt:**

```
Actúa como Software Architect.

Con estos requisitos e historias:
[contenido previo]

Define:
- Arquitectura
- Stack tecnológico recomendado
- Diagrama lógico (texto)
- Modelos de datos
- Endpoints API
- Decisiones técnicas clave
```

📁 Guardar en:

```
/docs/architecture.md
```

---

## 4. 📋 Plan de tareas

**Prompt:**

```
Divide el diseño en tareas técnicas pequeñas.

Para cada tarea:
- Nombre
- Descripción
- Archivos involucrados
- Dependencias
- Estimación (S/M/L)

Ordenalas por prioridad.
```

📁 Guardar en:

```
/docs/tasks.md
```

---

# ⚙️ Estructura del proyecto

```
/project
  /docs
    requirements.md
    user-stories.md
    architecture.md
    tasks.md
```

---

# 🧠 Tips clave

## 1. No mezclar etapas

Evitar pedir requisitos + código juntos.

---

## 2. Usar contexto mínimo

Demasiado contexto reduce precisión.

---

## 3. Iterar en pasos chicos

Mejor calidad que generar todo de una.

---

## 4. Versionar todo

Usar Git para cada documento.

---
