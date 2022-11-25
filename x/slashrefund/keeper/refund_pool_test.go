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

func createNRefundPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.RefundPool {
	items := make([]types.RefundPool, n)
	for i := range items {
		items[i].OperatorAddress = strconv.Itoa(i)

		keeper.SetRefundPool(ctx, items[i])
	}
	return items
}

func TestRefundPoolGet(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNRefundPool(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetRefundPool(ctx,
			item.OperatorAddress,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestRefundPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNRefundPool(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveRefundPool(ctx,
			item.OperatorAddress,
		)
		_, found := keeper.GetRefundPool(ctx,
			item.OperatorAddress,
		)
		require.False(t, found)
	}
}

func TestRefundPoolGetAll(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNRefundPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllRefundPool(ctx)),
	)
}
