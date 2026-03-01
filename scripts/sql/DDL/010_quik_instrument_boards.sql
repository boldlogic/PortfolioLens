IF NOT EXISTS (
  SELECT
    *
  FROM
    sys.tables t
    INNER JOIN sys.schemas s ON t.schema_id = s.schema_id
  WHERE
    s.name = N'quik'
    AND t.name = N'instrument_boards'
) 

BEGIN
CREATE TABLE
  quik.instrument_boards (
    instrument_id BIGINT NOT NULL,
    board_id tinyint NOT NULL,

    type_id       tinyint NOT NULL,
    subtype_id       tinyint NULL,

    currency_id BIGINT NULL,
    base_currency_id BIGINT NULL,
    quote_currency_id BIGINT NULL,
    counter_currency_id BIGINT NULL,
    is_primary bit NULL,
    CONSTRAINT PK_quik_instrument_boards PRIMARY KEY CLUSTERED (instrument_id, board_id),
    CONSTRAINT FK_quik_instrument_boards_instrument FOREIGN KEY (instrument_id) REFERENCES quik.instruments (instrument_id),
    CONSTRAINT FK_quik_instrument_boards_board FOREIGN KEY (board_id) REFERENCES quik.boards (board_id),

    CONSTRAINT FK_quik_instrument_boards_type FOREIGN KEY (type_id) REFERENCES quik.instrument_types (type_id),
    CONSTRAINT FK_quik_instrument_boards_subtype FOREIGN KEY (subtype_id) REFERENCES quik.instrument_subtypes (subtype_id),
    CONSTRAINT FK_quik_instrument_boards_currency FOREIGN KEY (currency_id) REFERENCES dbo.currencies (iso_code),
    CONSTRAINT FK_quik_instrument_boards_base_currency FOREIGN KEY (base_currency_id) REFERENCES dbo.currencies (iso_code),
    CONSTRAINT FK_quik_instrument_boards_quote_currency FOREIGN KEY (quote_currency_id) REFERENCES dbo.currencies (iso_code),
    CONSTRAINT FK_quik_instrument_boards_counter_currency FOREIGN KEY (counter_currency_id) REFERENCES dbo.currencies (iso_code),
  );


END 
GO