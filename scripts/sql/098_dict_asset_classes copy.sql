IF NOT EXISTS (
    SELECT * FROM sys.tables t
    INNER JOIN sys.schemas s ON t.schema_id = s.schema_id
    WHERE s.name = N'dict' AND t.name = N'asset_classes'
)
BEGIN
    CREATE TABLE dict.asset_classes (
        asset_id       smallint  IDENTITY (1, 1)          ,
       	title char(150) NOT NULL,
        CONSTRAINT PK_dict_asset_classes PRIMARY KEY CLUSTERED (asset_id),
        CONSTRAINT UQ_dict_asset_classes_title UNIQUE (title)
    );


END

GO
