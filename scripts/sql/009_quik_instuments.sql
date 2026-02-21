IF NOT EXISTS (SELECT * FROM sys.schemas WHERE name = N'quik')
    EXEC (N'CREATE SCHEMA quik');

GO

IF NOT EXISTS (
    SELECT * FROM sys.tables t
    INNER JOIN sys.schemas s ON t.schema_id = s.schema_id
    WHERE s.name = N'quik' AND t.name = N'instruments'
)
BEGIN
    CREATE TABLE quik.instruments (
        instrument_id       BIGINT   IDENTITY (1, 1), 
        ticker              NVARCHAR(12)   NOT NULL,                    -- Код инструмента (char 15)
        registration_number NVARCHAR(30)   NULL,                     -- Рег.номер (char 250)
        full_name           NVARCHAR(100)   NULL,                    -- Инструмент (char 250)
        short_name          NVARCHAR(50)   NULL,                     -- Инструмент сокр. (char 100)
        --class_code          NVARCHAR(15)    NULL,                    -- Код класса (char 20)
       -- class_name          NVARCHAR(60)   NULL,                     -- Класс (char 200)
        isin                NVARCHAR(12)    NULL,                    -- ISIN (char 15)
        face_value          FLOAT       NULL,                    -- Номинал (float)
        base_currency       bigint    NULL,                    -- Валюта / Базовая валюта
        quote_currency      bigint     NULL,                    -- Котир.валюта
        counter_currency    bigint    NULL,                    -- Сопр.валюта
        maturity_date       DATE        NULL,                    -- Погашение (date)
        coupon_duration     INT         NULL,                    -- Длит. купона (int)

        type_id       smallint NOT NULL,
        subtype_id       smallint NULL, 
        rw          ROWVERSION      NOT NULL,
        
        CONSTRAINT PK_quik_instruments PRIMARY KEY CLUSTERED (instrument_id),
        CONSTRAINT FK_quik_instruments_base_currency FOREIGN KEY (base_currency) REFERENCES dbo.currencies (iso_code),
        CONSTRAINT FK_quik_instruments_quote_currency FOREIGN KEY (quote_currency) REFERENCES dbo.currencies (iso_code),
        CONSTRAINT FK_quik_instruments_counter_currency FOREIGN KEY (counter_currency) REFERENCES dbo.currencies (iso_code),

    );


END
IF OBJECT_ID(N'quik.instruments', N'U') IS NOT NULL
BEGIN
    DROP INDEX IF EXISTS UQ_quik_instruments_ticker
    ON quik.instruments;

    CREATE UNIQUE NONCLUSTERED INDEX UQ_quik_instruments_ticker
    ON quik.instruments (ticker);

END
GO


GO
