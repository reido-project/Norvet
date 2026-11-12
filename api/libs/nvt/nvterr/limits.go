package nvterr

import (
	"net/http"
)

// Domain "limits": throttling, quotas, plan limits and feature gating.

// 429
var (
	ErrRateLimited              = define("rate_limited", http.StatusTooManyRequests, "too many requests")
	ErrConcurrencyLimitExceeded = defineChild(ErrRateLimited, "concurrency_limit_exceeded", "too many concurrent requests")
	ErrVerificationRateLimited  = defineChild(ErrRateLimited, "verification_rate_limited", "please wait before requesting another verification")

	ErrTooManyAttempts               = define("too_many_attempts", http.StatusTooManyRequests, "too many attempts, try again later")
	ErrLoginAttemptsExceeded         = defineChild(ErrTooManyAttempts, "login_attempts_exceeded", "too many failed login attempts, try again later")
	ErrOTPAttemptsExceeded           = defineChild(ErrTooManyAttempts, "otp_attempts_exceeded", "too many incorrect codes, try again later")
	ErrPasswordResetAttemptsExceeded = defineChild(ErrTooManyAttempts, "password_reset_attempts_exceeded", "too many password reset requests, try again later")
)

// 403
var (
	ErrQuotaExceeded        = define("quota_exceeded", http.StatusForbidden, "quota exceeded")
	ErrDailyQuotaExceeded   = defineChild(ErrQuotaExceeded, "daily_quota_exceeded", "daily quota exceeded")
	ErrMonthlyQuotaExceeded = defineChild(ErrQuotaExceeded, "monthly_quota_exceeded", "monthly quota exceeded")

	ErrLimitExceeded       = define("limit_exceeded", http.StatusForbidden, "limit exceeded")
	ErrSessionLimitReached = defineChild(ErrLimitExceeded, "session_limit_reached", "maximum number of active sessions reached")
	ErrAPIKeyLimitReached  = defineChild(ErrLimitExceeded, "api_key_limit_reached", "maximum number of API keys reached")

	ErrFeatureUnavailable  = define("feature_unavailable", http.StatusForbidden, "feature is not available")
	ErrFeatureDisabled     = defineChild(ErrFeatureUnavailable, "feature_disabled", "feature is disabled")
	ErrPlanUpgradeRequired = defineChild(ErrFeatureUnavailable, "plan_upgrade_required", "your plan does not include this feature")
)
