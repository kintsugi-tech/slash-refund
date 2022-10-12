package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/made-in-block/slash-refund/testutil/keeper"
	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func createNUnbondingDeposit(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.UnbondingDeposit {
	items := make([]types.UnbondingDeposit, n)
	for i := range items {
		items[i].Id = keeper.AppendUnbondingDeposit(ctx, items[i])
	}
	return items
}

func TestUnbondingDepositGet(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetUnbondingDeposit(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestUnbondingDepositRemove(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveUnbondingDeposit(ctx, item.Id)
		_, found := keeper.GetUnbondingDeposit(ctx, item.Id)
		require.False(t, found)
	}
}

func TestUnbondingDepositGetAll(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllUnbondingDeposit(ctx)),
	)
}

func TestUnbondingDepositCount(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetUnbondingDepositCount(ctx))
}
