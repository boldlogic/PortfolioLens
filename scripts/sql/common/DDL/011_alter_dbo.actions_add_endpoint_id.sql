IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID(N'dbo.actions') AND name = N'endpoint_id'
)
BEGIN
    ALTER TABLE dbo.actions
        ADD endpoint_id INT NULL
        CONSTRAINT FK_actions_endpoint FOREIGN KEY REFERENCES dbo.endpoints (id);
END
GO
