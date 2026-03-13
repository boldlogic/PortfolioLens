IF OBJECT_ID(N'dbo.actions', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.actions (
        id   TINYINT        NOT NULL IDENTITY(1, 1),
        code NVARCHAR(50)  NOT NULL,
        name NVARCHAR(150) NULL,
        CONSTRAINT PK_actions PRIMARY KEY CLUSTERED (id)
    );
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCLU_actions_code' AND object_id = OBJECT_ID(N'dbo.actions'))
BEGIN
    CREATE UNIQUE NONCLUSTERED INDEX NCLU_actions_code ON dbo.actions (code);
END
GO
