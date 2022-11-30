package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (r Refund) GetValidatorAddr() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(r.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (r *Refund) GetShares() sdk.Dec {
	if r != nil {
		return r.Shares
	}
	return sdk.Dec{}
}

func NewRefund(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdk.Dec) Refund {
	return Refund{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           shares,
	}
}
