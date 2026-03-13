package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"go.uber.org/zap"
)

const mergeExternalCodes = `
	WITH src AS (
		SELECT
			ext_system_id   = @p1,
			ext_code        = @p2,
			ext_code_type_id = @p3,
			internal_id     = @p4
	)
	MERGE INTO dbo.external_codes AS tgt
	USING src
	ON tgt.ext_system_id = src.ext_system_id
		AND tgt.ext_code = src.ext_code
		AND tgt.ext_code_type_id = src.ext_code_type_id
	WHEN MATCHED AND tgt.internal_id <> src.internal_id
	THEN UPDATE SET tgt.internal_id = src.internal_id
	WHEN NOT MATCHED BY TARGET
	THEN INSERT (ext_system_id, ext_code, ext_code_type_id, internal_id)
		VALUES (src.ext_system_id, src.ext_code, src.ext_code_type_id, src.internal_id);`

func (r *Repository) MergeExternalCodes(ctx context.Context, codes []models.ExternalCode) error {
	if len(codes) == 0 {
		return nil
	}

	for _, c := range codes {
		_, err := r.Db.ExecContext(ctx, mergeExternalCodes,
			uint8(c.ExternalSystemId),
			c.Code,
			uint8(c.Type),
			c.IntId,
		)
		if err != nil {
			if r.isShutdown(err) {
				return err
			}
			r.Logger.Error("ошибка сохранения external_code",
				zap.String("ext_code", c.Code),
				zap.Uint8("ext_system_id", uint8(c.ExternalSystemId)),
				zap.Error(err))
			return models.ErrSavingData
		}
	}
	return nil
}

const selectExternalCodeByCurrency = `
	SELECT ec.ext_code
	FROM dbo.currencies c
	JOIN dbo.external_codes ec
		ON ec.internal_id     = c.iso_code
		AND ec.ext_code_type_id = @p2
		AND ec.ext_system_id   = (
			SELECT u.external_system_id
			FROM dbo.endpoints e
			JOIN dbo.external_system_urls u ON u.id = e.external_system_url_id
			JOIN dbo.actions a              ON a.endpoint_id = e.id
			WHERE a.id = @p3
		)
	WHERE c.iso_char_code = @p1`

func (r *Repository) SelectExternalCodeByCurrency(ctx context.Context, isoCharCode string, extCodeTypeId uint8, actionId uint8) (string, error) {
	var code string
	row := r.Db.QueryRowContext(ctx, selectExternalCodeByCurrency, isoCharCode, extCodeTypeId, actionId)
	err := row.Scan(&code)
	if err != nil {
		if r.isShutdown(err) {
			return "", err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Debug("внешний код не найден",
				zap.String("iso_char_code", isoCharCode),
				zap.Uint8("ext_code_type_id", extCodeTypeId))
			return "", models.ErrNotFound
		}
		r.Logger.Error("ошибка при получении внешнего кода",
			zap.String("iso_char_code", isoCharCode),
			zap.Uint8("ext_code_type_id", extCodeTypeId),
			zap.Error(err))
		return "", models.ErrRetrievingData
	}
	return code, nil
}
