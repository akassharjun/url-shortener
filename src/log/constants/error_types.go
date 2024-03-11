package constants

type ErrorTypes string

const (
	ApplicationInit ErrorTypes = "APPLICATION_INIT"
	Application     ErrorTypes = "APPLICATION"
	InternalServer  ErrorTypes = "INTERNAL_SERVER"
	ExternalService ErrorTypes = "EXTERNAL_SERVICE"
)
