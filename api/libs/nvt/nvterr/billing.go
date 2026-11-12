package nvterr

import (
	"net/http"
)

// Domain "billing": payment-required (402) situations.

var (
	ErrPaymentRequired       = define("payment_required", http.StatusPaymentRequired, "payment required")
	ErrSubscriptionRequired  = defineChild(ErrPaymentRequired, "subscription_required", "an active subscription is required")
	ErrSubscriptionExpired   = defineChild(ErrPaymentRequired, "subscription_expired", "subscription has expired")
	ErrTrialExpired          = defineChild(ErrPaymentRequired, "trial_expired", "free trial has ended")
	ErrPaymentMethodRequired = defineChild(ErrPaymentRequired, "payment_method_required", "a payment method is required")
	ErrPaymentDeclined       = defineChild(ErrPaymentRequired, "payment_declined", "payment was declined")
	ErrInsufficientFunds     = defineChild(ErrPaymentRequired, "insufficient_funds", "insufficient funds")
)
