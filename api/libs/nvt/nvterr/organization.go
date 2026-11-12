package nvterr

// Domain "organization": organizations, members, roles, invites, tenants.

var (
	ErrOrganizationNotFound = defineChild(ErrNotFound, "organization_not_found", "organization not found")
	ErrMemberNotFound       = defineChild(ErrNotFound, "member_not_found", "member not found")
	ErrRoleNotFound         = defineChild(ErrNotFound, "role_not_found", "role not found")
	ErrInviteNotFound       = defineChild(ErrNotFound, "invite_not_found", "invitation not found")
	ErrTenantNotFound       = defineChild(ErrNotFound, "tenant_not_found", "tenant not found")

	ErrOrganizationAlreadyExists = defineChild(ErrAlreadyExists, "organization_already_exists", "organization already exists")
	ErrMemberAlreadyExists       = defineChild(ErrAlreadyExists, "member_already_exists", "user is already a member")
	ErrRoleAlreadyExists         = defineChild(ErrAlreadyExists, "role_already_exists", "role already exists")
	ErrInviteAlreadyExists       = defineChild(ErrAlreadyExists, "invite_already_exists", "an invitation is already pending for this user")

	ErrInviteExpired         = defineChild(ErrExpired, "invite_expired", "invitation has expired")
	ErrInviteRevoked         = defineChild(ErrGone, "invite_revoked", "invitation has been revoked")
	ErrInviteAlreadyAccepted = defineChild(ErrAlreadyAccepted, "invite_already_accepted", "invitation was already accepted")

	ErrRoleInUse = defineChild(ErrInUse, "role_in_use", "role is assigned to members and cannot be deleted")

	ErrOrganizationSuspended = defineChild(ErrForbidden, "organization_suspended", "organization is suspended")
	ErrTenantDisabled        = defineChild(ErrForbidden, "tenant_disabled", "tenant is disabled")
	ErrTenantMismatch        = defineChild(ErrForbidden, "tenant_mismatch", "resource belongs to a different tenant")

	ErrTenantContextMissing = defineChild(ErrMissingHeader, "tenant_context_missing", "tenant could not be determined for this request")

	ErrMemberLimitReached = defineChild(ErrLimitExceeded, "member_limit_reached", "member limit reached")

	ErrCannotDeleteLastAdmin = defineChild(ErrBusinessRuleViolation, "cannot_delete_last_admin", "cannot remove the last administrator")
	ErrCannotRemoveLastOwner = defineChild(ErrBusinessRuleViolation, "cannot_remove_last_owner", "cannot remove the last owner")
)
