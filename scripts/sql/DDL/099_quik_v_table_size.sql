CREATE OR ALTER VIEW quik.v_table_size
AS
SELECT
    s.name AS schema_name,
    t.name AS table_name,

    SUM(
        CASE 
            WHEN i.index_id IN (0, 1)
            THEN p.rows
            ELSE 0
        END
    ) AS row_count,

    SUM(
        CASE 
            WHEN i.index_id IN (0, 1)
            THEN au.total_pages
            ELSE 0
        END
    ) * 8 / 1024.0 AS data_mb,

    SUM(
        CASE 
            WHEN i.index_id = 1
            THEN au.total_pages
            ELSE 0
        END
    ) * 8 / 1024.0 AS clustered_index_mb,

    SUM(
        CASE 
            WHEN i.index_id > 1
            THEN au.total_pages
            ELSE 0
        END
    ) * 8 / 1024.0 AS nonclustered_indexes_mb,

    SUM(au.total_pages) * 8 / 1024.0 AS total_mb

FROM sys.tables t
JOIN sys.schemas s
    ON t.schema_id = s.schema_id
JOIN sys.indexes i
    ON t.object_id = i.object_id
JOIN sys.partitions p
    ON i.object_id = p.object_id
   AND i.index_id  = p.index_id
JOIN sys.allocation_units au
    ON p.partition_id = au.container_id

WHERE s.name = 'quik'

GROUP BY
    s.name,
    t.name;
GO
