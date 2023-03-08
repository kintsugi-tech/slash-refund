package testslashrefund

import (
	"testing"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/stretchr/testify/require"
)

// Helper is a structure which wraps the slashrefund message server
// and provides methods useful in tests
type Helper struct {
	t       *testing.T
	k       keeper.Keeper
	msgSrvr types.MsgServer
	ctx     sdk.Context
}

func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper) *Helper {
	return &Helper{t, k, keeper.NewMsgServerImpl(k), ctx}
}

func (srh *Helper) Deposit(depAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Int) {
	coin := sdk.NewCoin(srh.k.AllowedTokens(srh.ctx)[0], amount)
	msg := types.NewMsgDeposit(depAddr.String(), valAddr.String(), coin)
	res, err := srh.msgSrvr.Deposit(sdk.WrapSDKContext(srh.ctx), msg)
	require.NoError(srh.t, err)
	require.NotNil(srh.t, res)
}

func (srh *Helper) Withdraw(depAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Int) {
	coin := sdk.NewCoin(srh.k.AllowedTokens(srh.ctx)[0], amount)
	msg := types.NewMsgWithdraw(depAddr.String(), valAddr.String(), coin)
	res, err := srh.msgSrvr.Withdraw(sdk.WrapSDKContext(srh.ctx), msg)
	require.NoError(srh.t, err)
	require.NotNil(srh.t, res)
}

// createNDeposit creates N random deposits.
func CreateNDeposit(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Deposit {
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

// createNDepositForValidator creates N random deposits for a single validator.
func CreateNDepositForValidator(
	keeper *keeper.Keeper, 
	ctx sdk.Context, 
	n int,
) ([]types.Deposit, sdk.ValAddress) {
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

func CreateNDepositPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.DepositPool {
	items := make([]types.DepositPool, n)
	for i := range items {
		valPubk := secp256k1.GenPrivKey().PubKey()
		valAddr := sdk.ValAddress(valPubk.Address())
		items[i].OperatorAddress = valAddr.String()
		items[i].Shares = sdk.NewDec(int64(1000 * i))
		items[i].Tokens = sdk.NewInt64Coin("stake", int64(1000*i))
		keeper.SetDepositPool(ctx, items[i])
	}
	return items
}

func CreateNEntries(n int) []types.UnbondingDepositEntry {

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

// Creates n different unbonding deposit in the store, each with nEntries entries.
func CreateNUnbondingDeposit(keeper *keeper.Keeper, ctx sdk.Context, n int, nEntries int) []types.UnbondingDeposit {
	items := make([]types.UnbondingDeposit, n)
	for i := range items {
		items[i].DepositorAddress = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
		items[i].ValidatorAddress = sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
		items[i].Entries = CreateNEntries(nEntries)
		keeper.SetUnbondingDeposit(ctx, items[i])
	}
	return items
}