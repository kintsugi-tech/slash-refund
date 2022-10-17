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

func createNDeposit(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Deposit {
	items := make([]types.Deposit, n)
	for i := range items {
		items[i].DepositorAddress = strconv.Itoa(i)
		items[i].ValidatorAddress = strconv.Itoa(i)

		keeper.SetDeposit(ctx, items[i])
	}
	return items
}

func TestDepositGet(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNDeposit(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetDeposit(ctx,
			sdk.AccAddress(item.DepositorAddress),
			sdk.ValAddress(item.ValidatorAddress),
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestDepositRemove(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNDeposit(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveDeposit(ctx,
			items[0],
		)
		_, found := keeper.GetDeposit(ctx,
			sdk.AccAddress(item.DepositorAddress),
			sdk.ValAddress(item.ValidatorAddress),
		)
		require.False(t, found)
	}
}

func TestDepositGetAll(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNDeposit(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllDeposit(ctx)),
	)
}
