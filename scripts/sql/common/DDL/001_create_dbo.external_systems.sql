IF OBJECT_ID (N'dict.external_systems', N'U') IS NULL BEGIN
CREATE TABLE
    dict.external_systems (
        external_system_id tinyint IDENTITY (1, 1),
        code varchar(12) NOT NULL,
        name varchar(128) CONSTRAINT PK_quik_firms PRIMARY KEY CLUSTERED (firm_id),
    );

END 
GO 
IF OBJECT_ID (N'dict.external_systems', N'U') IS NOT NULL 
BEGIN
DROP INDEX IF EXISTS UQ_firms_code ON dict.external_systems;

CREATE UNIQUE NONCLUSTERED INDEX UQ_firms_code ON dict.external_systems (code);

END 

GO