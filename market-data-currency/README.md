# Market Data Currency

Сервис получает справочник валют и курсы из ЦБ РФ по API и терминала QUIK (чтение из зеркала таблицы Текущие торги в БД), сохраняет их в БД. 
В монорепо Portfolio Lens он отвечает за рыночные данные по валютам: portfolio берёт отсюда курсы для пересчёта позиций в рубли.

---

## Функции

| Компонент | Описание |
|-----------|----------|
| **Воркер задач** | Берёт задачи из `dbo.tasks`, по настройкам эндпоинтов (`dbo.endpoints`, `dbo.actions`) строит HTTP-запрос, вызывает API ЦБ, парсит XML и сохраняет курсы и справочник в БД |
| **Воркер QUIK** | Раз в минуту сливает кросс-курсы из `quik.current_quotes` в `dbo.fx_cbr_rates` |
| **HTTP API** | `GET /api/v1/currencies` отдаёт справочник валют, `GET /api/v1/currencies/{code}` — валюту по ISO 4217 |

---

## Типы задач (actions)

- `currency.cb.fetch.currency_list` — справочник валют ЦБ
- `currency.cb.fetch.rates_today` — курсы на текущую дату или на завтра (в зависимости от времени выполнения запроса)
- `currency.cb.fetch.historical_rates` — история курсов (в параметрах нужен `char_code`)

Задачи создаёт task-manager по запросу от клиента.

---

## Запуск

Нужен SQL Server с применёнными миграциями из `scripts/sql/` (common и market-data-currency) и конфиг в `market-data-currency/internal/configs/config.yaml`. Пароль БД можно задать через `DB_PASSWORD` (см. `pkg/config`).

Из корня репозитория:

```bash
go run ./market-data-currency/cmd
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
  port: 6000
  timeout: 20
db:
  host: "localhost"
  port: 1433
  user: "user"
  password: "***"
  db_name: "db_name"
client:
  timeout: 30
```

Путь по умолчанию: `market-data-currency/internal/configs/config.yaml`.

---


## Связанные сервисы

| Сервис | Роль |
|--------|------|
| **task-manager** | Создаёт задачи в `dbo.tasks` (POST /add_task) |
| **quik-portfolio** | Портфель и лимиты, читает `fx_cbr_rates` для пересчёта в рубли |
