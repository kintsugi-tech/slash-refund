package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func createNDeposit(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Deposit {
	items := make([]types.Deposit, n)
	for i := range items {
		depPubk := secp256k1.GenPrivKey().PubKey()
		depAddr := sdk.AccAddress(depPubk.Address())
		valPubk := secp256k1.GenPrivKey().PubKey()
		valAddr := sdk.ValAddress(valPubk.Address())
		items[i].DepositorAddress = depAddr.String()
		items[i].ValidatorAddress = valAddr.String()
		items[i].Shares = sdk.ZeroDec()
		keeper.SetDeposit(ctx, items[i])
	}
	return items
}

func createNDepositForValidator(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.Deposit, sdk.ValAddress) {
	items := make([]types.Deposit, n)
	valPubk := secp256k1.GenPrivKey().PubKey()
	valAddr := sdk.ValAddress(valPubk.Address())
	for i := range items {
		depPubk := secp256k1.GenPrivKey().PubKey()
		depAddr := sdk.AccAddress(depPubk.Address())
		items[i].DepositorAddress = depAddr.String()
		items[i].ValidatorAddress = valAddr.String()
		items[i].Shares = sdk.ZeroDec()
		keeper.SetDeposit(ctx, items[i])
	}
	return items, valAddr
}

func TestDepositGet(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	deposits := createNDeposit(keeper, ctx, 10)
	for _, deposit := range deposits {
		depAddr, _ := sdk.AccAddressFromBech32(deposit.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(deposit.ValidatorAddress)
		rst, found := keeper.GetDeposit(ctx, depAddr, valAddr)
		require.True(t, found)
		require.Equal(t, deposit, rst)
	}
}
func TestDepositRemove(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	deposits := createNDeposit(keeper, ctx, 10)
	for _, deposit := range deposits {
		keeper.RemoveDeposit(ctx, deposit)
		depAddr, _ := sdk.AccAddressFromBech32(deposit.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(deposit.ValidatorAddress)
		_, found := keeper.GetDeposit(ctx, depAddr, valAddr)
		require.False(t, found)
	}
}

func TestDepositGetAll(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNDeposit(keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllDeposit(ctx))
}

func TestGetValidatorDeposits(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items0, valAddr0 := createNDepositForValidator(keeper, ctx, 10)
	items1, valAddr1 := createNDepositForValidator(keeper, ctx, 10)
	require.ElementsMatch(t, items0, keeper.GetValidatorDeposits(ctx, valAddr0))
	require.ElementsMatch(t, items1, keeper.GetValidatorDeposits(ctx, valAddr1))
}
