package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"

	"github.com/made-in-block/slash-refund/app"
	"github.com/made-in-block/slash-refund/testutil/testsuite"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	"github.com/stretchr/testify/require"
)

type KeeperTestSuite struct {
	srApp          *app.App
	ctx            sdk.Context
	units          int64
	testAddrs      []sdk.AccAddress
	valAddrs       []sdk.ValAddress
	selfDelegation sdk.Int
	t              *testing.T
}

// Default initial state for all test. It creates two validators with a specified power.
func SetupTestSuite(t *testing.T, power int64) *KeeperTestSuite {

	srApp, ctx := testsuite.CreateTestApp(false)

	units := srApp.StakingKeeper.PowerReduction(ctx).Int64()

	initAmt := sdk.NewInt(int64(1000 * units))
	testAddrs, pubks := CreateNTestAccounts(srApp, ctx, 5, initAmt)

	selfDelegation := sdk.NewInt(power * units)

	// create 2 validators with consensous power equal to input power
	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)

	valAddrs := make([]sdk.ValAddress, 0, 2)

	for i := 0; i < 2; i++ {
		sth.CreateValidatorWithValPower(sdk.ValAddress(testAddrs[i]), pubks[i], selfDelegation.QuoRaw(units).Int64(), true)
		validator, found := srApp.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(testAddrs[i]))
		require.True(t, found)
		valAddr := validator.GetOperator()
		sd, found := srApp.StakingKeeper.GetDelegation(ctx, testAddrs[i], valAddr)
		require.True(t, found)
		require.Equal(t, selfDelegation, sd.Shares.TruncateInt())
		valAddrs = append(valAddrs, valAddr)
	}

	s := KeeperTestSuite{}
	s.srApp, s.ctx, s.units, s.testAddrs, s.valAddrs, s.selfDelegation, s.t = srApp, ctx, units, testAddrs, valAddrs, selfDelegation, t

	return &s
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
