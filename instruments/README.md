# instruments

Сервис ведёт справочники торговых инструментов и площадок. Источник данных — таблица текущих котировок QUIK (`quik.current_quotes`), в которую терминал пишет данные через ODBC.

---

## Функции

| Компонент | Описание |
|-----------|----------|
| **Воркер инструментов** | Каждую секунду берёт строки из `quik.current_quotes` без привязки к инструменту, создаёт или обновляет записи в `quik.instruments` и `quik.instrument_boards` |
| **Воркер справочников** | Раз в 60 секунд синхронизирует типы и подтипы инструментов, борды и торговые площадки из котировок. Разметка бордов по площадкам (MOEX, SPB, OTC и другие) происходит здесь же |

---

## API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/boards` | Список бордов с торговыми площадками |
| GET | `/api/v1/boards/{id}` | Борд по ID |
| GET | `/api/v1/tradepoints` | Список торговых площадок |
| GET | `/api/v1/tradepoints/{id}` | Торговая площадка по ID |

Дополнительно:

- `/health` — health check
- `/metrics` — Prometheus-метрики

---

## Запуск

SQL Server с применёнными DDL из `scripts/sql/` и конфигурационный файл `instruments/internal/configs/config.yaml`. Пароль и пользователя БД — через `DB_PASSWORD` и `DB_USER`, см. [`pkg/dbzap/db.go`](../pkg/dbzap/db.go).

Из корня репозитория:

```bash
go run ./instruments/cmd
go run ./instruments/cmd -config /path/to/config.yaml
```

---

## Конфигурация

Пример `config.yaml`:

```yaml
log:
  level: "info"
  format: "json"
  output_file: "app.log"
server:
  port: 8002
  timeout: 20
db:
  host: "localhost"
  port: 1433
  user: "user"
  password: "***"
  db_name: "db_name"
```

Путь по умолчанию: `instruments/internal/configs/config.yaml`.

---

## Объекты БД

Все таблицы находятся в схеме `quik`.

| Таблица | Описание |
|---------|----------|
| `quik.current_quotes` | Источник данных (пишет QUIK через ODBC) |
| `quik.instruments` | Торговые инструменты |
| `quik.instrument_boards` | Связь инструмента с бордом, типом и валютами |
| `quik.instrument_types` | Типы инструментов (Акции, Облигации, ...) |
| `quik.instrument_subtypes` | Подтипы инструментов |
| `quik.boards` | Торговые борды (режимы торгов) |
| `quik.trade_points` | Торговые площадки (MOEX, SPB, ...) |

---

## Связанные сервисы

| Сервис | Роль |
|--------|------|
| **quik-portfolio** | Читает борды и торговые площадки при пересчёте портфеля |
