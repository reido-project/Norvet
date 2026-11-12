package nvterr

import (
	"net/http"
)

// Domain "auth": authentication (401) and authorization (403).

// Authentication: "who are you?" (401).
var (
	ErrUnauthenticated = define("unauthenticated", http.StatusUnauthorized, "authentication required")
	ErrTokenMissing    = defineChild(ErrUnauthenticated, "token_missing", "authentication token is missing")
	ErrAPIKeyMissing   = defineChild(ErrUnauthenticated, "api_key_missing", "API key is missing")

	// Prefer plain ErrInvalidCredentials on login endpoints: telling the
	// client which part was wrong lets attackers enumerate accounts.
	ErrInvalidCredentials       = define("invalid_credentials", http.StatusUnauthorized, "invalid credentials")
	ErrCurrentPasswordIncorrect = defineChild(ErrInvalidCredentials, "current_password_incorrect", "current password is incorrect")
	ErrClientCredentialsInvalid = defineChild(ErrInvalidCredentials, "client_credentials_invalid", "client credentials are invalid")

	ErrTokenInvalid          = define("token_invalid", http.StatusUnauthorized, "token is invalid")
	ErrTokenMalformed        = defineChild(ErrTokenInvalid, "token_malformed", "token is malformed")
	ErrTokenSignatureInvalid = defineChild(ErrTokenInvalid, "token_signature_invalid", "token signature is invalid")
	ErrTokenIssuerInvalid    = defineChild(ErrTokenInvalid, "token_issuer_invalid", "token issuer is invalid")
	ErrTokenAudienceInvalid  = defineChild(ErrTokenInvalid, "token_audience_invalid", "token audience is invalid")
	ErrTokenNotYetValid      = defineChild(ErrTokenInvalid, "token_not_yet_valid", "token is not valid yet")
	ErrTokenTypeInvalid      = defineChild(ErrTokenInvalid, "token_type_invalid", "wrong token type for this request")
	ErrRefreshTokenInvalid   = defineChild(ErrTokenInvalid, "refresh_token_invalid", "refresh token is invalid")
	ErrRefreshTokenReused    = defineChild(ErrRefreshTokenInvalid, "refresh_token_reused", "refresh token was already used")

	ErrTokenExpired        = define("token_expired", http.StatusUnauthorized, "token has expired")
	ErrAccessTokenExpired  = defineChild(ErrTokenExpired, "access_token_expired", "access token has expired")
	ErrRefreshTokenExpired = defineChild(ErrTokenExpired, "refresh_token_expired", "refresh token has expired")

	ErrTokenRevoked        = define("token_revoked", http.StatusUnauthorized, "token has been revoked")
	ErrRefreshTokenRevoked = defineChild(ErrTokenRevoked, "refresh_token_revoked", "refresh token has been revoked")

	ErrSessionExpired = define("session_expired", http.StatusUnauthorized, "session has expired")

	ErrAPIKeyInvalid  = define("api_key_invalid", http.StatusUnauthorized, "API key is invalid")
	ErrAPIKeyExpired  = defineChild(ErrAPIKeyInvalid, "api_key_expired", "API key has expired")
	ErrAPIKeyRevoked  = defineChild(ErrAPIKeyInvalid, "api_key_revoked", "API key has been revoked")
	ErrAPIKeyDisabled = defineChild(ErrAPIKeyInvalid, "api_key_disabled", "API key is disabled")

	ErrSignatureInvalid        = define("signature_invalid", http.StatusUnauthorized, "signature is invalid")
	ErrWebhookSignatureInvalid = defineChild(ErrSignatureInvalid, "webhook_signature_invalid", "webhook signature is invalid")
	ErrRequestSignatureInvalid = defineChild(ErrSignatureInvalid, "request_signature_invalid", "request signature is invalid")
	ErrReplayDetected          = defineChild(ErrSignatureInvalid, "replay_detected", "request was already received or is too old")

	ErrMFAInvalid        = define("mfa_invalid", http.StatusUnauthorized, "invalid multi-factor authentication code")
	ErrTOTPCodeInvalid   = defineChild(ErrMFAInvalid, "totp_code_invalid", "authenticator code is invalid")
	ErrBackupCodeInvalid = defineChild(ErrMFAInvalid, "backup_code_invalid", "backup code is invalid or already used")
)

