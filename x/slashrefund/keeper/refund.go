package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (keeper Keeper) GetValidatorByConsAddrBytes(ctx sdk.Context, consAddr []byte) (validator stakingtypes.Validator, found bool) {

	// Decode address
	addr, _ := sdk.ConsAddressFromBech32(string(consAddr))

	// TODO: Handle error
	return keeper.stakingKeeper.GetValidatorByConsAddr(ctx, addr)
}

func (keeper Keeper) IterateValidatorDelegations(ctx sdk.Context, consAddr []byte) {

	validator, _ := keeper.GetValidatorByConsAddrBytes(ctx, consAddr)

	// TODO: Handle not found validator
	delegations := keeper.stakingKeeper.GetValidatorDelegations(ctx, validator.GetOperator())

	for _, delegation := range delegations {
		keeper.Logger(ctx).Error("delegation", "addr", delegation.DelegatorAddress, "shares", delegation.Shares)
	}
}

func (k Keeper) ProcessRefunds(ctx sdk.Context, slashEvents []types.SlashEvent) {

	k.Logger(ctx).Error("slashs", "amt", slashEvents[0].Amount, "reas", slashEvents[0].Reason, "val", slashEvents[0].Validator.GetOperator())
	for _, slash := range slashEvents {

		// read deposits for that validators
		deposits, total := k.GetDepositOfValidator(ctx, slash.Validator.GetOperator())

		k.Logger(ctx).Error("deposits", "dep", len(deposits), "tot", total)

		// skip if we don't have any deposit
		if len(deposits) == 0 || total.Amount.LTE(sdk.NewInt(0)) {
			continue
		}

		// Check how much we should refund
		amountToRefund := sdk.MinInt(slash.Amount, total.Amount)

		// Refund users
		k.Logger(ctx).Error("deposits", "refund amt", amountToRefund)
	}

}
