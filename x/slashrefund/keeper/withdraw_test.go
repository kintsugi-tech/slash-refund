package keeper_test

import (

	"testing"
	"time"

	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)



func createNUnbondingDepositForValidator(keeper *keeper.Keeper, ctx sdk.Context, n int, nEntries int, valAddress string) []types.UnbondingDeposit {
	items := make([]types.UnbondingDeposit, n)
	for i := range items {
		items[i].DepositorAddress = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
		items[i].ValidatorAddress = valAddress
		items[i].Entries = testslashrefund.CreateNEntries(nEntries)
		keeper.SetUnbondingDeposit(ctx, items[i])
	}
	return items
}

func Test_GetUnbondingDeposit(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := testslashrefund.CreateNUnbondingDeposit(keeper, ctx, 1, 5)
	for _, item := range items {
		depAddr, _ := sdk.AccAddressFromBech32(item.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		got, found := keeper.GetUnbondingDeposit(ctx, depAddr, valAddr)
		require.True(t, found)
		require.Equal(t, item, got)
	}
}
func Test_RemoveUnbondingDeposit(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := testslashrefund.CreateNUnbondingDeposit(keeper, ctx, 10, 5)
	for _, item := range items {
		keeper.RemoveUnbondingDeposit(ctx, item)
		depAddr, _ := sdk.AccAddressFromBech32(item.DepositorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		_, found := keeper.GetUnbondingDeposit(ctx, depAddr, valAddr)
		require.False(t, found)
	}
}

func Test_GetAllUnbondingDeposit(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := testslashrefund.CreateNUnbondingDeposit(keeper, ctx, 10, 5)
	require.ElementsMatch(t, items, keeper.GetAllUnbondingDeposit(ctx))
}

func Test_GetUnbondingDepositsFromValidator(t *testing.T) {
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

func Test_SetUnbondingDepositEntry(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	depAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
	creatinHeight := int64(10)
	minTime := time.Now().UTC()
	balance := sdk.NewInt(3)
	creatinHeight2 := int64(11)
	minTime2 := minTime.Add(time.Hour*1)
	balance2 := sdk.NewInt(5)

	// Test add entry to not existing unbonding deposit
	_ = keeper.SetUnbondingDepositEntry(ctx, depAddr, valAddr, creatinHeight, minTime, balance)
	ubd, found := keeper.GetUnbondingDeposit(ctx, depAddr, valAddr)
	require.Equal(t, found, true)
	require.Equal(
		t, 
		ubd.Entries,
		[]types.UnbondingDepositEntry{
			types.NewUnbondingDepositEntry(creatinHeight, minTime, balance),
		},
	)

	// Test add a new entry to existing unbonding deposit
	_ = keeper.SetUnbondingDepositEntry(ctx, depAddr, valAddr, creatinHeight2, minTime2, balance2)
	ubd, found = keeper.GetUnbondingDeposit(ctx, depAddr, valAddr)
	require.Equal(t, found, true)
	require.Equal(
		t, 
		ubd.Entries,
		[]types.UnbondingDepositEntry{
			types.NewUnbondingDepositEntry(creatinHeight, minTime, balance),
			types.NewUnbondingDepositEntry(creatinHeight2, minTime2, balance2),
		},
	)
}

func Test_GetUBDQueueTimeSlice(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	depAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
	creatinHeight := int64(10)
	minTime := time.Now().UTC()
	balance := sdk.NewInt(3)

	ubd := types.NewUnbondingDeposit(depAddr, valAddr, creatinHeight, minTime, balance)
	
	// Set a new timeslice to UBD queue
	keeper.InsertUBDQueue(ctx, ubd, minTime)
	dvPairs := keeper.GetUBDQueueTimeSlice(ctx, minTime)
	dvPair := types.DVPair{
		DepositorAddress: depAddr.String(), 
		ValidatorAddress: valAddr.String(),
	}
	require.Equal(t, dvPairs, []types.DVPair{dvPair})

	// Append a new dvPair to already existing timeslcie
	keeper.InsertUBDQueue(ctx, ubd, minTime)
	dvPairs = keeper.GetUBDQueueTimeSlice(ctx, minTime)
	require.Equal(t, dvPairs, []types.DVPair{dvPair, dvPair})
}