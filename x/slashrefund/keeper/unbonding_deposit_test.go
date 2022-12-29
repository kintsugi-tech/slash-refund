package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func createNEntries(n int) []types.UnbondingDepositEntry {
	creationHeight := n
	completionTime := time.Now().Add(time.Duration(n * 1000))
	balance := sdk.NewInt(1000000)
	var entries []types.UnbondingDepositEntry
	for i := 0; i < n; i++ {
		entry := types.NewUnbondingDepositEntry(int64(creationHeight), completionTime, balance)
		entries = append(entries, entry)
	}
	return entries
}

func createNUnbondingDeposit(keeper *keeper.Keeper, ctx sdk.Context, n int, nentries int) []types.UnbondingDeposit {
	items := make([]types.UnbondingDeposit, n)
	for i := range items {
		depPubk := secp256k1.GenPrivKey().PubKey()
		depAddr := sdk.AccAddress(depPubk.Address())
		valPubk := secp256k1.GenPrivKey().PubKey()
		valAddr := sdk.ValAddress(valPubk.Address())
		items[i].DepositorAddress = depAddr.String()
		items[i].ValidatorAddress = valAddr.String()
		items[i].Entries = createNEntries(nentries)
		keeper.SetUnbondingDeposit(ctx, items[i])
	}
	return items
}

func createNUnbondingDepositForValidator(keeper *keeper.Keeper, ctx sdk.Context, n int, nentries int, valAddr sdk.ValAddress) []types.UnbondingDeposit {
	items := make([]types.UnbondingDeposit, n)
	for i := range items {
		depPubk := secp256k1.GenPrivKey().PubKey()
		depAddr := sdk.AccAddress(depPubk.Address())
		items[i].DepositorAddress = depAddr.String()
		items[i].ValidatorAddress = valAddr.String()
		items[i].Entries = createNEntries(nentries)
		keeper.SetUnbondingDeposit(ctx, items[i])
	}
	return items
}

func TestUnbondingDepositGet(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10, 3)
	for _, item := range items {
		depAddr, _ := sdk.AccAddressFromBech32(item.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		got, found := keeper.GetUnbondingDeposit(ctx, depAddr, valAddr)
		require.True(t, found)
		require.Equal(t, nullify.Fill(&item), nullify.Fill(&got))
	}
}
func TestUnbondingDepositRemove(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10, 3)
	for _, item := range items {
		keeper.RemoveUnbondingDeposit(ctx, item)
		depAddr, _ := sdk.AccAddressFromBech32(item.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		_, found := keeper.GetUnbondingDeposit(ctx, depAddr, valAddr)
		require.False(t, found)
		_, found = keeper.GetUnbondingDepositByValIndexKey(ctx, valAddr, depAddr)
		require.False(t, found)

	}
}

func TestUnbondingDepositGetAll(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10, 3)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllUnbondingDeposit(ctx)),
	)
}

func TestUnbondingDepositGetUnbondingDepositsFromValidator(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	valPubk := secp256k1.GenPrivKey().PubKey()
	valAddr := sdk.ValAddress(valPubk.Address())
	items := createNUnbondingDepositForValidator(keeper, ctx, 10, 3, valAddr)
	got := keeper.GetUnbondingDepositsFromValidator(ctx, valAddr)
	require.ElementsMatch(t, nullify.Fill(items), nullify.Fill(got))
}
