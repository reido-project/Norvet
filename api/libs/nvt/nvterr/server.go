package nvterr

import (
	"net/http"
)

// Domain "server": server-side and upstream failures (5xx). Messages
// are deliberately generic: attach the real error with WithCause so it is
// logged but never returned to clients.

var (
	// 500
	ErrInternal             = define("internal", http.StatusInternalServerError, "internal server error")
	ErrDatabaseFailed       = defineChild(ErrInternal, "database_failed", "internal server error")
	ErrEncryptionFailed     = defineChild(ErrInternal, "encryption_failed", "internal server error")
	ErrDecryptionFailed     = defineChild(ErrInternal, "decryption_failed", "internal server error")
	ErrSerializationFailed  = defineChild(ErrInternal, "serialization_failed", "internal server error")
	ErrConfigurationInvalid = defineChild(ErrInternal, "configuration_invalid", "internal server error")

	// 501
	ErrNotImplemented = define("not_implemented", http.StatusNotImplemented, "not implemented")

	// 502: an upstream answered badly.
	ErrDependencyFailed          = define("dependency_failed", http.StatusBadGateway, "upstream service failed")
	ErrDependencyResponseInvalid = defineChild(ErrDependencyFailed, "dependency_response_invalid", "upstream service returned an invalid response")
	ErrPaymentProviderFailed     = defineChild(ErrDependencyFailed, "payment_provider_failed", "payment provider failed")
	ErrEmailProviderFailed       = defineChild(ErrDependencyFailed, "email_provider_failed", "email provider failed")
	ErrSMSProviderFailed         = defineChild(ErrDependencyFailed, "sms_provider_failed", "SMS provider failed")
	ErrOAuthProviderFailed       = defineChild(ErrDependencyFailed, "oauth_provider_failed", "identity provider failed")

	// 503
	ErrServiceUnavailable    = define("service_unavailable", http.StatusServiceUnavailable, "service temporarily unavailable")
	ErrMaintenance           = defineChild(ErrServiceUnavailable, "maintenance", "service is under maintenance")
	ErrShuttingDown          = defineChild(ErrServiceUnavailable, "shutting_down", "service is shutting down")
	ErrOverloaded            = defineChild(ErrServiceUnavailable, "overloaded", "service is overloaded, try again later")
	ErrCircuitBreakerOpen    = defineChild(ErrServiceUnavailable, "circuit_breaker_open", "service temporarily unavailable")
	ErrDependencyUnavailable = defineChild(ErrServiceUnavailable, "dependency_unavailable", "a required service is unavailable")
	ErrDatabaseUnavailable   = defineChild(ErrDependencyUnavailable, "database_unavailable", "database is unavailable")
	ErrCacheUnavailable      = defineChild(ErrDependencyUnavailable, "cache_unavailable", "cache is unavailable")
	ErrQueueUnavailable      = defineChild(ErrDependencyUnavailable, "queue_unavailable", "message queue is unavailable")
	ErrSearchUnavailable     = defineChild(ErrDependencyUnavailable, "search_unavailable", "search is unavailable")
	ErrStorageUnavailable    = defineChild(ErrDependencyUnavailable, "storage_unavailable", "file storage is unavailable")

	// 504
	ErrGatewayTimeout         = define("gateway_timeout", http.StatusGatewayTimeout, "upstream service timed out")
	ErrDatabaseTimeout        = defineChild(ErrGatewayTimeout, "database_timeout", "database timed out")
	ErrExternalServiceTimeout = defineChild(ErrGatewayTimeout, "external_service_timeout", "external service timed out")
)
