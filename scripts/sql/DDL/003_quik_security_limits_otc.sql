IF NOT EXISTS (
    SELECT
        *
    FROM
        sys.tables t
        INNER JOIN sys.schemas s ON t.schema_id=s.schema_id
    WHERE
        s.name=N'quik'
        AND t.name=N'security_limits_otc'
) BEGIN
CREATE TABLE
    quik.security_limits_otc (
        load_date date NOT NULL DEFAULT (getdate ()),
        client_code varchar(12) NOT NULL,
        ticker varchar(15) NOT NULL,
        trade_account varchar(15) NOT NULL DEFAULT 'OTC',
        settle_code varchar(2) NOT NULL DEFAULT 'Tx',
        firm_code varchar(12) NOT NULL,
        firm_name varchar(128) NULL,
        balance float NULL,
        acquisition_ccy varchar(3) NULL,
        isin varchar(15) NULL,
        ts timestamp NOT NULL,
        
        CONSTRAINT PK_quik_security_limits_otc PRIMARY KEY CLUSTERED (
            load_date ASC,
            client_code ASC,
            ticker ASC,
            trade_account ASC,
            settle_code ASC,
            firm_code ASC
        ),
    );

END 
GO