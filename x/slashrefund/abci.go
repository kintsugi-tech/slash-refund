package slashrefund

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	// Get events and iterate through them in order to get the slashing event.
	events := ctx.EventManager().Events()

	for _, event := range events {

		// Check if we have a slashing event. Skip the jail event, that also is emitted
		// by the slashing module as a slashing event, but with different attributes
		// (first attribute of the jail event is "jailed").
		// TODO: handle jail slash event better.
		if event.Type == slashingtypes.EventTypeSlash && string(event.Attributes[0].GetKey()) != "jailed" {
			k.HandleRefundsFromSlash(ctx, event)
		}
	}
}

func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []types.DVPair {

	// Handle unbonding dequeue
	matureUnbonds := k.BlockUnbondingDepositUpdates(ctx)

	// TODO: Handle removed validator's deposit

	return matureUnbonds
}
