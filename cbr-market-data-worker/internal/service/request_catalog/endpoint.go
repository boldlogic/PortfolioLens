package request_catalog

type RequestType string

type RequestPlan struct {
	Url            string
	Method         string
	Params         []Param
	Headers        map[string]string
	RequestTimeout int
	RetryPolicy    string
	RetryCount     int
}

type Param struct {
	Name string
	Type string
}
