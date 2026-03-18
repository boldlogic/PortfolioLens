# quik-portfolio

Сервис хранит лимиты из QUIK-терминала (денежные, по бумагам, OTC) и предоставляет портфель с пересчётом стоимости позиций в целевую валюту по курсам ЦБ из `market-data-currency`.

---

## Функции

| Компонент | Описание |
|-----------|----------|
| **HTTP API** | Приём лимитов (POST), выдача по дате (GET), портфель с FX-пересчётом, справочник фирм |
| **Roll-forward** | Перенос лимитов на текущую дату при пропуске торговых дней (воркеры: money, securities, otc) |
| **Actualize firms** | Синхронизация справочника фирм из загруженных лимитов |

---

## API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/quik/limits/money?date=` | Лимиты по денежным средствам на дату |
| POST | `/quik/limits/money` | Добавить лимит по денежным средствам |
| GET | `/quik/limits/securities?date=` | Лимиты по бумагам на дату |
| POST | `/quik/limits/securities` | Добавить лимит по бумагам |
| GET | `/quik/limits/securities/otc?date=` | OTC-лимиты по бумагам на дату |
| POST | `/quik/limits/securities/otc` | Добавить OTC-лимит |
| GET | `/quik/limits?date=` | Сводные лимиты (деньги + бумаги) на дату |
| GET | `/quik/portfolio?targetCcy=` | Портфель с FX-пересчётом в целевую валюту (по умолчанию RUB) |
| GET | `/quik/firms` | Список фирм |
| POST | `/quik/firms` | Добавить фирму |
| GET | `/quik/firms/{id}` | Фирма по ID |
| PATCH | `/quik/firms/{id}` | Обновить наименование фирмы |

---

## Запуск

**Предварительные требования:**
- SQL Server с применёнными DDL из `scripts/sql/quik-portfolio/`
- Конфигурационный файл `quik-portfolio/internal/configs/config.yaml`

Пароль БД можно передать через переменную окружения `DB_PASSWORD` (см. `pkg/config`).

Из корня репозитория:

```bash
go run ./quik-portfolio/cmd
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
  port: 8000
  timeout: 20
db:
  host: "localhost"
  port: 1433
  user: "user"
  password: "***"
  db_name: "db_name"
```

Путь по умолчанию: `quik-portfolio/internal/configs/config.yaml`.

---

## Объекты БД

Все таблицы находятся в схеме `quik`.

### quik.security_limits — лимиты по бумагам

| Поле | Тип | Ключ | По умолчанию | Описание |
|------|-----|------|--------------|----------|
| load_date | date | PK | getdate() | Дата загрузки лимита |
| client_code | varchar(12) | PK | — | Код клиента |
| ticker | varchar(12) | PK | — | Код инструмента |
| trade_account | varchar(12) | PK | — | Торговый счёт |
| settle_code | varchar(5) | PK | `'Tx'` | Код расчётов |
| firm_code | varchar(12) | PK | — | Код фирмы |
| firm_name | varchar(128) | — | NULL | Наименование фирмы |
| balance | DECIMAL(19,4) | — | NULL | Остаток в лотах или штуках |
| acquisition_ccy | varchar(4) | — | NULL | Валюта приобретения |
| isin | varchar(12) | — | NULL | ISIN |
| ts | timestamp | — | — | Метка времени строки |

PK: `(load_date, client_code, ticker, trade_account, settle_code, firm_code)`

### quik.money_limits — лимиты по денежным средствам

| Поле | Тип | Ключ | По умолчанию | Описание |
|------|-----|------|--------------|----------|
| load_date | date | PK | getdate() | Дата загрузки лимита |
| client_code | varchar(12) | PK | — | Код клиента |
| ccy | varchar(4) | PK | — | Код валюты |
| position_code | varchar(4) | PK | — | Код позиции |
| settle_code | varchar(5) | PK | — | Код расчётов |
| firm_code | varchar(12) | PK | — | Код фирмы |
| firm_name | varchar(128) | — | NULL | Наименование фирмы |
| balance | DECIMAL(19,4) | — | NULL | Остаток денежных средств |
| ts | timestamp | — | — | Метка времени строки |

PK: `(load_date, client_code, ccy, position_code, settle_code, firm_code)`

### quik.security_limits_otc — OTC-лимиты по бумагам

