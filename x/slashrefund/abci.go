package slashrefund

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	logger := k.Logger(ctx)

	logger.Error("Height", "height", ctx.BlockHeight())

}

func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []types.DVPair {

	logger := k.Logger(ctx)
	logger.Error("|_ End blocker")

	//Handle unbonding dequeue
	matureUnbonds := k.BlockUnbondingDepositUpdates(ctx)
	if matureUnbonds != nil {
		logger.Error("    |_ found and processed mature unbonds")
	}

	//Handle slashing event
	events := ctx.EventManager().Events()

	// Iterate all events in this block
	for _, event := range events {

		// Check if we have a slashing event
		if event.Type == slashingtypes.EventTypeSlash {
			k.HandleRefundsFromSlash(ctx, event)
		}
	}

	return matureUnbonds
}
