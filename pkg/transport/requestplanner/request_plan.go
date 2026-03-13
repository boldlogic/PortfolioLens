package requestplanner

import "fmt"

type RequestPlan struct {
	Url    string
	Method string
	Params []EndpointParam
}

type EndpointParam struct {
	ExternalName string
	Location     ParamLocation
	Format       string
	IsRequired   bool
	Value        string
}

type ParamLocation string

const (
	ParamLocationQuery  ParamLocation = "query"
	ParamLocationHeader ParamLocation = "header"
)

func (s ParamLocation) Valid() bool {
	switch s {
	case ParamLocationQuery, ParamLocationHeader:
		return true
	default:
		return false
	}
}

func fillParams(params []EndpointParam) (map[string]string, map[string]string, error) {
	if len(params) == 0 {
		return nil, nil, nil
	}

	query := make(map[string]string)
	headers := make(map[string]string)

	for _, param := range params {
		if err := checkParam(param); err != nil {
			return nil, nil, err
		}

		value := param.Value
		if value == "" {
			continue
		}

		var err error
		value, err = formatParamValue(value, param.Format)
		if err != nil {
			return nil, nil, err
		}

		switch param.Location {
		case ParamLocationQuery:
			query[param.ExternalName] = value
		case ParamLocationHeader:
			headers[param.ExternalName] = value
		}
	}
	return query, headers, nil
}

func checkParam(p EndpointParam) error {
	if p.ExternalName == "" {
		return fmt.Errorf("не заполнено имя параметра")
	}
	if !p.Location.Valid() {
		return fmt.Errorf("неизвестное расположение параметра")
	}
	if p.IsRequired && p.Value == "" {
		return fmt.Errorf("не заполнен обязательный параметр '%s'", p.ExternalName)
	}
	return nil
}
