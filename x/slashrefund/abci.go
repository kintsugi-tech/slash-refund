package slashrefund

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	// Handle slashing event
	events := ctx.EventManager().Events()

	// Iterate all events in this block
	for _, event := range events {
		// Check if we have a slashing event and that the event is not coming from a jail action
		if event.Type == slashingtypes.EventTypeSlash && string(event.Attributes[0].GetKey()) != "jailed" {
			// TODO: handle jail slash event better.
			_, err := k.HandleRefundsFromSlash(ctx, event)
			if err != nil {
				// TODO: handle the error
			}
		}
	}
}

func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []types.DVPair {

	// Handle unbonding dequeue
	matureUnbonds := k.BlockUnbondingDepositUpdates(ctx)

	// TODO: Handle removed validator's deposit

	return matureUnbonds
}
