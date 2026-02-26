IF OBJECT_ID (N'quik.boards', N'U') IS NULL 

BEGIN
CREATE TABLE
    quik.boards (
        board_id tinyint IDENTITY (1, 1),
        code nvarchar (15) NOT NULL,
        name nvarchar (60) NOT NULL,
        trade_point_id tinyint NULL,
        is_traded bit NOT NULL DEFAULT 0,
        CONSTRAINT PK_quik_boards PRIMARY KEY CLUSTERED (board_id),
        CONSTRAINT FK_quik_boards_trade_point_id FOREIGN KEY (trade_point_id) REFERENCES quik.trade_points (point_id),
    );

END 
GO 

IF OBJECT_ID (N'quik.boards', N'U') IS NOT NULL 
BEGIN
DROP INDEX IF EXISTS UQ_boards_code ON quik.boards;

CREATE UNIQUE NONCLUSTERED INDEX UQ_boards_code ON quik.boards (code);
END 

GO