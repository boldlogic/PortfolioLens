SET IDENTITY_INSERT dbo.endpoint_params ON;

MERGE INTO dbo.endpoint_params AS tgt
USING (
    SELECT id, endpoint_id, param_id, external_name, param_location, ext_code_type_id, format, is_required, default_value
    FROM (VALUES
        (1,  1, NULL, 'User-Agent', 'header', NULL, NULL, 1, N'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 YaBrowser/25.12.0.0 Safari/537.36'),
        (2,  2, NULL, 'User-Agent', 'header', NULL, NULL, 1, N'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 YaBrowser/25.12.0.0 Safari/537.36'),
        (3,  3, 1,    'VAL_NM_RQ',  'query',  1,    NULL,         1, NULL), 
        (4,  3, 2,    'date_req1',  'query',  NULL, 'dd/MM/yyyy', 1, NULL),  
        (5,  3, 3,    'date_req2',  'query',  NULL, 'dd/MM/yyyy', 1, NULL), 
        (6,  3, NULL, 'User-Agent', 'header', NULL, NULL,         1, N'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 YaBrowser/25.12.0.0 Safari/537.36')
    ) AS t(id, endpoint_id, param_id, external_name, param_location, ext_code_type_id, format, is_required, default_value)
) AS src ON tgt.id = src.id
WHEN NOT MATCHED BY TARGET THEN
    INSERT (id, endpoint_id, param_id, external_name, param_location, ext_code_type_id, format, is_required, default_value)
    VALUES (src.id, src.endpoint_id, src.param_id, src.external_name, src.param_location, src.ext_code_type_id, src.format, src.is_required, src.default_value)
WHEN MATCHED THEN
    UPDATE SET
        endpoint_id      = src.endpoint_id,
        param_id         = src.param_id,
        external_name    = src.external_name,
        param_location   = src.param_location,
        ext_code_type_id = src.ext_code_type_id,
        format           = src.format,
        is_required      = src.is_required,
        default_value    = src.default_value;

SET IDENTITY_INSERT dbo.endpoint_params OFF;
GO