// Sessions and API keys that do not exist (404). They match ErrNotFound.
var (
	ErrSessionNotFound = defineChild(ErrNotFound, "session_not_found", "session not found")
	ErrAPIKeyNotFound  = defineChild(ErrNotFound, "api_key_not_found", "API key not found")
)

// One-time codes and links (400): email/phone verification, password reset.
var (
	ErrVerificationCodeInvalid      = define("verification_code_invalid", http.StatusBadRequest, "verification code is invalid")
	ErrEmailVerificationCodeInvalid = defineChild(ErrVerificationCodeInvalid, "email_verification_code_invalid", "email verification code is invalid")
	ErrPhoneVerificationCodeInvalid = defineChild(ErrVerificationCodeInvalid, "phone_verification_code_invalid", "phone verification code is invalid")
	ErrPasswordResetTokenInvalid    = defineChild(ErrVerificationCodeInvalid, "password_reset_token_invalid", "password reset link is invalid")
	ErrPasswordResetTokenUsed       = defineChild(ErrVerificationCodeInvalid, "password_reset_token_used", "password reset link was already used")

	ErrVerificationCodeExpired      = define("verification_code_expired", http.StatusBadRequest, "verification code has expired")
	ErrEmailVerificationCodeExpired = defineChild(ErrVerificationCodeExpired, "email_verification_code_expired", "email verification code has expired")
	ErrPhoneVerificationCodeExpired = defineChild(ErrVerificationCodeExpired, "phone_verification_code_expired", "phone verification code has expired")
	ErrPasswordResetTokenExpired    = defineChild(ErrVerificationCodeExpired, "password_reset_token_expired", "password reset link has expired")

	ErrOAuthStateInvalid = define("oauth_state_invalid", http.StatusBadRequest, "OAuth state is invalid")
	ErrOAuthCodeInvalid  = define("oauth_code_invalid", http.StatusBadRequest, "OAuth authorization code is invalid")
	ErrCaptchaInvalid    = define("captcha_invalid", http.StatusBadRequest, "captcha verification failed")
)

// Authorization: "you can't do that" (403). Everything here matches ErrForbidden.
var (
	ErrForbidden = define("forbidden", http.StatusForbidden, "you do not have access to this resource")

	ErrInsufficientPermissions = defineChild(ErrForbidden, "insufficient_permissions", "insufficient permissions")
	ErrInsufficientRole        = defineChild(ErrInsufficientPermissions, "insufficient_role", "your role does not allow this action")
	ErrInsufficientScope       = defineChild(ErrInsufficientPermissions, "insufficient_scope", "token does not have the required scope")

	ErrNotResourceOwner    = defineChild(ErrForbidden, "not_resource_owner", "you do not own this resource")
	ErrCrossTenantAccess   = defineChild(ErrForbidden, "cross_tenant_access", "access to another tenant's data is not allowed")
	ErrSelfActionForbidden = defineChild(ErrForbidden, "self_action_forbidden", "you cannot perform this action on your own account")
	ErrCannotModifyOwnRole = defineChild(ErrSelfActionForbidden, "cannot_modify_own_role", "you cannot change your own role")

	ErrIPNotAllowed     = defineChild(ErrForbidden, "ip_not_allowed", "access from this IP address is not allowed")
	ErrOriginNotAllowed = defineChild(ErrForbidden, "origin_not_allowed", "origin is not allowed")
	ErrCSRFInvalid      = defineChild(ErrForbidden, "csrf_invalid", "CSRF token is missing or invalid")
	ErrCaptchaRequired  = defineChild(ErrForbidden, "captcha_required", "captcha required")
	ErrMFARequired      = defineChild(ErrForbidden, "mfa_required", "multi-factor authentication required")

	ErrAccountNotVerified = defineChild(ErrForbidden, "account_not_verified", "account is not verified")
	ErrEmailNotVerified   = defineChild(ErrAccountNotVerified, "email_not_verified", "email address is not verified")
	ErrPhoneNotVerified   = defineChild(ErrAccountNotVerified, "phone_not_verified", "phone number is not verified")

	ErrAccountDisabled    = defineChild(ErrForbidden, "account_disabled", "account is disabled")
	ErrAccountSuspended   = defineChild(ErrAccountDisabled, "account_suspended", "account is suspended")
	ErrAccountDeactivated = defineChild(ErrAccountDisabled, "account_deactivated", "account is deactivated")
	ErrAccountBanned      = defineChild(ErrAccountDisabled, "account_banned", "account is banned")
	ErrAccountLocked      = defineChild(ErrForbidden, "account_locked", "account is locked")
)
