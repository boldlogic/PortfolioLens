DELETE FROM dbo.task_statuses;

SET IDENTITY_INSERT dbo.task_statuses ON;

INSERT INTO dbo.task_statuses (id, name)
SELECT t.id, t.name
FROM (VALUES
    (0, N'scheduled'),
    (1, N'in_progress'),
    (2, N'completed'),
    (3, N'error')
) AS t(id, name);

SET IDENTITY_INSERT dbo.task_statuses OFF;
GO
