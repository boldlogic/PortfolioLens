WITH src AS (
    SELECT code, name
    FROM (VALUES
        ('MOEX', N'Основной рынок МБ'),
        ('SPB', N'Основной рынок СпбБ'),
        ('MOEX_OTC', N'Внебиржевой рынок МБ'),
        ('SPB_OTC', N'Внебиржевой рынок СпбБ'),
        ('BCS_OTC', N'Внебиржевые торги брокера БКС'),
        ('T_OTC', N'Внебиржевые торги брокера Т-банк'),
        ('VTB_OTC', N'Внебиржевые торги брокера ВТБ'),
        ('FORTS', N'Срочный рынок МБ'),
        ('SPB_BLK', N'Неторговый раздел СпбБ')
    ) AS t(code, name)
)
MERGE INTO quik.trade_points AS tgt
USING src ON tgt.code = src.code
WHEN MATCHED AND tgt.name <> src.name THEN
    UPDATE SET tgt.name = src.name
WHEN NOT MATCHED BY TARGET THEN
    INSERT (code, name)
    VALUES (src.code, src.name);

GO
