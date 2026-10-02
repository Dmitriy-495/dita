# DITA 0.1.0 — Architecture

## Core principle

**Meta-Layer = Platform.**

DITA separates:

```text
Meta-object
    ↓
Configuration
    ↓
Runtime
    ↓
Application data
```

## Initial Meta-Layer

The first meta-object type implemented for the foundation is `Catalog`.

```text
Catalog
├── id
├── name
└── prefix
```

`Catalog` describes an application catalog. It does not contain catalog elements.

## Explicitly deferred

The following are outside the first foundation step:

- runtime execution;
- PostgreSQL persistence;
- HTTP API;
- dynamic fields;
- catalog elements;
- business-code generation;
- users;
- system attributes;
- soft delete;
- EAV storage.

These mechanisms will be introduced only when required by the next verified vertical slice.

## Architectural constraints

DITA remains:

- independent from GigaWatt-3;
- Go-based;
- PostgreSQL-oriented;
- SQL-first;
- a modular monolith;
- free of ORM as an architectural center;
- free of EAV as the primary storage model;
- free of domain-specific GigaWatt entities.
