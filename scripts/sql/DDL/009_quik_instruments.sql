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
        trade_point_id tinyint NOT NULL,  
        ticker              NVARCHAR(12)   NOT NULL,                    -- Код инструмента (char 15)
        registration_number NVARCHAR(30)   NULL,                     -- Рег.номер (char 250)
        full_name           NVARCHAR(100)   NULL,                    -- Инструмент (char 250)
        short_name          NVARCHAR(50)   NULL,                     -- Инструмент сокр. (char 100)
        isin                NVARCHAR(12)    NULL,                    -- ISIN (char 15)
        face_value          FLOAT       NULL,                    -- Номинал (float)
        maturity_date       DATE        NULL,                    -- Погашение (date)
        coupon_duration     INT         NULL,                    -- Длит. купона (int)
          
        rw          ROWVERSION      NOT NULL,
        
        
        CONSTRAINT PK_quik_instruments PRIMARY KEY CLUSTERED (instrument_id),
        CONSTRAINT FK_quik_instruments_trade_point_id FOREIGN KEY (trade_point_id) REFERENCES quik.trade_points (point_id),


    );


END
IF OBJECT_ID(N'quik.instruments', N'U') IS NOT NULL
BEGIN
    DROP INDEX IF EXISTS UQ_quik_instruments_ticker_trade_point_id
    ON quik.instruments;

    CREATE UNIQUE NONCLUSTERED INDEX UQ_quik_instruments_ticker_trade_point_id
    ON quik.instruments (ticker, trade_point_id);

END
GO


GO
