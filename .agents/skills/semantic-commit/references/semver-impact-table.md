# SemVer Impact — go-semantic-release

Tabla completa de tipos de commit y su impacto en versión con go-semantic-release.

---

## Tabla de impacto por defecto

| Tipo de commit | Impacto SemVer | Genera release | Aparece en CHANGELOG |
|---------------|---------------|----------------|---------------------|
| `feat` | MINOR (0.**X**.0) | ✅ Sí | ✅ Sí — sección "Features" |
| `fix` | PATCH (0.0.**X**) | ✅ Sí | ✅ Sí — sección "Bug Fixes" |
| `perf` | PATCH (0.0.**X**) | ✅ Sí | ✅ Sí — sección "Performance" |
| `revert` | PATCH (0.0.**X**) | ✅ Sí | ✅ Sí — sección "Reverts" |
| `refactor` | — | ❌ No | ❌ No |
| `docs` | — | ❌ No | ❌ No |
| `style` | — | ❌ No | ❌ No |
| `test` | — | ❌ No | ❌ No |
| `build` | — | ❌ No | ❌ No |
| `ci` | — | ❌ No | ❌ No |
| `chore` | — | ❌ No | ❌ No |
| `BREAKING CHANGE` (cualquier tipo) | **MAJOR** (**X**.0.0) | ✅ Sí | ✅ Sí — sección "Breaking Changes" |

---

## Lógica de prioridad de versión

Cuando hay múltiples commits desde la última release:

```
BREAKING CHANGE > feat > fix/perf/revert > sin release
```

Ejemplo con commits mixtos desde v1.2.3:
```
feat(auth): agregar OAuth          → candidato MINOR → v1.3.0
fix(login): corregir redirect      → candidato PATCH
chore: actualizar dependencias     → sin release
```
go-semantic-release tomará el MAYOR impacto → **v1.3.0**

---

## Configuración en `.semrelrc` para personalizar

go-semantic-release usa el archivo `.semrelrc` (JSON) para configurar qué commits generan releases:

```json
{
  "plugins": {
    "commit-analyzer": {
      "name": "default",
      "options": {
        "typeMap": {
          "feat": "minor",
          "fix": "patch",
          "perf": "patch",
          "refactor": "patch",
          "docs": "patch",
          "build": "patch"
        }
      }
    }
  }
}
```

Con esta config, `refactor`, `docs` y `build` también generarían releases PATCH.

---

## Visualización del impacto de versión

### Escenarios habituales

```
Versión actual: v0.5.3

feat          → v0.6.0   (MINOR)
fix           → v0.5.4   (PATCH)
perf          → v0.5.4   (PATCH)
feat!         → v1.0.0   (MAJOR — breaking change)
fix + feat    → v0.6.0   (el mayor gana: MINOR)
chore         → v0.5.3   (sin cambio)
```

### Versión en 0.x.x (pre-1.0)

En proyectos pre-estables (< v1.0.0), go-semantic-release sigue respetando la lógica:
- `feat` → MINOR (0.**X**.0)  
- `fix` → PATCH (0.0.**X**)
- `BREAKING CHANGE` → MAJOR (**X**.0.0) → pasaría a v1.0.0

---

## Prefijo de tag

Por defecto go-semantic-release usa el prefijo `v`:
- Tag creado: `v1.4.0`
- Configurable en `.semrelrc`: `"tagPrefix": ""`

---

## Commits que NO deben llegar a main

Estos patrones bloquean el release si go-semantic-release los detecta como el único cambio:

- `chore(wip):` — trabajo en progreso
- `style:` — solo formato
- `docs:` — solo documentación (a menos que se configure)

Comunicar esto al usuario cuando genera este tipo de commits para que sea consciente.