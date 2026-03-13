SET NOCOUNT ON;

DECLARE @action_id TINYINT = (SELECT id FROM dbo.actions WHERE code = N'currency.cb.fetch.historical_rates');
DECLARE @date_from NVARCHAR(10) = N'2026-03-01';
DECLARE @date_to   NVARCHAR(10) = N'2026-03-10';

DECLARE @tasks_with_cc TABLE (id BIGINT PRIMARY KEY, currency NVARCHAR(3) NOT NULL);

;WITH currencies(currency) AS (
    SELECT t.currency
    FROM (VALUES
        (N'AED'), (N'AMD'), (N'AUD'), (N'BYN'), (N'CAD'), (N'CHF'), (N'CNY'),
        (N'EUR'), (N'GBP'), (N'HKD'), (N'JPY'), (N'KGS'), (N'KZT'), (N'PLN'),
        (N'SEK'), (N'SGD'), (N'TJS'), (N'TRY'), (N'USD'), (N'UZS'), (N'ZAR')
    ) AS t(currency)
)
MERGE INTO dbo.tasks AS t
USING (SELECT currency FROM currencies) AS s
ON 1 = 0
WHEN NOT MATCHED BY TARGET THEN
    INSERT (uuid, action_id)
    VALUES (NEWID(), @action_id)
OUTPUT inserted.id, s.currency INTO @tasks_with_cc(id, currency);

INSERT INTO dbo.task_params (task_id, param_id, value)
SELECT id, 1, currency FROM @tasks_with_cc
UNION ALL
SELECT id, 2, @date_from FROM @tasks_with_cc
UNION ALL
SELECT id, 3, @date_to   FROM @tasks_with_cc;

SELECT id AS task_id, currency FROM @tasks_with_cc ORDER BY currency;
