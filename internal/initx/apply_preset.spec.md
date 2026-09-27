<!-- @anchors
  code: APPRP
  updated_at: 2026-09-26
  layer: infra
-->
# ApplyPreset — writes a stack preset's layers into the configuration and deduces one identity prefix per module

> **Code**: `APPRP`

## Overview

When `anchors init` offers a stack preset, the chosen preset has to become real layers in the project's
configuration, and a modular preset also needs one identity prefix per module so that the codes generated
later say which module a unit belongs to. This unit does both.

The preset's layers are copied into the configuration as they are: a layer of the same name is replaced by
the preset's, and every layer of another name (the artifact layers chosen elsewhere, a layer the user already
declared) stays. The prefixes are not written into the layers: the identity prefix is per module, not per
layer, so the unit returns the module-to-prefix mapping for the caller to show and use.

Each module receives the two-letter prefix the identity code derives from the module's directory name. When
two modules would share a prefix, the one that comes later in alphabetical order keeps its first letter and
takes the first second letter still free; when all 26 of that initial are taken, it takes the first free
prefix of the other initials, so no two modules ever share an identity. Two modules with the same folder name
under different parents are two modules: each is keyed by its path, and each gets its own prefix. Modules are processed in alphabetical order, so the same set of
modules always yields the same mapping, whatever order the detection returned them in.

The unit is pure: it works on the module paths it receives and never looks at the disk.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | a configuration, with or without a layer set | no configuration at all | the caller: init always builds one before applying a preset |
| the preset | an entry of the preset catalog | a preset with no layers | the preset catalog (INCTN-I01) |
| the modules | module directory paths found under the preset's module directory, possibly none | paths that are not module directories; more than 676 modules (every two-letter prefix taken) | the caller: it globs the preset's module directory |

## Effects

| Effect | Description |
| --- | --- |
| `APPRP-B01` | Applying a preset adds every layer of the preset to the configuration, creating the layer set when the configuration has none, and keeps every layer of another name. |
| `APPRP-B02` | A layer whose name matches a preset layer is replaced by the preset's declaration. |
| `APPRP-B03` | Each module receives a two-letter prefix (`DeduceModulePrefixes`), keyed by the last segment of its directory path (a trailing slash is ignored) when no other module shares that segment, equal to the identity code's module prefix when no other module took it. |
| `APPRP-B04` | When a module's prefix is already taken, the module keeps the first letter and takes the first free second letter from A to Z; with all 26 taken, it takes the first free prefix of the other first letters, from A. |
| `APPRP-B05` | Applying a preset returns the module-to-prefix mapping deduced for the detected modules, and an empty mapping when there is none. |
| `APPRP-B06` | Modules that share a folder name are each keyed by their path (a trailing slash is ignored) and each receives its own prefix. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `APPRP-I01` | The same set of modules gives the same prefixes, whatever order they arrive in. | deduces the prefixes of the same modules in two orders and compares the mappings |
| `APPRP-I02` | No two modules share a prefix while a free two-letter prefix exists. | deduces the prefixes of thirty modules with the same initial and checks every prefix is unique |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `APPRP-X01` | Prefixes are deduced from the given paths alone; the disk is not read. | The caller already found the modules; reading again would make the result depend on the moment it runs, and would stop the function from being testable without a tree. |

## Errors

none — the unit copies layers and computes prefixes from strings it receives; there is no input it rejects and no resource it can fail to reach.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/code/code.go` | `ModulePrefix` | core — the identity code's module prefix |
| DEP2 | `internal/config/config.go` | `Config`, `Layer` | core — project configuration |
| DEP3 | `internal/initx/presets.go` | `Preset`, `ToLayers` | infra — the preset catalog (INCTN) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
