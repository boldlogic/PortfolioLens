IF NOT EXISTS (
    SELECT * FROM sys.tables t
    INNER JOIN sys.schemas s ON t.schema_id = s.schema_id
    WHERE s.name = N'quik' AND t.name = N'instrument_types'
)
BEGIN
    CREATE TABLE quik.instrument_types (
        type_id       tinyint  IDENTITY (1, 1)          ,
       	title nvarchar(150) NOT NULL,
        CONSTRAINT PK_quik_instrument_types PRIMARY KEY CLUSTERED (type_id),
        CONSTRAINT UQ_quik_instrument_types_title UNIQUE (title)
    );


END

GO
