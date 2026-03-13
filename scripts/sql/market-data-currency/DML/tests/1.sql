WITH
	c AS (
		SELECT distinct
			(q.currency)
		FROM
			quik.current_quotes q
		UNION
		SELECT distinct
			(q.base_currency)
		FROM
			quik.current_quotes q
		UNION
		SELECT distinct
			(q.counter_currency)
		FROM
			quik.current_quotes q
		UNION
		SELECT distinct
			(q.quote_currency)
		FROM
			quik.current_quotes q
	),
	cur as (
		SELECT distinct
			currency=CASE
				WHEN c.currency IN ('SUR', 'RUR', 'RUB') THEN 'RUB'
				else c.currency
			end
		FROM
			c
	)
SELECT
	cur.currency,
	currency_name=coalesce(q.full_name, q.short_name)
FROM
	cur
	LEFT JOIN quik.current_quotes q ON q.ticker=cur.currency
	AND q.class_code='CROSSRATE'
WHERE
	cur.currency is not null
	AND len (cur.currency)<=3
	AND not exists (
		SELECT
			1
		FROM
			dbo.currencies c
		WHERE
			c.iso_char_code=CASE
				WHEN cur.currency IN ('SUR', 'RUR', 'RUB') THEN 'RUB'
				else cur.currency
			end
	)
	---- 

MERGE INTO dbo.currencies AS tgt USING src ON tgt.iso_code=src.iso_code
	AND tgt.iso_char_code=src.iso_char_code 
WHEN MATCHED
	AND (
		tgt.currency_name<>src.currency_name
		OR tgt.lat_name<>src.lat_name
		OR tgt.minor_units<>src.minor_units
	) THEN
		UPDATE
			SET
				tgt.currency_name=src.currency_name,
				tgt.lat_name=src.lat_name,
				tgt.minor_units<>src.minor_units,
				tgt.updated_at=getdate () 
WHEN NOT MATCHED BY TARGET 
	THEN 
		INSERT (
				iso_code,
				iso_char_code,
				currency_name,
				lat_name,
				minor_units,
				updated_at,
				ext_system_id
				)
		VALUES
			(
				src.iso_code,
				src.iso_char_code,
				src.currency_name,
				src.lat_name,
				src.minor_units,
				src.updated_at,
				src.ext_system_id
			);