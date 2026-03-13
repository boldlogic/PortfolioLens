IF OBJECT_ID(N'dbo.external_system_urls', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.external_system_urls (
        id                INT            NOT NULL IDENTITY(1, 1),
        external_system_id TINYINT       NOT NULL,
        description       NVARCHAR(250)  NULL,
        proto             VARCHAR(10)    NOT NULL,
        host              NVARCHAR(150)  NOT NULL,
        port              INT            NULL,
        is_active         BIT            NOT NULL DEFAULT 1,
        CONSTRAINT PK_external_system_urls PRIMARY KEY CLUSTERED (id),
        CONSTRAINT FK_external_system_urls_ext_system FOREIGN KEY (external_system_id) REFERENCES dbo.external_systems (ext_system_id)
    );
END
GO
