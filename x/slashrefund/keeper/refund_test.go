package keeper_test

import (
	//"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func createNRefund(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Refund {

	items := make([]types.Refund, n)
	for i := range items {
		delPubk := secp256k1.GenPrivKey().PubKey()
		delAddr := sdk.AccAddress(delPubk.Address())
		valPubk := secp256k1.GenPrivKey().PubKey()
		valAddr := sdk.ValAddress(valPubk.Address())
		items[i].DelegatorAddress = delAddr.String()
		items[i].ValidatorAddress = valAddr.String()
		items[i].Shares = sdk.NewDec(int64(1000 * i))
		keeper.SetRefund(ctx, items[i])
	}
	return items
}

func TestRefundGet(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefund(keeper, ctx, 10)
	for _, item := range items {
		delAddr, _ := sdk.AccAddressFromBech32(item.DelegatorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		rst, found := keeper.GetRefund(ctx, delAddr, valAddr)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestRefundRemove(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefund(keeper, ctx, 10)
	for _, item := range items {
		delAddr, _ := sdk.AccAddressFromBech32(item.DelegatorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		refund, found := keeper.GetRefund(ctx, delAddr, valAddr)
		keeper.RemoveRefund(ctx, refund)
		_, found = keeper.GetRefund(ctx, delAddr, valAddr)
		require.False(t, found)
	}
}

func TestRefundGetAll(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefund(keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllRefund(ctx))
}
