# task-manager

Сервис принимает задачи по API и записывает их в SQL Server очередь. Обработчики задач живут в отдельных сервисах - task-manager отвечает только за постановку.

## API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/v1/tasks` | Создать задачу |

Дополнительно:

- `/health` - health check
- `/metrics` - Prometheus-метрики

### Пример запроса

```json
POST /api/v1/tasks
{
    "action": "currency.cb.fetch.historical_rates",
    "params": {
        "ccyCode": "USD",
        "dateFrom": "1991-01-01",
        "dateTo": "2026-03-21"
    }
}
```

### Коды ответов

| Код | Ситуация |
|-----|----------|
| 201 | Задача создана |
| 400 | Неизвестный `action_code` или невалидный параметр |
| 409 | Задача с таким UUID уже существует |
| 500 | Ошибка записи в БД |

---

## Запуск

SQL Server с применёнными DDL из `scripts/sql/` и конфигурационный файл `task-manager/internal/configs/config.yaml`. Пароль и пользователя БД — через `DB_PASSWORD` и `DB_USER`, см. [`pkg/dbzap/db.go`](../pkg/dbzap/db.go).

Из корня репозитория:

```bash
go run ./task-manager/cmd
go run ./task-manager/cmd -config /path/to/config.yaml
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
  port: 8001
  timeout: 20
db:
  host: "localhost"
  port: 1433
  user: "user"
  password: "***"
  db_name: "db_name"
```

Путь по умолчанию: `task-manager/internal/configs/config.yaml`.

---

## Объекты БД

Все таблицы находятся в схеме `dbo`.

| Таблица | Описание |
|---------|----------|
| `dbo.tasks` | Очередь задач. Атомарная выборка |
| `dbo.task_params` | Параметры задачи (code - value) |
| `dbo.actions` | Справочник типов действий (`action_code`) |
| `dbo.params` | Справочник допустимых параметров |
