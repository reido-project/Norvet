# nvterr error reference

267 errors. Indentation shows inheritance: `errors.Is(child, parent)` is true, never the reverse.

| Code | Status | Message |
|---|---|---|
| `api_version_unsupported` | 400 | API version is not supported |
| `bad_request` | 400 | bad request |
| `captcha_invalid` | 400 | captcha verification failed |
| `empty_body` | 400 | request body is required |
| &nbsp;&nbsp;&nbsp;&nbsp;`file_empty` | 400 | uploaded file is empty |
| `invalid_header` | 400 | header is invalid |
| `invalid_parameter` | 400 | parameter is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_filter` | 400 | invalid filter parameter |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_id` | 400 | identifier is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_uuid` | 400 | identifier is not a valid UUID |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_pagination` | 400 | invalid pagination parameters |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_cursor` | 400 | invalid or expired pagination cursor |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_page` | 400 | invalid page number |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_page_size` | 400 | invalid page size |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_search_query` | 400 | invalid search query |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_sort` | 400 | invalid sort parameter |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_sort_field` | 400 | cannot sort by this field |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_sort_order` | 400 | sort order must be asc or desc |
| `malformed_body` | 400 | request body is malformed |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_form` | 400 | request form data is malformed |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_json` | 400 | request body is not valid JSON |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_xml` | 400 | request body is not valid XML |
| `missing_header` | 400 | required header is missing |
| &nbsp;&nbsp;&nbsp;&nbsp;`idempotency_key_missing` | 400 | Idempotency-Key header is required |
| &nbsp;&nbsp;&nbsp;&nbsp;`missing_content_type` | 400 | Content-Type header is required |
| &nbsp;&nbsp;&nbsp;&nbsp;`tenant_context_missing` | 400 | tenant could not be determined for this request |
| `missing_parameter` | 400 | required parameter is missing |
| `oauth_code_invalid` | 400 | OAuth authorization code is invalid |
| `oauth_state_invalid` | 400 | OAuth state is invalid |
| `unknown_field` | 400 | request contains an unknown field |
| `verification_code_expired` | 400 | verification code has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`email_verification_code_expired` | 400 | email verification code has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`password_reset_token_expired` | 400 | password reset link has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`phone_verification_code_expired` | 400 | phone verification code has expired |
| `verification_code_invalid` | 400 | verification code is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`email_verification_code_invalid` | 400 | email verification code is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`password_reset_token_invalid` | 400 | password reset link is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`password_reset_token_used` | 400 | password reset link was already used |
| &nbsp;&nbsp;&nbsp;&nbsp;`phone_verification_code_invalid` | 400 | phone verification code is invalid |
| `api_key_invalid` | 401 | API key is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`api_key_disabled` | 401 | API key is disabled |
| &nbsp;&nbsp;&nbsp;&nbsp;`api_key_expired` | 401 | API key has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`api_key_revoked` | 401 | API key has been revoked |
| `invalid_credentials` | 401 | invalid credentials |
| &nbsp;&nbsp;&nbsp;&nbsp;`client_credentials_invalid` | 401 | client credentials are invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`current_password_incorrect` | 401 | current password is incorrect |
| `mfa_invalid` | 401 | invalid multi-factor authentication code |
| &nbsp;&nbsp;&nbsp;&nbsp;`backup_code_invalid` | 401 | backup code is invalid or already used |
| &nbsp;&nbsp;&nbsp;&nbsp;`totp_code_invalid` | 401 | authenticator code is invalid |
| `session_expired` | 401 | session has expired |
| `signature_invalid` | 401 | signature is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`replay_detected` | 401 | request was already received or is too old |
| &nbsp;&nbsp;&nbsp;&nbsp;`request_signature_invalid` | 401 | request signature is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`webhook_signature_invalid` | 401 | webhook signature is invalid |
| `token_expired` | 401 | token has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`access_token_expired` | 401 | access token has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`refresh_token_expired` | 401 | refresh token has expired |
| `token_invalid` | 401 | token is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`refresh_token_invalid` | 401 | refresh token is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`refresh_token_reused` | 401 | refresh token was already used |
| &nbsp;&nbsp;&nbsp;&nbsp;`token_audience_invalid` | 401 | token audience is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`token_issuer_invalid` | 401 | token issuer is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`token_malformed` | 401 | token is malformed |
| &nbsp;&nbsp;&nbsp;&nbsp;`token_not_yet_valid` | 401 | token is not valid yet |
| &nbsp;&nbsp;&nbsp;&nbsp;`token_signature_invalid` | 401 | token signature is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`token_type_invalid` | 401 | wrong token type for this request |
| `token_revoked` | 401 | token has been revoked |
| &nbsp;&nbsp;&nbsp;&nbsp;`refresh_token_revoked` | 401 | refresh token has been revoked |
| `unauthenticated` | 401 | authentication required |
| &nbsp;&nbsp;&nbsp;&nbsp;`api_key_missing` | 401 | API key is missing |
| &nbsp;&nbsp;&nbsp;&nbsp;`token_missing` | 401 | authentication token is missing |
| `payment_required` | 402 | payment required |
| &nbsp;&nbsp;&nbsp;&nbsp;`insufficient_funds` | 402 | insufficient funds |
| &nbsp;&nbsp;&nbsp;&nbsp;`payment_declined` | 402 | payment was declined |
| &nbsp;&nbsp;&nbsp;&nbsp;`payment_method_required` | 402 | a payment method is required |
| &nbsp;&nbsp;&nbsp;&nbsp;`subscription_expired` | 402 | subscription has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`subscription_required` | 402 | an active subscription is required |
| &nbsp;&nbsp;&nbsp;&nbsp;`trial_expired` | 402 | free trial has ended |
| `feature_unavailable` | 403 | feature is not available |
| &nbsp;&nbsp;&nbsp;&nbsp;`feature_disabled` | 403 | feature is disabled |
| &nbsp;&nbsp;&nbsp;&nbsp;`plan_upgrade_required` | 403 | your plan does not include this feature |
| `forbidden` | 403 | you do not have access to this resource |
| &nbsp;&nbsp;&nbsp;&nbsp;`account_disabled` | 403 | account is disabled |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`account_banned` | 403 | account is banned |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`account_deactivated` | 403 | account is deactivated |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`account_suspended` | 403 | account is suspended |
| &nbsp;&nbsp;&nbsp;&nbsp;`account_locked` | 403 | account is locked |
| &nbsp;&nbsp;&nbsp;&nbsp;`account_not_verified` | 403 | account is not verified |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`email_not_verified` | 403 | email address is not verified |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`phone_not_verified` | 403 | phone number is not verified |
| &nbsp;&nbsp;&nbsp;&nbsp;`captcha_required` | 403 | captcha required |
| &nbsp;&nbsp;&nbsp;&nbsp;`cross_tenant_access` | 403 | access to another tenant's data is not allowed |
| &nbsp;&nbsp;&nbsp;&nbsp;`csrf_invalid` | 403 | CSRF token is missing or invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`insufficient_permissions` | 403 | insufficient permissions |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`insufficient_role` | 403 | your role does not allow this action |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`insufficient_scope` | 403 | token does not have the required scope |
| &nbsp;&nbsp;&nbsp;&nbsp;`ip_not_allowed` | 403 | access from this IP address is not allowed |
| &nbsp;&nbsp;&nbsp;&nbsp;`mfa_required` | 403 | multi-factor authentication required |
| &nbsp;&nbsp;&nbsp;&nbsp;`not_resource_owner` | 403 | you do not own this resource |
| &nbsp;&nbsp;&nbsp;&nbsp;`organization_suspended` | 403 | organization is suspended |
| &nbsp;&nbsp;&nbsp;&nbsp;`origin_not_allowed` | 403 | origin is not allowed |
| &nbsp;&nbsp;&nbsp;&nbsp;`self_action_forbidden` | 403 | you cannot perform this action on your own account |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`cannot_modify_own_role` | 403 | you cannot change your own role |
| &nbsp;&nbsp;&nbsp;&nbsp;`tenant_disabled` | 403 | tenant is disabled |
| &nbsp;&nbsp;&nbsp;&nbsp;`tenant_mismatch` | 403 | resource belongs to a different tenant |
| `limit_exceeded` | 403 | limit exceeded |
| &nbsp;&nbsp;&nbsp;&nbsp;`api_key_limit_reached` | 403 | maximum number of API keys reached |
| &nbsp;&nbsp;&nbsp;&nbsp;`member_limit_reached` | 403 | member limit reached |
| &nbsp;&nbsp;&nbsp;&nbsp;`session_limit_reached` | 403 | maximum number of active sessions reached |
| `quota_exceeded` | 403 | quota exceeded |
| &nbsp;&nbsp;&nbsp;&nbsp;`daily_quota_exceeded` | 403 | daily quota exceeded |
| &nbsp;&nbsp;&nbsp;&nbsp;`monthly_quota_exceeded` | 403 | monthly quota exceeded |
| &nbsp;&nbsp;&nbsp;&nbsp;`storage_quota_exceeded` | 403 | storage quota exceeded |
| `not_found` | 404 | resource not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`api_key_not_found` | 404 | API key not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`file_not_found` | 404 | file not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`invite_not_found` | 404 | invitation not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`member_not_found` | 404 | member not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`organization_not_found` | 404 | organization not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`role_not_found` | 404 | role not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`route_not_found` | 404 | endpoint not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`session_not_found` | 404 | session not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`tenant_not_found` | 404 | tenant not found |
| &nbsp;&nbsp;&nbsp;&nbsp;`user_not_found` | 404 | user not found |
| `method_not_allowed` | 405 | method not allowed |
| `not_acceptable` | 406 | no acceptable representation available |
| `request_timeout` | 408 | request timed out |
| `conflict` | 409 | request conflicts with the current state of the resource |
| &nbsp;&nbsp;&nbsp;&nbsp;`already_exists` | 409 | resource already exists |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_registered` | 409 | already registered |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`email_already_registered` | 409 | email is already registered |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`phone_already_registered` | 409 | phone number is already registered |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`username_already_registered` | 409 | username is already registered |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`file_already_exists` | 409 | a file with this name already exists |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invite_already_exists` | 409 | an invitation is already pending for this user |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`member_already_exists` | 409 | user is already a member |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`name_taken` | 409 | name is already taken |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`organization_already_exists` | 409 | organization already exists |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`role_already_exists` | 409 | role already exists |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`slug_taken` | 409 | slug is already taken |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`social_account_already_linked` | 409 | this social account is already linked to another user |
| &nbsp;&nbsp;&nbsp;&nbsp;`concurrent_modification` | 409 | resource was modified by another request |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`version_conflict` | 409 | resource version does not match |
| &nbsp;&nbsp;&nbsp;&nbsp;`has_dependents` | 409 | resource has dependent resources |
| &nbsp;&nbsp;&nbsp;&nbsp;`in_use` | 409 | resource is in use and cannot be modified |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`role_in_use` | 409 | role is assigned to members and cannot be deleted |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_state` | 409 | resource is in an invalid state for this operation |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_accepted` | 409 | already accepted |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invite_already_accepted` | 409 | invitation was already accepted |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_active` | 409 | resource is already active |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_archived` | 409 | already archived |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_canceled` | 409 | already canceled |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_completed` | 409 | already completed |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_deleted` | 409 | already deleted |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`file_already_deleted` | 409 | file was already deleted |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_disabled` | 409 | already disabled |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_enabled` | 409 | already enabled |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_processed` | 409 | already processed |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`already_verified` | 409 | already verified |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`email_already_verified` | 409 | email is already verified |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`phone_already_verified` | 409 | phone number is already verified |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_transition` | 409 | invalid state transition |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`not_active` | 409 | resource is not active |
| &nbsp;&nbsp;&nbsp;&nbsp;`request_in_progress` | 409 | an identical request is already in progress |
| `expired` | 410 | resource has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`file_expired` | 410 | file has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`invite_expired` | 410 | invitation has expired |
| &nbsp;&nbsp;&nbsp;&nbsp;`link_expired` | 410 | link has expired |
| `gone` | 410 | resource is no longer available |
| &nbsp;&nbsp;&nbsp;&nbsp;`account_deleted` | 410 | account has been deleted |
| &nbsp;&nbsp;&nbsp;&nbsp;`api_version_sunset` | 410 | this API version is no longer available |
| &nbsp;&nbsp;&nbsp;&nbsp;`invite_revoked` | 410 | invitation has been revoked |
| `precondition_failed` | 412 | precondition failed |
| &nbsp;&nbsp;&nbsp;&nbsp;`etag_mismatch` | 412 | ETag does not match the current version |
| `payload_too_large` | 413 | payload too large |
| &nbsp;&nbsp;&nbsp;&nbsp;`file_too_large` | 413 | file is too large |
| &nbsp;&nbsp;&nbsp;&nbsp;`image_too_large` | 413 | image is too large |
| &nbsp;&nbsp;&nbsp;&nbsp;`video_too_large` | 413 | video is too large |
| `unsupported_media_type` | 415 | unsupported media type |
| &nbsp;&nbsp;&nbsp;&nbsp;`unsupported_file_extension` | 415 | file extension is not allowed |
| &nbsp;&nbsp;&nbsp;&nbsp;`unsupported_file_type` | 415 | file type is not allowed |
| &nbsp;&nbsp;&nbsp;&nbsp;`unsupported_image_format` | 415 | image format is not supported |
| `business_rule_violation` | 422 | business rule violation |
| &nbsp;&nbsp;&nbsp;&nbsp;`cannot_delete_last_admin` | 422 | cannot remove the last administrator |
| &nbsp;&nbsp;&nbsp;&nbsp;`cannot_remove_last_owner` | 422 | cannot remove the last owner |
| &nbsp;&nbsp;&nbsp;&nbsp;`profile_incomplete` | 422 | profile must be completed first |
| `idempotency_key_reused` | 422 | idempotency key was already used with a different request |
| `not_supported` | 422 | operation not supported |
| `unprocessable` | 422 | request cannot be processed |
| &nbsp;&nbsp;&nbsp;&nbsp;`file_corrupted` | 422 | file is corrupted or unreadable |
| &nbsp;&nbsp;&nbsp;&nbsp;`file_malware_detected` | 422 | file was rejected by the security scan |
| `validation_failed` | 422 | validation failed |
| &nbsp;&nbsp;&nbsp;&nbsp;`duplicate_value` | 422 | duplicate value |
| &nbsp;&nbsp;&nbsp;&nbsp;`image_dimensions_invalid` | 422 | image dimensions are not allowed |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_date` | 422 | date is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`date_in_future` | 422 | date must not be in the future |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`date_in_past` | 422 | date must not be in the past |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_date_range` | 422 | end date must not be before start date |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_file_name` | 422 | file name is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_format` | 422 | field has an invalid format |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_base64` | 422 | value is not valid base64 |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_country_code` | 422 | country code is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_currency` | 422 | currency code is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_email` | 422 | email address is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_ip_address` | 422 | IP address is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_locale` | 422 | locale is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_phone_number` | 422 | phone number is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_postal_code` | 422 | postal code is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_timezone` | 422 | timezone is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_url` | 422 | URL is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_username` | 422 | username is invalid |
| &nbsp;&nbsp;&nbsp;&nbsp;`invalid_length` | 422 | value has an invalid length |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`value_too_long` | 422 | value is too long |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`value_too_short` | 422 | value is too short |
| &nbsp;&nbsp;&nbsp;&nbsp;`password_reused` | 422 | password was used recently and cannot be reused |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`password_same_as_current` | 422 | new password must differ from the current password |
| &nbsp;&nbsp;&nbsp;&nbsp;`password_too_weak` | 422 | password does not meet the strength requirements |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`password_compromised` | 422 | password appears in a known data breach |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`password_contains_user_info` | 422 | password must not contain your name or email |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`password_too_long` | 422 | password is too long |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`password_too_short` | 422 | password is too short |
| &nbsp;&nbsp;&nbsp;&nbsp;`required` | 422 | field is required |
| &nbsp;&nbsp;&nbsp;&nbsp;`too_few_items` | 422 | too few items |
| &nbsp;&nbsp;&nbsp;&nbsp;`too_many_items` | 422 | too many items |
| &nbsp;&nbsp;&nbsp;&nbsp;`value_not_allowed` | 422 | value is not allowed |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`invalid_enum` | 422 | value is not one of the allowed options |
| &nbsp;&nbsp;&nbsp;&nbsp;`value_out_of_range` | 422 | value is out of range |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`value_too_large` | 422 | value is above the maximum |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`value_too_small` | 422 | value is below the minimum |
| `locked` | 423 | resource is locked |
| `precondition_required` | 428 | precondition header is required |
| `rate_limited` | 429 | too many requests |
| &nbsp;&nbsp;&nbsp;&nbsp;`concurrency_limit_exceeded` | 429 | too many concurrent requests |
| &nbsp;&nbsp;&nbsp;&nbsp;`verification_rate_limited` | 429 | please wait before requesting another verification |
| `too_many_attempts` | 429 | too many attempts, try again later |
| &nbsp;&nbsp;&nbsp;&nbsp;`login_attempts_exceeded` | 429 | too many failed login attempts, try again later |
| &nbsp;&nbsp;&nbsp;&nbsp;`otp_attempts_exceeded` | 429 | too many incorrect codes, try again later |
| &nbsp;&nbsp;&nbsp;&nbsp;`password_reset_attempts_exceeded` | 429 | too many password reset requests, try again later |
| `client_closed_request` | 499 | client closed the request |
| `internal` | 500 | internal server error |
| &nbsp;&nbsp;&nbsp;&nbsp;`configuration_invalid` | 500 | internal server error |
| &nbsp;&nbsp;&nbsp;&nbsp;`database_failed` | 500 | internal server error |
| &nbsp;&nbsp;&nbsp;&nbsp;`decryption_failed` | 500 | internal server error |
| &nbsp;&nbsp;&nbsp;&nbsp;`encryption_failed` | 500 | internal server error |
| &nbsp;&nbsp;&nbsp;&nbsp;`file_processing_failed` | 500 | internal server error |
| &nbsp;&nbsp;&nbsp;&nbsp;`serialization_failed` | 500 | internal server error |
| `not_implemented` | 501 | not implemented |
| `dependency_failed` | 502 | upstream service failed |
| &nbsp;&nbsp;&nbsp;&nbsp;`dependency_response_invalid` | 502 | upstream service returned an invalid response |
| &nbsp;&nbsp;&nbsp;&nbsp;`email_provider_failed` | 502 | email provider failed |
| &nbsp;&nbsp;&nbsp;&nbsp;`oauth_provider_failed` | 502 | identity provider failed |
| &nbsp;&nbsp;&nbsp;&nbsp;`payment_provider_failed` | 502 | payment provider failed |
| &nbsp;&nbsp;&nbsp;&nbsp;`sms_provider_failed` | 502 | SMS provider failed |
| `service_unavailable` | 503 | service temporarily unavailable |
| &nbsp;&nbsp;&nbsp;&nbsp;`circuit_breaker_open` | 503 | service temporarily unavailable |
| &nbsp;&nbsp;&nbsp;&nbsp;`dependency_unavailable` | 503 | a required service is unavailable |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`cache_unavailable` | 503 | cache is unavailable |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`database_unavailable` | 503 | database is unavailable |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`queue_unavailable` | 503 | message queue is unavailable |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`search_unavailable` | 503 | search is unavailable |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`storage_unavailable` | 503 | file storage is unavailable |
| &nbsp;&nbsp;&nbsp;&nbsp;`maintenance` | 503 | service is under maintenance |
| &nbsp;&nbsp;&nbsp;&nbsp;`overloaded` | 503 | service is overloaded, try again later |
| &nbsp;&nbsp;&nbsp;&nbsp;`shutting_down` | 503 | service is shutting down |
| `gateway_timeout` | 504 | upstream service timed out |
| &nbsp;&nbsp;&nbsp;&nbsp;`database_timeout` | 504 | database timed out |
| &nbsp;&nbsp;&nbsp;&nbsp;`external_service_timeout` | 504 | external service timed out |
