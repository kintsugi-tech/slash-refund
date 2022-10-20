package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (d Deposit) GetValidatorAddr() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func NewDeposit(depositorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdk.Dec) Deposit {
	return Deposit{
		DepositorAddress: depositorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           shares,
	}
}