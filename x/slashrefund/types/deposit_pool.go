package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDepositPool(validatorAddr sdk.ValAddress, tokens sdk.Coin, shares sdk.Dec) DepositPool {
	return DepositPool{
		OperatorAddress: validatorAddr.String(),
		Tokens: tokens,
		Shares: shares,
	}
}

func (d DepositPool) SharesFromTokens(tokens sdk.Coin) sdk.Dec {
	// TODO: manage error
	return d.Shares.MulInt(tokens.Amount).QuoInt(d.GetTokens().Amount)
}