package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDRefundPool(validatorAddr sdk.ValAddress, tokens sdk.Coin, shares sdk.Dec) RefundPool {
	return RefundPool{
		OperatorAddress: validatorAddr.String(),
		Tokens:          tokens,
		Shares:          shares,
	}
}
