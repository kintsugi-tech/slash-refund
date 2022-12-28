package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/made-in-block/slash-refund/testutil/keeper"
	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func createNDepositPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.DepositPool {
	items := make([]types.DepositPool, n)
	for i := range items {
		valPubk := secp256k1.GenPrivKey().PubKey()
		valAddr := sdk.ValAddress(valPubk.Address())
		items[i].OperatorAddress = valAddr.String()
		items[i].Shares = sdk.NewDec(int64(1000 * i))
		keeper.SetDepositPool(ctx, items[i])
	}
	return items
}

func TestDepositPoolGet(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNDepositPool(keeper, ctx, 10)
	for _, item := range items {
		valAddr, _ := sdk.ValAddressFromBech32(item.OperatorAddress)
		rst, found := keeper.GetDepositPool(ctx, valAddr)
		require.True(t, found)
		require.Equal(t, nullify.Fill(&item), nullify.Fill(&rst))
	}
}
func TestDepositPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNDepositPool(keeper, ctx, 10)
	for _, item := range items {
		valAddr, _ := sdk.ValAddressFromBech32(item.OperatorAddress)
		keeper.RemoveDepositPool(ctx, valAddr)
		_, found := keeper.GetDepositPool(ctx, valAddr)
		require.False(t, found)
	}
}

func TestDepositPoolGetAll(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	items := createNDepositPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllDepositPool(ctx)),
	)
}
