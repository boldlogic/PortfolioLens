package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/requestplan"
	"go.uber.org/zap"
)

const selectEndpointPlan = `
	SELECT
		u.proto,
		u.host,
		u.port,
		e.path,
		e.method,
		e.timeout_ms,
		e.retry_policy,
		e.retry_count,
		ep.id              AS ep_id,
		ep.param_id,
		p.code             AS param_code,
		ep.external_name,
		ep.param_location,
		ep.ext_code_type_id,
		ep.format,
		ep.is_required,
		ep.default_value
	FROM dbo.actions a
	JOIN dbo.endpoints e         ON e.id = a.endpoint_id
	JOIN dbo.external_system_urls u ON u.id = e.external_system_url_id
	LEFT JOIN dbo.endpoint_params ep ON ep.endpoint_id = e.id
	LEFT JOIN dbo.params p        ON p.id = ep.param_id
	WHERE a.id = @p1
		AND e.is_active = 1
		AND u.is_active = 1`

type rawEndpointRow struct {
	Proto         string
	Host          string
	Port          sql.NullInt32
	Path          string
	Method        string
	TimeoutMs     int
	RetryPolicy   string
	RetryCount    int
	EpId          sql.NullInt32
	ParamId       sql.NullInt32
	ParamCode     sql.NullString
	ExternalName  sql.NullString
	ParamLocation sql.NullString
	ExtCodeTypeId sql.NullByte
	Format        sql.NullString
	IsRequired    sql.NullBool
	DefaultValue  sql.NullString
}

func (r *Repository) SelectRequestPlan(ctx context.Context, actionId uint8) (requestplan.RequestPlan, error) {
	rows, err := r.Db.QueryContext(ctx, selectEndpointPlan, actionId)
	if err != nil {
		if r.isShutdown(err) {
			return requestplan.RequestPlan{}, err
		}
		r.Logger.Error("ошибка при загрузке плана запроса", zap.Uint8("action_id", actionId), zap.Error(err))
		return requestplan.RequestPlan{}, models.ErrRetrievingData
	}
	defer rows.Close()

	var plan requestplan.RequestPlan
	var params []requestplan.EndpointParam
	initialized := false

	for rows.Next() {
		var row rawEndpointRow
		err = rows.Scan(
			&row.Proto,
			&row.Host,
			&row.Port,
			&row.Path,
			&row.Method,
			&row.TimeoutMs,
			&row.RetryPolicy,
			&row.RetryCount,
			&row.EpId,
			&row.ParamId,
			&row.ParamCode,
			&row.ExternalName,
			&row.ParamLocation,
			&row.ExtCodeTypeId,
			&row.Format,
			&row.IsRequired,
			&row.DefaultValue,
		)
		if err != nil {
			if r.isShutdown(err) {
				return requestplan.RequestPlan{}, err
			}
			r.Logger.Error("ошибка при чтении строки плана запроса", zap.Error(err))
			return requestplan.RequestPlan{}, models.ErrRetrievingData
		}

		if !initialized {
			host := row.Host
			if row.Port.Valid {
				host = fmt.Sprintf("%s:%d", host, row.Port.Int32)
			}
			plan = requestplan.RequestPlan{
				Url:         row.Proto + "://" + host + "/" + row.Path,
				Method:      row.Method,
				TimeoutMs:   row.TimeoutMs,
				RetryPolicy: row.RetryPolicy,
				RetryCount:  row.RetryCount,
			}
			initialized = true
		}

		if !row.EpId.Valid {
			continue
		}

		ep := requestplan.EndpointParam{
			Id:            int(row.EpId.Int32),
			ExternalName:  row.ExternalName.String,
			ParamLocation: row.ParamLocation.String,
			IsRequired:    row.IsRequired.Bool,
		}
		if row.ParamId.Valid {
			id := int(row.ParamId.Int32)
			ep.ParamId = &id
			ep.ParamCode = row.ParamCode.String
		}
		if row.ExtCodeTypeId.Valid {
			v := row.ExtCodeTypeId.Byte
			ep.ExtCodeTypeId = &v
		}
		if row.Format.Valid {
			ep.Format = row.Format.String
		}
		if row.DefaultValue.Valid {
			ep.DefaultValue = row.DefaultValue.String
		}
		params = append(params, ep)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка при итерации строк плана запроса", zap.Error(rows.Err()))
		return requestplan.RequestPlan{}, models.ErrRetrievingData
	}
	if !initialized {
		r.Logger.Debug("план запроса не найден", zap.Uint8("action_id", actionId))
		return requestplan.RequestPlan{}, models.ErrNotFound
	}

	plan.Params = params
	return plan, nil
}
