IF EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCLU_currencies_iso_char_code' AND object_id = OBJECT_ID(N'dbo.currencies'))
BEGIN
    DROP INDEX NCLU_currencies_iso_char_code ON dbo.currencies;
END
GO
