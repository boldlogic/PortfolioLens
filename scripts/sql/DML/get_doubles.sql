
with doubles as (
select ticker from quik.current_quotes
group by ticker
having count (1)>1)
SELECT  instrument_class
      ,q.ticker
	  ,p.*
	  ,[trading_status]
      ,registration_number
      ,full_name
      ,short_name
      ,class_code
      ,class_name
      ,instrument_type
      ,instrument_subtype
      ,isin
      ,face_value
      ,base_currency
      ,quote_currency
      ,counter_currency
      ,maturity_date
      ,coupon_duration
      ,rw
      ,instrument_id
	  ,b.code
	  ,b.trade_point_id
	  
  FROM quik.current_quotes q
  join doubles on q.ticker=doubles.ticker
  left join quik.boards b on b.code=q.class_code
  left join quik.trade_points p on p.point_id=b.trade_point_id
  order by ticker

