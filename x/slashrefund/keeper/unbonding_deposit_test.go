package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/made-in-block/slash-refund/testutil/keeper"
	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNUnbondingDeposit(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.UnbondingDeposit {
	items := make([]types.UnbondingDeposit, n)
	for i := range items {
		items[i].DepositorAddress = strconv.Itoa(i)
		items[i].ValidatorAddress = strconv.Itoa(i)

		keeper.SetUnbondingDeposit(ctx, items[i])
	}
	return items
}

func TestUnbondingDepositGet(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetUnbondingDeposit(ctx,
			item.DepositorAddress,
			item.ValidatorAddress,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestUnbondingDepositRemove(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveUnbondingDeposit(ctx,
			item.DepositorAddress,
			item.ValidatorAddress,
		)
		_, found := keeper.GetUnbondingDeposit(ctx,
			item.DepositorAddress,
			item.ValidatorAddress,
		)
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
