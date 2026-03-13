IF OBJECT_ID(N'dbo.params', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.params (
        id        INT           NOT NULL IDENTITY(1, 1),
        code      VARCHAR(50)   NOT NULL,
        name      NVARCHAR(150) NULL,
        data_type VARCHAR(20)   NOT NULL,  -- string, date
        CONSTRAINT PK_params PRIMARY KEY CLUSTERED (id)
    );
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCLU_params_code' AND object_id = OBJECT_ID(N'dbo.params'))
BEGIN
    CREATE UNIQUE NONCLUSTERED INDEX NCLU_params_code ON dbo.params (code);
END
GO
