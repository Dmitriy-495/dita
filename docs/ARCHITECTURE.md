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

## Implemented foundation

### Meta-Layer

The initial meta-object is `Catalog`:

```text
Catalog
├── id
├── name
└── prefix
```

`Catalog` describes an application catalog and does not contain its runtime elements.

The Meta-Layer supports catalog creation, retrieval and update through a typed service and repository.

### PostgreSQL persistence

PostgreSQL is the persistent storage.

- SQL-first migrations;
- typed repositories;
- metadata stored in `meta_catalogs`;
- catalog elements stored in `meta_catalog_elements`;
- no ORM-centered architecture;
- no EAV as the primary storage model.

### Runtime

Runtime loads catalog metadata and supports:

- creating catalog elements;
- retrieving elements;
- updating element names;
- generating element codes from the catalog prefix and element ID.

Example:

```text
Catalog prefix: ОО
Element ID:     1
Generated code: ОО-000001
```

### HTTP API

The standard-library `net/http` package exposes the external API.

Implemented endpoints:

- `POST /api/v1/catalogs`
- `GET /api/v1/catalogs/{id}`
- `POST /api/v1/catalogs/{id}/elements`
- `GET /api/v1/catalogs/{id}/elements/{id}`
- `PATCH /api/v1/catalogs/{id}/elements/{id}`

Clients access application data through the API rather than internal packages or direct database access.

## Remaining work

The following capabilities are not yet confirmed as implemented:

- soft delete;
- platform-managed system attributes;
- dynamic metadata-driven fields;
- configuration-driven behavior beyond the current catalog slice.

They must be introduced deliberately and verified by tests.

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

New abstractions must be justified by an actual requirement and a verified vertical slice.
