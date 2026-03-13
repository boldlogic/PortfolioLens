SET IDENTITY_INSERT dbo.params ON;

MERGE INTO dbo.params AS tgt
USING (
    SELECT id, code, name, data_type
    FROM (VALUES
        (1, 'char_code',  N'Символьный код',          'string'),
        (2, 'date_from',  N'Начальная дата периода',   'date'),
        (3, 'date_to',    N'Конечная дата периода',    'date')
    ) AS t(id, code, name, data_type)
) AS src ON tgt.id = src.id
WHEN NOT MATCHED BY TARGET THEN
    INSERT (id, code, name, data_type)
    VALUES (src.id, src.code, src.name, src.data_type)
WHEN MATCHED THEN
    UPDATE SET
        code      = src.code,
        name      = src.name,
        data_type = src.data_type;

SET IDENTITY_INSERT dbo.params OFF;
GO
