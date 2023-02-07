package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func createNEntries(n int) []types.UnbondingDepositEntry {

	var entries []types.UnbondingDepositEntry
	for i := 0; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		r := rand.Int63n(1000000)
		creationHeight := r
		completionTime := time.Now().Add(time.Duration(r)).UTC()
		balance := sdk.NewInt(r)
		initBalance := balance.AddRaw(rand.Int63n(1000000))
		entry := types.NewUnbondingDepositEntry(int64(creationHeight), completionTime, initBalance)
		entry.Balance = balance
		entries = append(entries, entry)
	}
	return entries
}

func createNUnbondingDeposit(keeper *keeper.Keeper, ctx sdk.Context, n int, nentries int) []types.UnbondingDeposit {
	items := make([]types.UnbondingDeposit, n)
	for i := range items {
		items[i].DepositorAddress = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
		items[i].ValidatorAddress = sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
		items[i].Entries = createNEntries(nentries)
		keeper.SetUnbondingDeposit(ctx, items[i])
	}
	return items
}

func createNUnbondingDepositForValidator(keeper *keeper.Keeper, ctx sdk.Context, n int, nentries int, valAddress string) []types.UnbondingDeposit {
	items := make([]types.UnbondingDeposit, n)
	for i := range items {
		items[i].DepositorAddress = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
		items[i].ValidatorAddress = valAddress
		items[i].Entries = createNEntries(nentries)
		keeper.SetUnbondingDeposit(ctx, items[i])
	}
	return items
}

func TestUnbondingDepositGet(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10, 5)
	for _, item := range items {
		depAddr, _ := sdk.AccAddressFromBech32(item.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		got, found := keeper.GetUnbondingDeposit(ctx, depAddr, valAddr)
		require.True(t, found)
		require.Equal(t, item, got)
	}
}
func TestUnbondingDepositRemove(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNUnbondingDeposit(keeper, ctx, 10, 5)
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
	items := createNUnbondingDeposit(keeper, ctx, 10, 5)
	require.ElementsMatch(t, items, keeper.GetAllUnbondingDeposit(ctx))
}

func TestUnbondingDepositGetUnbondingDepositsFromValidator(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	valAddress0 := "cosmosvaloper12h6y5kn64xh6d6wsnw7098kc8k9kp2u6m03was"
	items0 := createNUnbondingDepositForValidator(keeper, ctx, 10, 5, valAddress0)
	valAddress1 := "cosmosvaloper1e8wanntwnsnvrz7eaj82hzhhp5lrz3gsfzkfly"
	items1 := createNUnbondingDepositForValidator(keeper, ctx, 10, 5, valAddress1)
	valAddr0, err := sdk.ValAddressFromBech32(valAddress0)
	require.NoError(t, err)
	valAddr1, err := sdk.ValAddressFromBech32(valAddress1)
	require.NoError(t, err)
	got0 := keeper.GetUnbondingDepositsFromValidator(ctx, valAddr0)
	require.ElementsMatch(t, items0, got0)
	got1 := keeper.GetUnbondingDepositsFromValidator(ctx, valAddr1)
	require.ElementsMatch(t, items1, got1)
}
