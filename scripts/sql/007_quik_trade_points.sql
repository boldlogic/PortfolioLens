IF OBJECT_ID (N'quik.trade_points', N'U') IS NULL 

BEGIN
CREATE TABLE
    quik.trade_points (
        point_id tinyint IDENTITY (1, 1),
        code nvarchar (15) NOT NULL,
        name nvarchar (60) NOT NULL,
        CONSTRAINT PK_quik_trade_points PRIMARY KEY CLUSTERED (point_id),
    );

END 

GO 
IF OBJECT_ID (N'quik.trade_points', N'U') 
IS NOT NULL 

BEGIN
DROP INDEX IF EXISTS UQ_trade_points_code ON quik.trade_points;

CREATE UNIQUE NONCLUSTERED INDEX UQ_trade_points_code ON quik.trade_points (code);

END 
GO