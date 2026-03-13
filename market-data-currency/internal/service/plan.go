package service

import (
	"context"
	"fmt"

	"github.com/boldlogic/PortfolioLens/pkg/models/requestplan"
	"github.com/boldlogic/PortfolioLens/pkg/transport/requestplanner"
)

func (s *Service) fillPlanForTask(ctx context.Context, rawPlan requestplan.RequestPlan, taskParams map[string]string, actionId uint8) (requestplanner.RequestPlan, error) {
	resolved := requestplanner.RequestPlan{
		Url:    rawPlan.Url,
		Method: rawPlan.Method,
	}

	for _, ep := range rawPlan.Params {
		value, skip, err := s.resolveParamValue(ctx, ep, taskParams, actionId)
		if err != nil {
			return requestplanner.RequestPlan{}, err
		}
		if skip {
			continue
		}

		resolved.Params = append(resolved.Params, requestplanner.EndpointParam{
			ExternalName: ep.ExternalName,
			Location:     requestplanner.ParamLocation(ep.ParamLocation),
			Format:       ep.Format,
			IsRequired:   ep.IsRequired,
			Value:        value,
		})
	}

	return resolved, nil
}

func (s *Service) resolveParamValue(
	ctx context.Context,
	ep requestplan.EndpointParam,
	taskParams map[string]string,
	actionId uint8,
) (value string, skip bool, err error) {
	switch {

	case ep.ParamId == nil:
		return ep.DefaultValue, false, nil

	case ep.ExtCodeTypeId != nil:
		rawValue, ok := taskParams[ep.ParamCode]
		if !ok {
			if ep.IsRequired {
				return "", false, fmt.Errorf("отсутствует обязательный параметр '%s'", ep.ParamCode)
			}
			return "", true, nil
		}
		extCode, err := s.schedulerRepo.SelectExternalCodeByCurrency(ctx, rawValue, *ep.ExtCodeTypeId, actionId)
		if err != nil {
			return "", false, fmt.Errorf("внешний код для '%s'='%s': %w", ep.ParamCode, rawValue, err)
		}
		return extCode, false, nil

	default:
		rawValue, ok := taskParams[ep.ParamCode]
		if !ok {
			if ep.IsRequired {
				return "", false, fmt.Errorf("отсутствует обязательный параметр '%s'", ep.ParamCode)
			}
			return "", true, nil
		}
		return rawValue, false, nil
	}
}
