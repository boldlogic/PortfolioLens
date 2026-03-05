IF OBJECT_ID(N'dbo.currencies', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.currencies (
        iso_code       SMALLINT          NOT NULL,
        iso_char_code  NVARCHAR(3)       NOT NULL,
        currency_name  NVARCHAR(100)     NULL,
        lat_name       NVARCHAR(100)     NULL,
        minor_units                      INT               NULL,
        created_at                       DATETIMEOFFSET(7) NOT NULL DEFAULT (getdate ()),
        updated_at                       DATETIMEOFFSET(7) NOT NULL DEFAULT (getdate ()),
        ext_system_id   TINYINT           NULL,  -- кто последний обновил
        CONSTRAINT PK_currencies PRIMARY KEY CLUSTERED (iso_code),
        CONSTRAINT FK_currencies_ext_system FOREIGN KEY (ext_system_id) REFERENCES dbo.external_systems (ext_system_id)
    );
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCLU_currencies_iso_char_code' AND object_id = OBJECT_ID(N'dbo.currencies'))
BEGIN
    CREATE UNIQUE NONCLUSTERED INDEX NCLU_currencies_iso_char_code ON dbo.currencies (iso_char_code);
END
GO
