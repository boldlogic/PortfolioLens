IF OBJECT_ID(N'dbo.task_statuses', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.task_statuses (
        id   TINYINT      NOT NULL IDENTITY(0, 1),
        name NVARCHAR(50) NOT NULL,
        CONSTRAINT PK_task_statuses PRIMARY KEY CLUSTERED (id)
    );
END
GO
