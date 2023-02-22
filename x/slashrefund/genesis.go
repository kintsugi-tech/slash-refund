package slashrefund

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// Initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all deposits
	for _, elem := range genState.DepositList {
		k.SetDeposit(ctx, elem)
	}
	// Set all deposit pools
	for _, elem := range genState.DepositPoolList {
		k.SetDepositPool(ctx, elem)
	}
	// Set all unbonding deposits
	for _, elem := range genState.UnbondingDepositList {
		k.SetUnbondingDeposit(ctx, elem)
	}
	// Set all refund pools
	for _, elem := range genState.RefundPoolList {
		k.SetRefundPool(ctx, elem)
	}
	// Set all refunds
	for _, elem := range genState.RefundList {
		k.SetRefund(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// Returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.DepositList = k.GetAllDeposit(ctx)
	genesis.DepositPoolList = k.GetAllDepositPool(ctx)
	genesis.UnbondingDepositList = k.GetAllUnbondingDeposit(ctx)
	genesis.RefundPoolList = k.GetAllRefundPool(ctx)
	genesis.RefundList = k.GetAllRefund(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
