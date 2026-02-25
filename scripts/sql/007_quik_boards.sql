IF OBJECT_ID(N'quik.boards', N'U') IS NULL
BEGIN
    CREATE TABLE quik.boards (
        board_id       smallint  IDENTITY (1, 1)          ,  
       	code nvarchar(15) NOT NULL,
        name            nvarchar(60)   NOT NULL,
                CONSTRAINT PK_quik_boards PRIMARY KEY CLUSTERED (board_id),

    );
END
GO

IF OBJECT_ID(N'quik.boards', N'U') IS NOT NULL
BEGIN
    DROP INDEX IF EXISTS UQ_boards_code
    ON quik.boards;

    CREATE UNIQUE NONCLUSTERED  INDEX UQ_boards_code
    ON quik.boards (code);
END
GO
