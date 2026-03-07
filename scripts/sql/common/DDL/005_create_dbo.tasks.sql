IF OBJECT_ID(N'dbo.tasks', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.tasks (
        id          BIGINT             NOT NULL IDENTITY(1, 1),
        uuid        UNIQUEIDENTIFIER   NOT NULL,
        action_id   TINYINT            NOT NULL,
        status_id   TINYINT           NOT NULL DEFAULT 0,
        created_at  DATETIMEOFFSET(7)  NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        started_at  DATETIMEOFFSET(7)  NULL,
        scheduled_at DATETIMEOFFSET(7) NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        completed_at DATETIMEOFFSET(7) NULL,
        updated_at  DATETIMEOFFSET(7)  NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        error       NVARCHAR(500)      NULL,
        CONSTRAINT PK_tasks PRIMARY KEY CLUSTERED (id),
        CONSTRAINT FK_tasks_action FOREIGN KEY (action_id) REFERENCES dbo.actions (id),
        CONSTRAINT FK_tasks_status FOREIGN KEY (status_id) REFERENCES dbo.task_statuses (id)
    );
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCLU_tasks_uuid' AND object_id = OBJECT_ID(N'dbo.tasks'))
BEGIN
    CREATE UNIQUE NONCLUSTERED INDEX NCLU_tasks_uuid ON dbo.tasks (uuid);
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCL_tasks_action_status' AND object_id = OBJECT_ID(N'dbo.tasks'))
BEGIN
    CREATE NONCLUSTERED INDEX NCL_tasks_action_status ON dbo.tasks (action_id, status_id);
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCL_tasks_scheduled_at' AND object_id = OBJECT_ID(N'dbo.tasks'))
BEGIN
    CREATE NONCLUSTERED INDEX NCL_tasks_scheduled_at ON dbo.tasks (scheduled_at);
END
GO
