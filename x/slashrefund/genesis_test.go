package slashrefund_test

import (
	"testing"

	keepertest "github.com/made-in-block/slash-refund/testutil/keeper"
	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		DepositList: []types.Deposit{
			{
				DepositorAddress: "0",
				ValidatorAddress: "0",
			},
			{
				DepositorAddress: "1",
				ValidatorAddress: "1",
			},
		},
		UnbondingDepositList: []types.UnbondingDeposit{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		UnbondingDepositCount: 2,
		DepositPoolList: []types.DepositPool{
			{
				OperatorAddress: "0",
			},
			{
				OperatorAddress: "1",
			},
		},
		UnbondingDepositList: []types.UnbondingDeposit{
			{
				DelegatorAddress: "0",
				ValidatorAddress: "0",
			},
			{
				DelegatorAddress: "1",
				ValidatorAddress: "1",
			},
		},
		RefundPoolList: []types.RefundPool{
			{
				OperatorAddress: "0",
			},
			{
				OperatorAddress: "1",
			},
		},
		RefundList: []types.Refund{
			{
				Delegator: "0",
				Validator: "0",
			},
			{
				Delegator: "1",
				Validator: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SlashrefundKeeper(t)
	slashrefund.InitGenesis(ctx, *k, genesisState)
	got := slashrefund.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.DepositList, got.DepositList)
	require.ElementsMatch(t, genesisState.UnbondingDepositList, got.UnbondingDepositList)
	require.Equal(t, genesisState.UnbondingDepositCount, got.UnbondingDepositCount)
	require.ElementsMatch(t, genesisState.DepositPoolList, got.DepositPoolList)
	require.ElementsMatch(t, genesisState.UnbondingDepositList, got.UnbondingDepositList)
	require.ElementsMatch(t, genesisState.RefundPoolList, got.RefundPoolList)
	require.ElementsMatch(t, genesisState.RefundList, got.RefundList)
	// this line is used by starport scaffolding # genesis/test/assert
}
