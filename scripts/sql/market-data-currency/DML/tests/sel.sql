DECLARE @quik_date DATE = '2026-03-10';
DECLARE @cbr_date  DATE = '2026-03-07';

;WITH quik_cc AS (
    SELECT
        ticker   = RTRIM(q.ticker),
        close_price,
        quote_date,
        base_cc  = RTRIM(UPPER(REPLACE(REPLACE(REPLACE(q.ticker, 'RUR','RUB'), 'SUR','RUB'), ' ', '')))
    FROM quik.current_quotes q
    WHERE q.class_code = 'CROSSRATE'
      AND CAST(q.quote_date AS DATE) = @quik_date
      AND LEN(RTRIM(q.ticker)) >= 2
)
SELECT
    q.ticker,
    q.quote_date   AS quik_date,
    f.date         AS cbr_date,
    q.base_cc      AS currency,
    f.rate_quote_per_base AS cbr_quote_per_unit,
    q.close_price         AS crossrate_close,
    (f.rate_quote_per_base - q.close_price) AS diff,
    CASE WHEN q.close_price <> 0
         THEN 100.0 * (f.rate_quote_per_base - q.close_price) / q.close_price
         ELSE NULL END AS pct_diff
FROM quik_cc q
JOIN dbo.currencies base_c ON base_c.iso_char_code = q.base_cc
JOIN dbo.fx_cbr_rates f
    ON f.date = @cbr_date
   AND f.base_iso_code  = base_c.iso_code
   AND f.quote_iso_code = 643   /* ЦБ всегда quote = RUB */
ORDER BY q.ticker;