package finout

import "github.com/joomcode/errorx"

var (
	FinoutError                   = errorx.NewNamespace("azure")
	InvalidAuthMethodForAction    = FinoutError.NewType("invalid_auth_method_for_action")
	InvalidAuthMethodForActionMsg = "Auth method does not supported the used action. Please try a different one"
)
