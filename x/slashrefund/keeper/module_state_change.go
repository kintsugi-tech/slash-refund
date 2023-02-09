package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	//abci "github.com/tendermint/tendermint/abci/types"
)

// This function checks the state of unbonding deposits in the UBDQueue and complete the unbonding
// if the unbodning time is reached.
func (k Keeper) BlockUnbondingDepositUpdates(ctx sdk.Context) []types.DVPair {

	ctxTime := ctx.BlockHeader().Time

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctxTime)
	for _, dvPair := range matureUnbonds {
		validatorAddress, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		depositorAddress := sdk.MustAccAddressFromBech32(dvPair.DepositorAddress)

		balances, err := k.CompleteUnbonding(ctx, ctxTime, depositorAddress, validatorAddress)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbond,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDepositor, dvPair.DepositorAddress),
			),
		)
	}

	return matureUnbonds
}
