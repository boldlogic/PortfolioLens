IF OBJECT_ID(N'dbo.fx_cbr_rates', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.fx_cbr_rates (
        date                date              NOT NULL,
        quote_iso_code      INT               NOT NULL,
        base_iso_code       INT               NOT NULL,

        rate_quote_per_base      DECIMAL(18,8)     NULL,
        rate_base_per_quote      DECIMAL(18,8)     NULL,
        created_at          DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        updated_at          DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        ext_system_id   TINYINT           NULL,  -- кто последний обновил
        CONSTRAINT PK_fx_cbr_rates PRIMARY KEY CLUSTERED (date, quote_iso_code, base_iso_code),
        CONSTRAINT FK_fx_cbr_rates_ext_system FOREIGN KEY (ext_system_id) REFERENCES dbo.external_systems (ext_system_id)
    );
END
GO
