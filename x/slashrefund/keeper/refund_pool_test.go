package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

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
		valPubk := secp256k1.GenPrivKey().PubKey()
		valAddr := sdk.ValAddress(valPubk.Address())
		items[i].OperatorAddress = valAddr.String()
		items[i].Shares = sdk.ZeroDec()
		items[i].Tokens.Amount = sdk.ZeroInt()
		items[i].Tokens.Denom = keeper.AllowedTokensList(ctx)[0]
		keeper.SetRefundPool(ctx, items[i])
	}
	return items
}

func TestRefundPoolGet(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNRefundPool(keeper, ctx, 10)
	for _, item := range items {
		valAddr, _ := sdk.ValAddressFromBech32(item.OperatorAddress)
		rst, found := keeper.GetRefundPool(ctx, valAddr)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestUpdateRefundPool(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	refPools := createNRefundPool(keeper, ctx, 10)
	for i, refPool := range refPools {
		valAddr, _ := sdk.ValAddressFromBech32(refPool.OperatorAddress)

		refPool.Tokens.Amount = refPool.Tokens.Amount.Add(sdk.NewInt(int64(i * 1000)))
		refPool.Shares = refPool.Shares.Add(sdk.NewDec(int64(i * 2000)))
		keeper.SetRefundPool(ctx, refPool)

		rst, found := keeper.GetRefundPool(ctx, valAddr)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&refPool),
			nullify.Fill(&rst),
		)
	}
}

func TestRefundPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNRefundPool(keeper, ctx, 10)
	for _, item := range items {
		valAddr, _ := sdk.ValAddressFromBech32(item.OperatorAddress)
		keeper.RemoveRefundPool(ctx, valAddr)
		_, found := keeper.GetRefundPool(ctx, valAddr)
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
