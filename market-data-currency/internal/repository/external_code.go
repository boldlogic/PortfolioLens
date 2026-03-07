package repository

import (
	"context"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
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
			if shutdown.IsExceeded(err) {
				return err
			}
			r.Logger.Error("ошибка сохранения external_code",
				zap.String("ext_code", c.Code),
				zap.Uint8("ext_system_id", uint8(c.ExternalSystemId)),
				zap.Error(err))
			return apperrors.ErrSavingData
		}
	}
	return nil
}
