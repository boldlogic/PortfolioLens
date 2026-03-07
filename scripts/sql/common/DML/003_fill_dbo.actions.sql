DELETE FROM dbo.actions;

SET IDENTITY_INSERT dbo.actions ON;

INSERT INTO dbo.actions (id, code, name)
SELECT t.id, t.code, t.name
FROM (VALUES
    (1, N'currency.cb.fetch.currency_list', N'Получение справочника валют из ЦБ по www.cbr.ru/scripts/XML_valFull.asp'),
    (2, N'currency.cb.fetch.rates_today', N'Получение текущих курсов валют из ЦБ по www.cbr.ru/scripts/XML_daily.asp'),
    (3, N'currency.cb.fetch.historical_rates', N'Получение истории курсов валют из ЦБ по www.cbr.ru/scripts/XML_dynamic.asp')
) AS t(id, code, name);

SET IDENTITY_INSERT dbo.actions OFF;
GO
