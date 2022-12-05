package types

const (
	EventTypeDeposit        = "deposit"
	EventTypeWithdraw       = "withdraw"
	EventTypeClaim          = "claim"
	EventTypeCompleteUnbond = "complete_unbond"

	AttributeKeyValidator      = "validator"
	AttributeKeyDepositor      = "depositor"
	AttributeKeyToken          = "token"
	AttributeKeyNewShares      = "new_shares"
	AttributeKeyCompletionTime = "completion_time"
	AttributeValueCategory     = ModuleName

	AttributeKeyDelegator = "delegator"
)
