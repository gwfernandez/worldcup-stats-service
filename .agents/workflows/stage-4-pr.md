---
description: Stage 4 — Commit y Pull Request
---

# Stage 4 — Commit y Pull Request

## Propósito

Commitear los cambios siguiendo Conventional Commits con descripciones en español, hacer push de la rama y abrir el Pull Request contra `main` con toda la información necesaria para el reviewer. Luego de que el PR sea mergeado, finalizar el issue asociado en GitHub.

---

## Pasos

### 1. Verificación final del estado del repo

```bash
git status
go build ./...
go vet ./...
go test ./... 2>&1 | tail -10
```

Todo debe estar en verde antes de commitear. Si hay algo pendiente → resolver antes de continuar.

---

### 2. Revisar el diff completo

```bash
git diff main --stat
git diff main
```

Confirmar que:
- Solo están los cambios relacionados con el issue
- No hay archivos accidentales (`.env`, `.coverage/`, binarios compilados)
- No hay `fmt.Println` ni código comentado olvidado
- El archivo `.env` no está incluido

---

### 3. Verificar documentación actualizada

Si el issue agregó o modificó endpoints, verificar que estén documentados:

- [ ] `README.md` — tabla de endpoints actualizada
- [ ] `AGENTS.md` — sección "Endpoints implementados" actualizada (si aplica)
- [ ] Godoc en todas las funciones públicas nuevas o modificadas

Si falta documentación → actualizar antes de commitear.

---

### 4. Stagear los cambios

```bash
git add .
git status
```

Confirmar que el staging es correcto. Si hay archivos que no deben ir → quitarlos con `git restore --staged <archivo>`.

---

### 5. Crear el commit

Formato obligatorio: **Conventional Commits**

```
<tipo>(<scope>): <descripción en español, presente, sin punto final>
```

**Regla de idioma:** el tipo y el scope van en inglés (convención técnica), la descripción **siempre en español**.

**Tipos válidos:** `feat` · `fix` · `perf` · `refactor` · `docs` · `style` · `test` · `build` · `ci` · `chore` · `revert`

**Para breaking changes:** usar `!` después del tipo o `BREAKING CHANGE:` en el footer.

```bash
git commit -m "<tipo>(<scope>): <descripción en español>"
```

Ejemplos correctos:
```bash
git commit -m "feat(api): agregar endpoint de goleadores por mundial"
git commit -m "fix(handler): corregir código de respuesta cuando el año no existe"
git commit -m "refactor(service): extraer lógica de validación de año a helper"
git commit -m "test(repository): agregar tests para caso de lista vacía"
git commit -m "docs(readme): documentar nuevos endpoints de goleadores"
```

---

### 6. Push de la rama

```bash
git push origin <nombre-de-la-rama>
```

---

### 7. Leer el baseline para el PR

```bash
cat /tmp/baseline.json
cat /tmp/stage3_coverage.txt | grep "total:"
```

---

### 8. Abrir el Pull Request

El título del PR sigue el mismo formato de Conventional Commits con descripción en español.
El cuerpo del PR se redacta **completamente en español**.

```bash
gh pr create \
  --base main \
  --title "<tipo>(<scope>): <descripción en español>" \
  --body "$(cat <<'EOF'
## Closes #<número-de-issue>

## ¿Qué hace este PR?

<descripción de 2-3 oraciones en español explicando qué se implementó y por qué>

## Cambios realizados

### Archivos creados
- `internal/domain/...` — descripción
- `internal/repository/...` — descripción
- `internal/service/...` — descripción
- `internal/handler/...` — descripción
- `db/queries/...` — descripción (si aplica)
- `db/migrations/...` — descripción (si aplica)

### Archivos modificados
- `cmd/main.go` — descripción del cambio (si aplica)
- `README.md` — descripción del cambio (si aplica)

## Nuevos endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/...` | ... |

(Omitir sección si no hay nuevos endpoints)

## Decisiones de implementación

<Si se tomó alguna decisión no obvia de arquitectura o diseño, documentarla aquí. Si todo sigue los patrones establecidos, omitir esta sección.>

## Testing

| Métrica | Antes | Después |
|---|---|---|
| Tests passing | <baseline> | <actual> |
| Tests nuevos | — | <cantidad> |
| Coverage handler | <baseline> | <actual> |
| Coverage service | <baseline> | <actual> |
| Coverage repository | <baseline> | <actual> |
| Coverage total | <baseline> | <actual> |

Edge cases cubiertos:
- <listar los edge cases del issue que tienen test>

## Puntos de atención para el reviewer

<Si hay algo que merece atención especial: decisiones de diseño no obvias, cambios de comportamiento, impacto en otros módulos. Si no hay nada, escribir "Ninguno".>

## Checklist

- [ ] El código compila sin errores ni warnings
- [ ] Coverage ≥ 90% en paquetes modificados
- [ ] Todas las funciones públicas tienen godoc
- [ ] Sin datos sensibles en logs
- [ ] Documentación actualizada (README, AGENTS.md si aplica)
- [ ] API de solo lectura respetada (solo GET)

EOF
)"
```

---

### 9. Verificar que el PR se creó correctamente

```bash
gh pr view --web
```

Confirmar que:
- El título sigue Conventional Commits con descripción en español
- El cuerpo está completamente en español
- Incluye `Closes #<número>`
- La rama base es `main`
- Las métricas de coverage están completas
- El checklist está presente

---

### 10. Publicar Walkthrough en el issue de GitHub

Publicar un comentario final en el issue usando la integración de GitHub, indicando que el
trabajo ha concluido:

```
✅ Trabajo completado — PR listo para review

## Resumen

<descripción de 2-3 oraciones de qué se implementó y qué problema resuelve>

## Cambios principales

- `<archivo>`: <descripción del cambio>
- `<archivo>`: <descripción del cambio>

## Testing

- Tests nuevos: X
- Coverage total: X% (sin regresiones)
- Edge cases cubiertos: <lista resumida>

## Pull Request

<URL del PR> — listo para review de @gwfernandez
```

---

### 11. Finalizar el issue luego del merge del PR

Una vez que el PR haya sido mergeado en `main`, cerrar/finalizar el issue asociado:

1. Verificar que el PR esté mergeado antes de modificar el estado final del issue.
2. Si el issue pertenece a un GitHub Project, actualizar el campo de estado del item a `Done` o `Finalizado`.
3. Cerrar el issue si no fue cerrado automáticamente por `Closes #<número-de-issue>`.
4. Publicar un comentario final en el issue:

```
✅ Issue finalizado

El PR <URL del PR> fue mergeado en `main`.
El issue queda cerrado/finalizado.
```

Si el PR todavía no fue mergeado, **no cerrar ni finalizar el issue**. Reportar que queda pendiente de merge/review.

---

### 12. Reportar al usuario

```
✅ Stage 4 completado

Commit:   feat(api): agregar endpoint de goleadores por mundial
Rama:     feat/42-endpoint-goleadores
PR:       https://github.com/gwfernandez/worldcup-stats-service/pull/<número>

  Closes #42
  Tests: 105 passing (+18 nuevos)
  Coverage total: 94% (sin cambios respecto al baseline)

El PR está listo para review de @gwfernandez.
```