| Поле | Тип | Ключ | По умолчанию | Описание |
|------|-----|------|--------------|----------|
| load_date | date | PK | getdate() | Дата загрузки лимита |
| client_code | varchar(12) | PK | — | Код клиента |
| ticker | varchar(12) | PK | — | Код инструмента |
| trade_account | varchar(12) | PK | `'OTC'` | Торговый счёт |
| settle_code | varchar(5) | PK | `'Tx'` | Код расчётов |
| firm_code | varchar(12) | PK | — | Код фирмы |
| firm_name | varchar(128) | — | NULL | Наименование фирмы |
| balance | DECIMAL(19,4) | — | NULL | Остаток |
| acquisition_ccy | varchar(3) | — | NULL | Валюта приобретения |
| isin | varchar(12) | — | NULL | ISIN |
| ts | timestamp | — | — | Метка времени строки |

PK: `(load_date, client_code, ticker, trade_account, settle_code, firm_code)`

### quik.firms — справочник фирм

| Поле | Тип | Ключ | По умолчанию | Описание |
|------|-----|------|--------------|----------|
| firm_id | tinyint IDENTITY | PK | — | Внутренний Id |
| code | varchar(12) | UNIQUE | — | Код участника торгов |
| name | varchar(128) | — | — | Наименование фирмы |

PK: `(firm_id)`. Уникальный некластеризованный индекс по `code`.

---

## Настройка экспорта из QUIK

Терминал QUIK поддерживает прямой экспорт данных в таблицы БД через ODBC. Лимиты по бумагам и денежные лимиты загружаются автоматически при каждом обновлении таблицы в терминале.

> **OTC-лимиты QUIK не экспортирует.** Их необходимо загружать вручную через API (`POST /quik/limits/securities/otc`).

### Предварительная настройка

1. Применить DDL-скрипты из `scripts/sql/quik-portfolio/` к целевой БД.
2. Создать ODBC-источник данных (DSN) для подключения к этой БД.
3. Убедиться, что пользователь ODBC имеет право на INSERT и UPDATE в схему `quik`.

### 1. Экспорт лимитов по бумагам (`quik.security_limits`)

Источник данных — **Таблица лимитов по бумагам** (аналог DepoLimits в МЭБИ QUIK).

**Шаги настройки в терминале QUIK:**

1. Открыть меню **Создать окно → Позиции по инструментам** — откроется Таблица лимитов по бумагам.
2. ПКМ по таблице → **Редактировать таблицу** → добавить параметры согласно колонке «Поле QUIK» таблицы ниже → **Да**.
3. ПКМ по таблице → **Вывод по ODBC**.
4. Выбрать созданный ранее DSN.
5. В поле **Таблица** выбрать `security_limits`.
6. Для каждого параметра QUIK назначить целевое поле согласно таблице ниже → **Да**.

| Поле QUIK | Целевое поле | Пример |
|-----------|-------------|--------|
| *(дата загрузки)* | `load_date` | 2026-03-17 |
| Код клиента | `client_code` | 11A2BC |
| Код инструмента | `ticker` | OBLG |
| Счет депо | `trade_account` | L01-00000F00 |
| Срок расчётов | `settle_code` | T0 |
| Фирма | `firm_code` | MC0003300000 |
| Наименование фирмы | `firm_name` | ВТБ |
| Баланс | `balance` | 500.0000 |
| Валюта цены приобретения | `acquisition_ccy` | SUR |
| ISIN | `isin` | RU000A1002S8 |

### 2. Экспорт лимитов по денежным средствам (`quik.money_limits`)

Источник данных — **Таблица лимитов по денежным средствам** (аналог MoneyLimits в МЭБИ QUIK).

**Шаги настройки в терминале QUIK:**

1. Открыть меню **Создать окно → Позиции по деньгам** — откроется Таблица лимитов по денежным средствам.
2. ПКМ по таблице → **Редактировать таблицу** → добавить параметры согласно колонке «Поле QUIK» таблицы ниже →  **Да**.
3. ПКМ по таблице → **Вывод по ODBC**.
4. Выбрать созданный ранее DSN.
5. В поле **Таблица** выбрать `money_limits`.
6. Для каждого параметра QUIK назначить целевое поле согласно таблице ниже → **Да**.

| Поле QUIK          | Целевое поле    | Пример       |
| ------------------ | --------------- | ------------ |
| *(дата загрузки)*  | `load_date`     | 2026-03-17   |
| Код клиента        | `client_code`   | 11A2BC       |
| Валюта             | `ccy`           | SUR          |
| Код позиции        | `position_code` | EQTV         |
| Срок расчётов      | `settle_code`   | T0           |
| Фирма              | `firm_code`     | MC0003300000 |
| Наименование фирмы | `firm_name`     | ВТБ          |
| Баланс             | `balance`       | 100.0000     |
