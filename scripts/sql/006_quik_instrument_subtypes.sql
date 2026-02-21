IF NOT EXISTS (
    SELECT * FROM sys.tables t
    INNER JOIN sys.schemas s ON t.schema_id = s.schema_id
    WHERE s.name = N'quik' AND t.name = N'instrument_subtypes'
)
BEGIN
    CREATE TABLE quik.instrument_subtypes (
        subtype_id       smallint  IDENTITY (1, 1), 
        type_id       smallint          ,
       	title nvarchar(150) NOT NULL,
        CONSTRAINT PK_quik_instrument_subtypes PRIMARY KEY CLUSTERED (subtype_id),
         CONSTRAINT FK_quik_instrument_type FOREIGN KEY (type_id) REFERENCES quik.instrument_types (type_id),
        CONSTRAINT UQ_quik_instrument_subtypes_type_title UNIQUE (type_id, title)
    );


END

GO
