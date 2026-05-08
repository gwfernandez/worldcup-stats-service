# Conventional Commits — Spec Completa y Edge Cases

Referencia detallada para casos que no están cubiertos en el SKILL.md principal.

---

## Estructura completa del mensaje

```
<type>[optional scope][optional !]: <description>
<blank line>
[optional body]
<blank line>
[optional footer(s)]
```

### Reglas de la spec (v1.0.0)

1. El commit DEBE iniciar con un tipo (sustantivo: `feat`, `fix`, etc.)
2. El scope es OPCIONAL, entre paréntesis: `feat(parser):`
3. El `!` es OPCIONAL después del tipo/scope, indica breaking change
4. El separador `:` + espacio es OBLIGATORIO después del tipo/scope/!
5. La descripción es OBLIGATORIA después del separador
6. El body es OPCIONAL, separado por línea en blanco
7. Los footers son OPCIONALES, separados del body por línea en blanco
8. `BREAKING CHANGE:` en footer es equivalente a `!` en el tipo

---

## Footers válidos

Formato: `token: value` o `token #value`

```
BREAKING CHANGE: descripción del cambio incompatible
Reviewed-by: @username
Closes #123
Refs #456, #789
Co-authored-by: Nombre <email@example.com>
```

Reglas de footers:
- `BREAKING CHANGE` DEBE ser en mayúsculas
- `BREAKING-CHANGE` es sinónimo aceptado
- Múltiples footers: uno por línea

---

## Casos edge frecuentes

### Scope con múltiples palabras
Usar kebab-case:
```
feat(user-profile): agregar campo de biografía
fix(payment-gateway): manejar timeout en Stripe
```

### Scope de archivo o módulo
```
refactor(utils/date): extraer función formatDate
test(api/users): agregar casos para endpoint DELETE
```

### Breaking change sin `!` explícito (solo footer)
```
feat(config): cambiar formato de archivo de configuración

BREAKING CHANGE: el archivo ahora debe ser YAML en lugar de JSON.
Renombrar `.config.json` a `.config.yml` y adaptar la estructura.
```
Ambas formas son válidas. Se recomienda usar `!` Y el footer para máxima claridad.

### Múltiples breaking changes
```
refactor(api)!: rediseñar endpoints de autenticación

BREAKING CHANGE: /auth/login ahora requiere JSON body en lugar de form-data
BREAKING CHANGE: el token de respuesta cambia de JWT a opaque token
```

### Revert con referencia
```
revert: revert "feat(user): agregar avatar"

This reverts commit a1b2c3d4.
Refs #204
```

### Commit de merge (evitar si es posible)
```
chore: merge branch 'feature/login' into main
```

### Sin scope cuando el cambio es transversal
```
feat: agregar soporte para dark mode
fix: corregir comportamiento en Safari
```

---

## Qué NO hacer

```bash
# ❌ Tipo en mayúscula
Feat(auth): agregar login

# ❌ Punto final en descripción
feat(auth): agregar login.

# ❌ Descripción muy vaga
fix: arreglar bug

# ❌ Header muy largo (>72 chars)
feat(authentication-module): agregar soporte completo para OAuth2 con Google y Facebook

# ❌ Pasado en lugar de imperativo
feat(auth): agregué login con Google
feat(auth): se agregó login con Google

# ❌ Breaking change sin documentar
feat(api)!: cambiar respuesta
# Falta footer BREAKING CHANGE con descripción

# ❌ Mezclar múltiples cambios sin relación
feat(auth): agregar login + fix(ui): corregir botón
```

---

## Idioma del commit

go-semantic-release es agnóstico al idioma. Usar el idioma del equipo/proyecto:

- Si el proyecto tiene commits previos en inglés → inglés
- Si el equipo trabaja en español → español
- No mezclar idiomas en el mismo repositorio
- El tipo SIEMPRE en inglés (es parte de la spec técnica)

```bash
# Español válido
feat(pagos): agregar soporte para Mercado Pago

# Inglés válido
feat(payments): add Mercado Pago support
```