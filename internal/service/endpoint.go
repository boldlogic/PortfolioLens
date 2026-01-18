package service

type RequestType string

type Endpoint struct {
	Url            string
	Method         string
	Headers        map[string]string
	RequestTimeout int
	RetryPolicy    string
	RetryCount     int
}
