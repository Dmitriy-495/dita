# DITA 0.2.0 — Checkpoint 002: Meta Core

## Состояние

**Итерация:** 002 — Meta Core
**Статус:** начата; аудит, архитектура и первый пакет Core Foundation подготовлены.
**Базовый commit:** `04d8352ef89c65ebccb997b36ce7e67353bf5be6`.
**Ветка:** `master`.

## Фактическое состояние

На GitHub master обнаружены Go 1.25.0 и pgx/v5 5.11.0; `cmd/dita/main.go`; legacy `internal/meta` Catalog; `internal/runtime` Catalog/Element; `internal/api/http`; `internal/storage/postgres`; миграции `000002_create_meta_catalogs` и `000003_create_catalog_elements`.

Старые миграции создают `public.meta_catalogs` и `public.meta_catalog_elements`; `code` элемента — `VARCHAR(255)`. Целевая модель DITA 0.2.0 использует схемы `meta`/`data`, отдельную таблицу на тип объекта и BIGINT бизнес-код. Старые таблицы нельзя напрямую переиспользовать без смены их смысла.

## Направление

1. Legacy-код и миграции сохраняются.
2. Новая модель — `internal/meta/configuration`.
3. Planner/Installer — отдельные доменные пакеты без прямой зависимости от pgx.
4. PostgreSQL persistence нового Meta — `internal/storage/postgres/meta`.
5. Новая схема добавляется аддитивной forward-only миграцией `000004_create_meta_core.up.sql`.
6. Сброс БД, DROP старых объектов и перенос данных не выполняются.

## Подготовленные артефакты и проверки

- Архитектура Iteration 002 и checklist.
- `internal/meta/configuration/model.go` — базовая JSON-модель стандартных типов и реквизитов.
- `internal/meta/configuration/model_test.go` — unit-тесты модели.
- `migrations/000004_create_meta_core.up.sql` — аддитивная миграция системной схемы Meta.
- Изолированный запуск `GO111MODULE=off go test ./internal/meta/configuration` завершился успешно на Go 1.23.2.
- `gofmt -d` не показал изменений.

Локальный checkout пользователя и PostgreSQL недоступны. Полный `go test ./...`, `go build ./cmd/dita` и SQL integration tests не запускались. Миграция не применялась и не проверялась на реальном PostgreSQL.

## Следующие шаги

Завершить Core Foundation, проверить Meta Schema в изолированной БД, реализовать Parser/Validator, dry-run Planner, Installer и integration tests. Обновить checkpoint фактическими результатами.

## Безопасность

Не выполнять `DROP DATABASE`, `DROP SCHEMA`, `TRUNCATE`, reset или автоматический перенос legacy-данных. Не запускать миграции, пока не подтверждено, что `DITA_DATABASE_URL` указывает на отдельную БД DITA, а не на `giga_watt` или иную БД GigaWatt.
