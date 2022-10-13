package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type SlashEvent struct {
	Validator stakingtypes.Validator
	Amount    sdk.Int
	Reason    string
}
