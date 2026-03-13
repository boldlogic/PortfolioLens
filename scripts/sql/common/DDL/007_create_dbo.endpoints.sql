IF OBJECT_ID(N'dbo.endpoints', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.endpoints (
        id                      INT           NOT NULL IDENTITY(1, 1),
        external_system_url_id  INT           NOT NULL,
        description             NVARCHAR(250) NULL,
        path                    NVARCHAR(250) NOT NULL,
        method                  VARCHAR(10)   NOT NULL,
        timeout_ms              INT           NOT NULL DEFAULT 20000,
        retry_policy            VARCHAR(20)   NOT NULL DEFAULT 'fixed',
        retry_count             TINYINT       NOT NULL DEFAULT 0,
        is_active               BIT           NOT NULL DEFAULT 1,
        CONSTRAINT PK_endpoints PRIMARY KEY CLUSTERED (id),
        CONSTRAINT FK_endpoints_ext_system_url FOREIGN KEY (external_system_url_id) REFERENCES dbo.external_system_urls (id)
    );
END
GO
