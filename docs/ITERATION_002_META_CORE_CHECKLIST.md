# DITA 0.2.0 — Iteration 002 Checklist: Meta Core

## Аудит
- [x] Проверен актуальный commit `04d8352` на `master`.
- [x] Проверены `go.mod`, `cmd/dita/main.go`, legacy Meta/Runtime, HTTP и PostgreSQL repositories.
- [x] Проверены миграции `000002` и `000003`.
- [x] Зафиксировано несовпадение legacy-схемы с целевой моделью DITA 0.2.0.
- [x] Принято решение сохранять legacy-код и миграции на переходном этапе.

## 002.1 Core Foundation
- [x] Подготовить доменную модель JSON-конфигурации в `internal/meta/configuration/model.go`.
- [x] Определить типы стандартных объектов и реквизитов.
- [x] Добавить unit-тесты модели.
- [x] Проверить `gofmt` и компиляцию пакета изолированно (`go test`).

## 002.2 Meta Schema
- [x] Подготовить `000004_create_meta_core.up.sql` (проверка на PostgreSQL ещё не выполнена).
- [ ] Создать `meta` и `data`.
- [ ] Создать actors, configurations, configuration_versions, object_types, attributes, number_counters.
- [ ] Добавить PK/FK/UNIQUE/CHECK и индексы.
- [ ] Проверить на отдельной тестовой БД DITA и повторное применение runner.

## 002.3 Parser
- [ ] Строгое UTF-8 JSON decoding.
- [ ] Отклонение неизвестных полей и повторяющихся ключей.
- [ ] Проверка format_version.
- [ ] Позитивные и негативные parser tests.

## 002.4 Validator
- [ ] UUID и SemVer.
- [ ] Уникальность ID/code.
- [ ] Стандартные виды и типы реквизитов.
- [ ] Применимость type-specific параметров.
- [ ] Length, precision/scale, default и ограничения.
- [ ] Ссылки и принадлежность target той же конфигурации.
- [ ] Индексы и ссылки на существующие атрибуты.
- [ ] Детерминированные ошибки с JSON path.

## 002.5 Planner
- [ ] Сравнение desired configuration с установленной.
- [ ] Классификация безопасных/условных/разрушительных изменений.
- [ ] Dry-run без SQL.
- [ ] Блокировка неизвестных изменений.
- [ ] Детерминированный порядок действий.

## 002.6 Installer
- [ ] Advisory lock по configuration UUID.
- [ ] Одна транзакция на установку.
- [ ] Сначала object tables, затем FK.
- [ ] Индексы/ограничения только из проверенных метаданных.
- [ ] Запись JSONB version и checksum.
- [ ] installed_version обновляется последним.
- [ ] Destructive DDL только по отдельному подтверждённому плану.

## 002.7 Verification
- [ ] Unit tests model/parser/validator/planner.
- [ ] PostgreSQL integration tests.
- [ ] Идемпотентность повторной установки.
- [ ] FK и циклические ссылки.
- [ ] Конкурентная установка и нумерация.
- [ ] `gofmt`, `go test ./...`, `go build ./cmd/dita`, `git diff --check`.
- [ ] Обновить checkpoint фактическими результатами.
- [ ] Опубликовать отдельным коммитом после проверки.

## Ограничения

Не изменять `000002`/`000003`; не выполнять DROP/RESET/пересоздание БД; не запускать миграции на неизвестной БД; не удалять legacy Catalog/Runtime; не добавлять IAM, audit, workflow, registers, EAV или GigaWatt-специфику.
