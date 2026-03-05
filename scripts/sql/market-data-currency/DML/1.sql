WITH c AS
(
	SELECT  distinct (q.currency)
	FROM quik.current_quotes q
	UNION
	SELECT  distinct (q.base_currency)
	FROM quik.current_quotes q
	UNION
	SELECT  distinct (q.counter_currency)
	FROM quik.current_quotes q
	UNION
	SELECT  distinct (q.quote_currency)
	FROM quik.current_quotes q
), cur as(
SELECT  distinct currency = CASE WHEN c.currency IN ('SUR','RUR','RUB') THEN 'RUB' else c.currency end
FROM c )
SELECT  cur.currency
       ,currency_name = coalesce(q.full_name,q.short_name)
FROM cur
LEFT JOIN quik.current_quotes q
ON q.ticker = cur.currency AND q.class_code = 'CROSSRATE'
WHERE cur.currency is not null
AND len (cur.currency) <= 3
AND not exists (
SELECT  1
FROM dbo.currencies c
WHERE c.iso_char_code = CASE WHEN cur.currency IN ('SUR', 'RUR', 'RUB') THEN 'RUB' else cur.currency end )