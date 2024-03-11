package constants

type RequestKey string

const (
	NewRelicTransaction RequestKey = "newRelicTransaction"
	RequestId           RequestKey = "x-request-id"
)
