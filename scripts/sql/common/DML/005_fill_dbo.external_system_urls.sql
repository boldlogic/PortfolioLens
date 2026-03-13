SET IDENTITY_INSERT dbo.external_system_urls ON;

MERGE INTO dbo.external_system_urls AS tgt
USING (
    SELECT id, external_system_id, description, proto, host, port, is_active
    FROM (VALUES
        (1, 1, N'ЦБ РФ — основной адрес', 'https', N'www.cbr.ru', NULL, 1)
    ) AS t(id, external_system_id, description, proto, host, port, is_active)
) AS src ON tgt.id = src.id
WHEN NOT MATCHED BY TARGET THEN
    INSERT (id, external_system_id, description, proto, host, port, is_active)
    VALUES (src.id, src.external_system_id, src.description, src.proto, src.host, src.port, src.is_active)
WHEN MATCHED THEN
    UPDATE SET
        external_system_id = src.external_system_id,
        description        = src.description,
        proto              = src.proto,
        host               = src.host,
        port               = src.port,
        is_active          = src.is_active;

SET IDENTITY_INSERT dbo.external_system_urls OFF;
GO
