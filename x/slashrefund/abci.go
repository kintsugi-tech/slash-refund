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

	var slashEvents []types.SlashEvent

	events := ctx.EventManager().Events()

	// Iterate all events in this block
	for _, event := range events {

		// Check if we have a slashing event
		if event.Type == slashingtypes.EventTypeSlash {

			slashEvent := types.SlashEvent{}

			// Iterate attributes to fill event details in list
			for _, attr := range event.Attributes {

				// Convert validtor address
				if string(attr.GetKey()) == "address" {
					validator, _ := k.GetValidatorByConsAddrBytes(ctx, attr.GetValue())
					// TODO: handle not ok
					slashEvent.Validator = validator
				}

				// Convert slashed amount
				if string(attr.GetKey()) == "burned_coins" {
					amount, _ := sdk.NewIntFromString(string(attr.GetValue()))

					// TODO handle not ok
					slashEvent.Amount = amount
				}

				// Copy reason
				if string(attr.GetKey()) == "reason" {
					slashEvent.Reason = string(attr.GetValue())
				}
			}

			// append to the list
			slashEvents = append(slashEvents, slashEvent)
		}
	}

	// Process refunds
	if len(slashEvents) > 0 {
		k.ProcessRefunds(ctx, slashEvents)
	}
}
