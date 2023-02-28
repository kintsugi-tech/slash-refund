package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (d Deposit) MustGetValidatorAddr() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (d Deposit) MustGetDepositorAddr() sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(d.DepositorAddress)
	return addr
}

func (d *Deposit) GetShares() sdk.Dec {
	if d != nil {
		return d.Shares
	}
	return sdk.Dec{}
}

func NewDeposit(depositorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdk.Dec) Deposit {
	return Deposit{
		DepositorAddress: depositorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           shares,
	}
}

// -------------------------------------------------------------------------------------------------
// Deposit pool
// -------------------------------------------------------------------------------------------------

func NewDepositPool(validatorAddr sdk.ValAddress, tokens sdk.Coin, shares sdk.Dec) DepositPool {
	return DepositPool{
		OperatorAddress: validatorAddr.String(),
		Tokens:          tokens,
		Shares:          shares,
	}
}

// Returns the amount of shares given an amount of tokens through a proportion
//
//            pool_shares
// shares = --------------- * tokens
//            pool_tokens
//
func (d DepositPool) SharesFromTokens(tokens sdk.Coin) (sdk.Dec, error) {
	if d.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientTokens
	}
	return d.Shares.MulInt(tokens.Amount).QuoInt(d.GetTokens().Amount), nil
}

func (d DepositPool) SharesFromTokensTruncated(tokens sdk.Coin) (sdk.Dec, error) {
	if d.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientTokens
	}
	return d.Shares.MulInt(tokens.Amount).QuoTruncate(sdk.NewDecFromInt(d.GetTokens().Amount)), nil
}

func (d DepositPool) TokensFromShares(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(d.Tokens.Amount)).Quo(d.Shares)
}