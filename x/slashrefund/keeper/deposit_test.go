package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

// -------------------------------------------------------------------------------------------------
// Test deposit
// -------------------------------------------------------------------------------------------------

func TestGetDeposit(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	deposits := testslashrefund.CreateNDeposit(keeper, ctx, 10)
	for _, deposit := range deposits {
		depAddr, _ := sdk.AccAddressFromBech32(deposit.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(deposit.ValidatorAddress)
		rst, found := keeper.GetDeposit(ctx, depAddr, valAddr)
		require.True(t, found)
		require.Equal(t, deposit, rst)
	}
}

func TestRemoveDeposit(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	deposits := testslashrefund.CreateNDeposit(keeper, ctx, 10)
	for _, deposit := range deposits {
		keeper.RemoveDeposit(ctx, deposit)
		depAddr, _ := sdk.AccAddressFromBech32(deposit.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(deposit.ValidatorAddress)
		_, found := keeper.GetDeposit(ctx, depAddr, valAddr)
		require.False(t, found)
	}
}

func TestGetAllDeposit(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := testslashrefund.CreateNDeposit(keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllDeposit(ctx))
}

func TestGetValidatorDeposits(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items0, valAddr0 := testslashrefund.CreateNDepositForValidator(keeper, ctx, 10)
	items1, valAddr1 := testslashrefund.CreateNDepositForValidator(keeper, ctx, 10)
	require.ElementsMatch(t, items0, keeper.GetValidatorDeposits(ctx, valAddr0))
	require.ElementsMatch(t, items1, keeper.GetValidatorDeposits(ctx, valAddr1))
}

// -------------------------------------------------------------------------------------------------
// Test deposit pool
// -------------------------------------------------------------------------------------------------

func TestGetDepositPool(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := testslashrefund.CreateNDepositPool(keeper, ctx, 10)
	for _, item := range items {
		valAddr, _ := sdk.ValAddressFromBech32(item.OperatorAddress)
		rst, found := keeper.GetDepositPool(ctx, valAddr)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}

func TestRemoveDepositPool(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := testslashrefund.CreateNDepositPool(keeper, ctx, 10)
	for _, item := range items {
		valAddr, _ := sdk.ValAddressFromBech32(item.OperatorAddress)
		keeper.RemoveDepositPool(ctx, valAddr)
		_, found := keeper.GetDepositPool(ctx, valAddr)
		require.False(t, found)
	}
}

func TestGetAllDepositPool(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := testslashrefund.CreateNDepositPool(keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllDepositPool(ctx))
}

func TestAddDepPoolTokensAndShares(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	depPool := testslashrefund.CreateNDepositPool(keeper, ctx, 1)[0]
	tokensToAdd := sdk.NewInt(100) 
	keeper.AddDepPoolTokensAndShares(
		ctx, 
		depPool, 
		sdk.NewCoin(types.DefaultAllowedTokens[0], tokensToAdd),
	)
	valAddr, err := sdk.ValAddressFromBech32(depPool.OperatorAddress)
	require.NoError(t, err)
	foundDepPool, found := keeper.GetDepositPool(ctx, valAddr)
	require.True(t, found)
	updatedTokensAmount := depPool.Tokens.AddAmount(tokensToAdd).Amount
	updatedShares := depPool.Shares.Add(sdk.NewDecFromInt(tokensToAdd))

	require.Equal(t, foundDepPool.Tokens.Amount, updatedTokensAmount)
	require.Equal(t, foundDepPool.Shares, updatedShares)
}

func TestRemoveDepPoolTokensAndShares(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	// The first pool has 0 tokens and shares
	depPool := testslashrefund.CreateNDepositPool(keeper, ctx, 3)[2]
	sharesToRemove := sdk.NewInt(100)
	keeper.RemoveDepPoolTokensAndShares(
		ctx, 
		depPool, 
		sdk.NewDecFromInt(sharesToRemove),
	)
	valAddr, err := sdk.ValAddressFromBech32(depPool.OperatorAddress)
	require.NoError(t, err)
	foundDepPool, found := keeper.GetDepositPool(ctx, valAddr)
	require.True(t, found)
	updatedTokensAmount := depPool.Tokens.SubAmount(sharesToRemove).Amount
	updatedShares := depPool.Shares.Sub(sdk.NewDecFromInt(sharesToRemove))

	require.Equal(t, foundDepPool.Tokens.Amount, updatedTokensAmount)
	require.Equal(t, foundDepPool.Shares, updatedShares)
}