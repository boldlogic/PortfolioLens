IF NOT EXISTS (
    SELECT
        *
    FROM
        sys.tables t
        INNER JOIN sys.schemas s ON t.schema_id=s.schema_id
    WHERE
        s.name=N'quik'
        AND t.name=N'instrument_board'
) BEGIN
CREATE TABLE
    quik.instrument_board (
        instrument_id BIGINT,
        board_id smallint,
        base_currency bigint NULL,
        is_primary bit NULL,
        CONSTRAINT PK_quik_instrument_board PRIMARY KEY CLUSTERED (instrument_id, board_id),
        CONSTRAINT FK_quik_instrument_board_instrument FOREIGN KEY (instrument_id) REFERENCES quik.instruments (instrument_id),
        CONSTRAINT FK_quik_instrument_board_board FOREIGN KEY (board_id) REFERENCES quik.boards (board_id),
        CONSTRAINT FK_quik_instrument_board_base_currency FOREIGN KEY (base_currency) REFERENCES dbo.currencies (iso_code),
    );

END GO