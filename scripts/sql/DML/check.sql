-- Почему котировки без instrument_id не попадают в выборку
SELECT
    CASE
        WHEN b.board_id IS NULL THEN 'Нет борда (class_code не в boards)'
        WHEN p.point_id IS NULL THEN 'Нет площадки у борда'
        WHEN it.type_id IS NULL THEN 'Нет типа инструмента (instrument_type не в instrument_types)'
        ELSE 'Можно привязать (ещё не обработано)'
    END AS reason,
    COUNT(*) AS cnt
FROM quik.current_quotes q
LEFT JOIN quik.boards b ON b.code = q.class_code
LEFT JOIN quik.trade_points p ON p.point_id = b.trade_point_id
LEFT JOIN quik.instrument_types it ON it.title = q.instrument_type
WHERE q.instrument_id IS NULL
GROUP BY
    CASE
        WHEN b.board_id IS NULL THEN 'Нет борда (class_code не в boards)'
        WHEN p.point_id IS NULL THEN 'Нет площадки у борда'
        WHEN it.type_id IS NULL THEN 'Нет типа инструмента (instrument_type не в instrument_types)'
        ELSE 'Можно привязать (ещё не обработано)'
    END
ORDER BY cnt DESC;