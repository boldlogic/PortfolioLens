IF OBJECT_ID(N'dbo.task_params', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.task_params (
        task_id  BIGINT        NOT NULL,
        param_id INT           NOT NULL,
        value    NVARCHAR(250) NOT NULL,
        CONSTRAINT PK_task_params PRIMARY KEY CLUSTERED (task_id, param_id),
        CONSTRAINT FK_task_params_task  FOREIGN KEY (task_id)  REFERENCES dbo.tasks (id),
        CONSTRAINT FK_task_params_param FOREIGN KEY (param_id) REFERENCES dbo.params (id)
    );
END
GO
