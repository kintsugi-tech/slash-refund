package slashrefund

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	// slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	logger := k.Logger(ctx)

	logger.Error("Height", "height", ctx.BlockHeight())

	matureUnbonds := k.BlockUnbondingDepositUpdates(ctx)
	if matureUnbonds != nil {
		logger.Error("    found and processed mature unbonds")
	}

	//events := ctx.EventManager().Events()
	//
	//// Iterate all events in this block
	//for _, event := range events {
	//
	//	// Check if we have a slashing event
	//	if event.Type == slashingtypes.EventTypeSlash {
	//
	//		// Iterate attributes to find which validators has been slashed
	//		for _, attr := range event.Attributes {
	//
	//			// Check if validator has a deposit ready to use as refund
	//			if string(attr.GetKey()) == "address" {
	//				validator, _ := k.GetValidatorByConsAddrBytes(ctx, attr.GetValue())
	//				deposits, total := k.GetDepositOfValidator(ctx, validator.GetOperator())
	//
	//				logger.Error("deposits", "dep", len(deposits), "tot", total)
	//
	//				// skip if we don't have any deposit
	//				if len(deposits) == 0 || total.Amount.LTE(sdk.NewInt(0)) {
	//					return
	//				}
	//
	//				// Check how much we should refund
	//
	//				// Refund users
	//
	//				return
	//			}
	//		}
	//	}
	//}
	//
}

func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []types.DVPair {

	logger := k.Logger(ctx)
	logger.Error("HEY!!")
	//logger.Error("End blocker for block %d", ctx.BlockHeight())

	matureUnbonds := k.BlockUnbondingDepositUpdates(ctx)
	if matureUnbonds != nil {
		logger.Error("    found and processed mature unbonds")
	}

	return matureUnbonds
}
