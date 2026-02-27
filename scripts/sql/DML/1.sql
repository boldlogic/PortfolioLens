SELECT --TOP (1)
    q.instrument_class,
    q.ticker,
    q.isin,
    q.registration_number,
    q.full_name,
    q.short_name,
    q.face_value,
    q.maturity_date,
    q.coupon_duration,

    b.trade_point_id,
    b.board_id,
    q.class_code,
    q.class_name,

    it.type_id       AS instrument_type_id,
    q.instrument_type,
    ist.subtype_id   AS instrument_subtype_id,
    q.instrument_subtype,

    curr.iso_code    AS currency_id,
    q.currency,
    base_curr.iso_code AS base_currency_id,
    q.base_currency,
    quote_curr.iso_code AS quote_currency_id,
    q.quote_currency,
    counter_curr.iso_code AS counter_currency_id,
    q.counter_currency
FROM
    quik.current_quotes q
    LEFT JOIN quik.boards b ON b.code = q.class_code
    LEFT JOIN quik.trade_points p ON p.point_id = b.trade_point_id
    LEFT JOIN quik.instrument_types it ON it.title = q.instrument_type
    LEFT JOIN quik.instrument_subtypes ist ON ist.type_id = it.type_id AND ist.title = q.instrument_subtype
    LEFT JOIN dbo.currencies curr ON curr.iso_char_code = RTRIM(CASE WHEN q.currency IN ('SUR', 'RUR') THEN 'RUB' ELSE q.currency END)
    LEFT JOIN dbo.currencies base_curr ON base_curr.iso_char_code = RTRIM(CASE WHEN q.base_currency IN ('SUR', 'RUR') THEN 'RUB' ELSE q.base_currency END)
    LEFT JOIN dbo.currencies quote_curr ON quote_curr.iso_char_code = RTRIM(CASE WHEN q.quote_currency IN ('SUR', 'RUR') THEN 'RUB' ELSE q.quote_currency END)
    LEFT JOIN dbo.currencies counter_curr ON counter_curr.iso_char_code = RTRIM(CASE WHEN q.counter_currency IN ('SUR', 'RUR') THEN 'RUB' ELSE q.counter_currency END)
WHERE
    q.instrument_id IS NULL


