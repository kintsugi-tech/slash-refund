package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/app"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

type KeeperTestSuite struct {
	srApp          *app.App
	ctx            sdk.Context
	testAddrs      []sdk.AccAddress
	valAddrs       []sdk.ValAddress
	selfDelegation sdk.Int
	t              *testing.T
}

func (s KeeperTestSuite) RequireNoRefund(addr sdk.AccAddress, valAddr sdk.ValAddress) {
	_, found := s.srApp.SlashrefundKeeper.GetRefund(s.ctx, addr, valAddr)
	require.False(s.t, found, "refund found")
}

func (s KeeperTestSuite) RequireRefund(addr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) types.Refund {
	refund, found := s.srApp.SlashrefundKeeper.GetRefund(s.ctx, addr, valAddr)
	require.True(s.t, found, "refund not found")
	require.Equal(s.t, shares, refund.Shares, "shares not equal")

	return refund
}

func (s KeeperTestSuite) RequireNoRefundPool(valAddr sdk.ValAddress) {

	_, found := s.srApp.SlashrefundKeeper.GetRefundPool(s.ctx, valAddr)
	require.False(s.t, found, "refund pool found")
}

func (s KeeperTestSuite) RequireRefundPool(valAddr sdk.ValAddress, tokens sdk.Int, shares sdk.Dec, refunds []types.Refund,
) types.RefundPool {

	pool, found := s.srApp.SlashrefundKeeper.GetRefundPool(s.ctx, valAddr)
	require.True(s.t, found, "refund pool not found")
	require.Equal(s.t, tokens, pool.Tokens.Amount, "tokens not equal")
	require.Equal(s.t, shares, pool.Shares, "shares not equal")

	// check associated refunds shares
	total := sdk.NewDec(0)
	for i := 0; i < len(refunds); i++ {
		total = total.Add(refunds[i].Shares)
	}
	require.Equal(s.t, total, pool.Shares, "refunds and refund pool mismatch")

	return pool
}

func (s KeeperTestSuite) RequireNoDeposit(addr sdk.AccAddress, valAddr sdk.ValAddress) {
	_, found := s.srApp.SlashrefundKeeper.GetDeposit(s.ctx, addr, valAddr)
	require.False(s.t, found, "deposit found")

}

func (s KeeperTestSuite) RequireDeposit(addr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) types.Deposit {
	deposit, found := s.srApp.SlashrefundKeeper.GetDeposit(s.ctx, addr, valAddr)
	require.True(s.t, found, "deposit not found")
	require.Equal(s.t, shares, deposit.Shares)

	return deposit
}

func (s KeeperTestSuite) RequireNoDepositPool(valAddr sdk.ValAddress) {
	_, found := s.srApp.SlashrefundKeeper.GetDepositPool(s.ctx, valAddr)
	require.False(s.t, found, "deposit pool found")
}

func (s KeeperTestSuite) RequireDepositPool(valAddr sdk.ValAddress, tokens sdk.Int, shares sdk.Dec, deposits []types.Deposit,
) types.DepositPool {

	pool, found := s.srApp.SlashrefundKeeper.GetDepositPool(s.ctx, valAddr)
	require.True(s.t, found, "deposit pool not found")
	require.Equal(s.t, tokens, pool.Tokens.Amount, "tokens not equal")
	require.Equal(s.t, shares, pool.Shares, "shares not equal")

	// check associated deposits shares
	total := sdk.NewDec(0)
	for i := 0; i < len(deposits); i++ {
		total = total.Add(deposits[i].Shares)
	}
	require.Equal(s.t, total, pool.Shares, "deposits and deposit pool mismatch")

	return pool
}
