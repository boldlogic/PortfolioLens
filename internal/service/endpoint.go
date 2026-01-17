package service

type RequestType string

type Endpoint struct {
	url            string
	method         string
	headers        map[string]string
	requestTimeout int
	retryPolicy    string
	retryCount     int
}
