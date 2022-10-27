package types

const (
	EventTypeDeposit  = "deposit"
	EventTypeWithdraw = "withdraw"

	AttributeKeyValidator      = "validator"
	AttributeKeyDepositor      = "depositor"
	AttributeKeyToken          = "token"
	AttributeKeyNewShares      = "new_shares"
	AttributeKeyCompletionTime = "completion_time"
	AttributeValueCategory     = ModuleName

	EventTypeCompleteUnbond = "complete_unbond"
)
