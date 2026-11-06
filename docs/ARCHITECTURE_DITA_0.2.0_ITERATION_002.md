# DITA 0.2.0 — Iteration 002: Meta Core Architecture

## 1. Цель

Реализовать минимальное системное ядро Meta: реестр конфигураций и версий, метаданные типов/реквизитов, строгую обработку JSON, dry-run планирование и контролируемую установку типизированной схемы PostgreSQL. GigaWatt остаётся отдельной прикладной конфигурацией.

## 2. Фактическая исходная точка

Проверенный GitHub commit: `04d8352ef89c65ebccb997b36ce7e67353bf5be6` (`docs: complete DITA 0.2.0 iteration 001`).

В репозитории обнаружены `cmd/dita/main.go`, legacy `internal/meta` Catalog, `internal/runtime` Catalog/Element, `internal/api/http`, `internal/storage/postgres` и миграции `000002_create_meta_catalogs` / `000003_create_catalog_elements`. Старые таблицы находятся в `public`; `meta_catalog_elements.code` — `VARCHAR(255)`. Новая модель требует `meta`/`data`, отдельные таблицы типов и BIGINT бизнес-код. Старый код и миграции сохраняются на переходном этапе.

## 3. Решения итерации

1. Не переписывать legacy `internal/meta`, `internal/runtime`, HTTP API и миграции `000002`/`000003` на первом шаге.
2. Новая JSON-модель — `internal/meta/configuration`.
3. Dry-run планирование — `internal/meta/planner`.
4. Применение плана — `internal/meta/installer`.
5. PostgreSQL persistence нового Meta — `internal/storage/postgres/meta`.
6. Доменная модель, parser, validator и planner не зависят от `pgxpool` и не исполняют SQL.
7. Существующий migration runner и `public.schema_migrations` сохраняются.
8. Базовая схема добавляется forward-only миграцией `000004_create_meta_core.up.sql`; она только создаёт `meta`, `data` и системные таблицы, не удаляя старые объекты.
9. Конфигурация декодируется и проверяется до DDL. Installer применяет проверенный план транзакционно под advisory lock.
10. Сброс БД, DROP старых таблиц и перенос данных не входят в эту итерацию.

## 4. Целевая структура

```text
cmd/dita/main.go                 # composition root
internal/api/http/               # текущий transport, пока без изменений
internal/meta/                   # legacy Catalog API, временно сохраняется
internal/meta/configuration/     # JSON-модель, parser, validator
internal/meta/planner/           # desired vs installed, dry-run
internal/meta/installer/         # применение проверенного плана
internal/runtime/                # legacy Runtime, позже переводится на Meta Core
internal/storage/postgres/       # существующий pool/migration runner
internal/storage/postgres/meta/  # repositories нового Meta Core
migrations/000004_create_meta_core.up.sql
```

Зависимости направлены от инфраструктуры к доменным интерфейсам. `configuration`, `planner` и `installer` не должны импортировать `pgxpool`; SQL выполняет инфраструктурный слой.

## 5. Первая схема Meta

Создаются `meta.actors`, `meta.configurations`, `meta.configuration_versions`, `meta.object_types`, `meta.attributes`, `meta.number_counters` и schema `data`. Исходный JSONB версии — канонический документ; нормализованные типы и реквизиты — индекс для Runtime/DDL. Для Catalog в `number_counters` используется `year = 0`, для Document/Planning — календарный год. Это внутренняя деталь, не публичная generic numbering subsystem.

## 6. Порядок реализации

| Этап | Результат |
|---|---|
| 002.1 Core Foundation | Доменная JSON-модель и базовые типы |
| 002.2 Meta Schema | Аддитивная миграция системных таблиц |
| 002.3 Parser | Строгое декодирование, неизвестные поля и дубли ключей отклоняются |
| 002.4 Validator | UUID, SemVer, виды, типы, ограничения и ссылки |
| 002.5 Planner | Детерминированный dry-run без SQL |
| 002.6 Installer | Транзакционный DDL под advisory lock |
| 002.7 Verification | Unit/integration tests, повторная установка, FK и конкурентные сценарии |

Идём последовательно; не создавать абстракции для дальних этапов раньше времени.

## 7. Безопасность

- Старые миграции `000002`/`000003` не изменять.
- `000004` только добавляет schema/table/index/constraint.
- Не применять миграцию, пока не подтверждено, что `DITA_DATABASE_URL` указывает на отдельную БД DITA.
- Не выполнять `DROP DATABASE`, `DROP SCHEMA`, `TRUNCATE`, reset или автоматический перенос legacy-данных.
- SQL identifiers строить только из проверенных UUID; значения передавать параметрами.
- Destructive changes не генерировать автоматически.

## 8. Критерии завершения

Строгая конфигурация декодируется и валидируется до DDL; Meta schema воспроизводима; planner не меняет БД; installer применяет только проверенный план; повторная установка идемпотентна; тесты подтверждают ограничения и транзакционность. Выполнить `gofmt`, `go test ./...`, `go build ./cmd/dita`, `git diff --check` и записать фактические результаты.

## 9. Не входит в объём

IAM, пароли, сессии, роли/права, аудит, регистры, workflow, табличные части, вложения, EAV, GigaWatt-специфические типы, импорт legacy-данных и удаление legacy-кода.
