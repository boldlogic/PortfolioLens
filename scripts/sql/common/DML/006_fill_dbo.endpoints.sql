SET IDENTITY_INSERT dbo.endpoints ON;

MERGE INTO dbo.endpoints AS tgt
USING (
    SELECT id, external_system_url_id, description, path, method, timeout_ms, retry_policy, retry_count, is_active
    FROM (VALUES
        (1, 1, N'Справочник валют ЦБ',          'scripts/XML_valFull.asp', 'GET', 20000, 'fixed', 0, 1),
        (2, 1, N'Текущие курсы валют ЦБ',        'scripts/XML_daily.asp',   'GET', 20000, 'fixed', 0, 1),
        (3, 1, N'Исторические курсы валют ЦБ',   'scripts/XML_dynamic.asp', 'GET', 20000, 'fixed', 0, 1)
    ) AS t(id, external_system_url_id, description, path, method, timeout_ms, retry_policy, retry_count, is_active)
) AS src ON tgt.id = src.id
WHEN NOT MATCHED BY TARGET THEN
    INSERT (id, external_system_url_id, description, path, method, timeout_ms, retry_policy, retry_count, is_active)
    VALUES (src.id, src.external_system_url_id, src.description, src.path, src.method, src.timeout_ms, src.retry_policy, src.retry_count, src.is_active)
WHEN MATCHED THEN
    UPDATE SET
        external_system_url_id = src.external_system_url_id,
        description            = src.description,
        path                   = src.path,
        method                 = src.method,
        timeout_ms             = src.timeout_ms,
        retry_policy           = src.retry_policy,
        retry_count            = src.retry_count,
        is_active              = src.is_active;

SET IDENTITY_INSERT dbo.endpoints OFF;
GO
