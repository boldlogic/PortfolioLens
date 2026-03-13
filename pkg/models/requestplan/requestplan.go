package requestplan

type RequestPlan struct {
	Url         string
	Method      string
	TimeoutMs   int
	RetryPolicy string
	RetryCount  int
	Params      []EndpointParam
}

type EndpointParam struct {
	Id            int
	EndpointId    int
	ParamId       *int
	ParamCode     string
	ExternalName  string
	ParamLocation string
	ExtCodeTypeId *uint8
	Format        string
	IsRequired    bool
	DefaultValue  string
}
