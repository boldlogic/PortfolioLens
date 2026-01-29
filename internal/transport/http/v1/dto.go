package v1

import "time"

type APIResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}

type Body struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type currencyDTO struct {
	ISOCode  int    `json:"isoCode"`
	CharCode string `json:"code"`
	NameRu   string `json:"nameRu,omitempty"`
	NameEn   string `json:"nameEn,omitempty"`
}

type newTaskDTO struct {
	Action string `json:"action"`
	Uuid   string `json:"uuid,omitempty"`

	Params params `json:"params,omitempty"`
	//ScheduledAt string `json:"scheduledAt,omitempty"`
}
type newTaskRespDTO struct {
	Id          int       `json:"id"`
	Uuid        string    `json:"uuid"`
	CreatedAt   time.Time `json:"createdAt"`
	ScheduledAt time.Time `json:"scheduledAt"`
}

type params struct {
	CcyCode  string `json:"ccyCode,omitempty"`
	DateFrom string `json:"dateFrom,omitempty"`
	DateTo   string `json:"dateTo,omitempty"`
}
