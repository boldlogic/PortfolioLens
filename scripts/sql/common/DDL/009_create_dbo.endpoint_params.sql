IF OBJECT_ID(N'dbo.endpoint_params', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.endpoint_params (
        id              INT           NOT NULL IDENTITY(1, 1),
        endpoint_id     INT           NOT NULL,
        param_id        INT           NULL,      
        external_name   VARCHAR(50)   NOT NULL,  
        param_location  VARCHAR(10)   NOT NULL, 
        ext_code_type_id TINYINT      NULL,      
        format          VARCHAR(50)   NULL,    
        is_required     BIT           NOT NULL DEFAULT 1,
        default_value   NVARCHAR(500) NULL,     
        CONSTRAINT PK_endpoint_params PRIMARY KEY CLUSTERED (id),
        CONSTRAINT FK_endpoint_params_endpoint FOREIGN KEY (endpoint_id) REFERENCES dbo.endpoints (id),
        CONSTRAINT FK_endpoint_params_param    FOREIGN KEY (param_id)    REFERENCES dbo.params (id)
    );
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCLU_endpoint_params_endpoint_name' AND object_id = OBJECT_ID(N'dbo.endpoint_params'))
BEGIN
    CREATE UNIQUE NONCLUSTERED INDEX NCLU_endpoint_params_endpoint_name
        ON dbo.endpoint_params (endpoint_id, external_name);
END
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'NCL_endpoint_params_endpoint' AND object_id = OBJECT_ID(N'dbo.endpoint_params'))
BEGIN
    CREATE NONCLUSTERED INDEX NCL_endpoint_params_endpoint
        ON dbo.endpoint_params (endpoint_id)
        INCLUDE (param_id, external_name, param_location, ext_code_type_id, format, is_required, default_value);
END
GO
