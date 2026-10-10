# DITA 0.1.0 — Iteration 002 Checklist

## Goal

Implement PostgreSQL persistence for catalog metadata and deliver the first working vertical slice through Runtime and HTTP API.

## Scope and results

- [x] Define the typed `Catalog` model.
- [x] Add PostgreSQL migration and constraints for catalog metadata.
- [x] Implement the catalog repository.
- [x] Support catalog creation and retrieval.
- [x] Implement persistence for catalog elements.
- [x] Generate element codes from catalog prefix and element ID.
- [x] Implement Runtime catalog loading.
- [x] Implement element creation, retrieval and update.
- [x] Expose catalog and element operations through HTTP API.
- [x] Return `404 Not Found` for missing catalogs and elements.
- [x] Add and run Go tests.
- [x] Run `gofmt`.
- [x] Run `go test ./...`.
- [x] Run `git diff --check`.
- [x] Verify the vertical slice manually against PostgreSQL.
- [x] Update architecture documentation.
- [x] Create `CHECKPOINT_002.md`.

## Manually verified scenario

1. Create catalog `Объекты эксплуатации`, prefix `ОО`.
2. Retrieve the catalog through Runtime.
3. Create element `Котельная №1`.
4. Verify generated code `ОО-000001`.
5. Retrieve the element.
6. Update its name and verify that its code remains unchanged.
7. Verify `404` responses for a missing catalog and a missing element.

## Deferred

- Soft delete.
- Platform-managed system attributes.
- Dynamic metadata-driven fields.
- Authentication and authorization.
- Additional configuration-driven runtime behavior.

## Completion criteria

The catalog metadata and elements persist in PostgreSQL, the first vertical slice works through HTTP API, and automated tests pass.
