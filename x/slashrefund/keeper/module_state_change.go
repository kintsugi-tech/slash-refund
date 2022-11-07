package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	//abci "github.com/tendermint/tendermint/abci/types"
)

// BlockUnbondingDepositUpdates check state of unbonding deposits in the UBDQueue
func (k Keeper) BlockUnbondingDepositUpdates(ctx sdk.Context) []types.DVPair {

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		validatorAddress, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		depositorAddress := sdk.MustAccAddressFromBech32(dvPair.DepositorAddress)

		balances, err := k.CompleteUnbonding(ctx, depositorAddress, validatorAddress)
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
