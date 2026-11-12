package nvterr

// Domain "account": user identity, sign-up and profile errors.

// Unique identity fields. Each matches ErrAlreadyRegistered, ErrAlreadyExists
// and ErrConflict too.
var (
	ErrAlreadyRegistered         = defineChild(ErrAlreadyExists, "already_registered", "already registered")
	ErrUsernameAlreadyRegistered = defineChild(ErrAlreadyRegistered, "username_already_registered", "username is already registered")
	ErrEmailAlreadyRegistered    = defineChild(ErrAlreadyRegistered, "email_already_registered", "email is already registered")
	ErrPhoneAlreadyRegistered    = defineChild(ErrAlreadyRegistered, "phone_already_registered", "phone number is already registered")

	ErrSocialAccountAlreadyLinked = defineChild(ErrAlreadyExists, "social_account_already_linked", "this social account is already linked to another user")
)

var (
	ErrUserNotFound = defineChild(ErrNotFound, "user_not_found", "user not found")

	ErrAccountDeleted = defineChild(ErrGone, "account_deleted", "account has been deleted")

	ErrEmailAlreadyVerified = defineChild(ErrAlreadyVerified, "email_already_verified", "email is already verified")
	ErrPhoneAlreadyVerified = defineChild(ErrAlreadyVerified, "phone_already_verified", "phone number is already verified")

	ErrProfileIncomplete = defineChild(ErrBusinessRuleViolation, "profile_incomplete", "profile must be completed first")
)
