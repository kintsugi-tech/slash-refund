package slashrefund

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	logger := k.Logger(ctx)

	logger.Error("Begin Blocker")

	events := ctx.EventManager().Events()

	logger.Error("B Events:", "len", len(events))

	for _, event := range events {
		logger.Error("Ricevuto evento", "type", event.Type)

		if event.Type == slashingtypes.EventTypeSlash {
			logger.Error("Attributi", "attr", event.Attributes)
			for _, attr := range event.Attributes {
				logger.Error("Attribute", "key", attr.GetKey(), "value", attr.GetValue())
			}
		}
	}

}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

	logger := k.Logger(ctx)

	logger.Error("End Blocker")

	events := ctx.EventManager().Events()

	logger.Error("Events: ", "len", len(events))

	for _, event := range events {
		logger.Error("Ricevuto evento", "type", event.Type)
	}

	return
}
