UPDATE dbo.actions SET endpoint_id = 1 WHERE id = 1 AND endpoint_id IS NULL;
UPDATE dbo.actions SET endpoint_id = 2 WHERE id = 2 AND endpoint_id IS NULL;
UPDATE dbo.actions SET endpoint_id = 3 WHERE id = 3 AND endpoint_id IS NULL;
GO
