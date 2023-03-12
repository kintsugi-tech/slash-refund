package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/made-in-block/slash-refund/app"
	"github.com/made-in-block/slash-refund/testutil/testsuite"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stretchr/testify/require"
)

var (
	numDelAddrs = 1
	numValAddrs = 1
	numDepAddrs = 1
)

func SetupMsgServerTest() (
	*app.App,
	sdk.Context,
	[]sdk.AccAddress,
	[]sdk.ValAddress,
	[]sdk.AccAddress,
) {

	// Setup delegators
	delAddrs := testsuite.GenerateNAddresses(numDelAddrs)
	delAccs := testsuite.ConvertAddressesToAccAddr(delAddrs)
	balances := testsuite.GenerateBalances(delAccs)

	// Setup validators
	valAddrs := testsuite.GenerateNAddresses(numValAddrs)
	valAccs := testsuite.ConvertAddressesToValAddr(valAddrs)

	// Setup depositors
	depAddrs := testsuite.GenerateNAddresses(numDepAddrs)
	depAccs := testsuite.ConvertAddressesToAccAddr(depAddrs)
	depBalances := testsuite.GenerateBalances(depAccs)

	balances = append(balances, depBalances...)

	testInputs := testsuite.TestInputs{
		DelAddrs: delAccs,
		ValAddrs: valAccs,
		Balances: balances,
	}

	app, ctx := testsuite.CreateTestApp(testInputs, false)

	return app, ctx, delAccs, valAccs, depAccs
}

// -------------------------------------------------------------------------------------------------
// Deposit
// -------------------------------------------------------------------------------------------------
func TestDepositValid(t *testing.T) {
	app, ctx, _, validators, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgDeposit{
		DepositorAddress: depositors[0].String(),
		ValidatorAddress: validators[0].String(),
		Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(1)),
	}
	_, err := msgServer.Deposit(ctx, msg)
	require.NoError(t, err)
}

func TestDepositNotValidatorBech32(t *testing.T) {
	app, ctx, _, _, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgDeposit{
		DepositorAddress: depositors[0].String(),
		ValidatorAddress: depositors[0].String(),
		Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(1)),
	}
	_, err := msgServer.Deposit(ctx, msg)
	require.Error(t, err)
}

func TestDepositNotValidator(t *testing.T) {
	app, ctx, _, _, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	val := sdk.MustBech32ifyAddressBytes(bech32PrefixValAddr, depositors[0])

	msg := &types.MsgDeposit{
		DepositorAddress: depositors[0].String(),
		ValidatorAddress: val,
		Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(1)),
	}
	_, err := msgServer.Deposit(ctx, msg)
	require.ErrorIs(t, err, stakingtypes.ErrNoValidatorFound)
}

func TestDepositNotAccountBech32(t *testing.T) {
	app, ctx, _, validators, _ := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgDeposit{
		DepositorAddress: validators[0].String(),
		ValidatorAddress: validators[0].String(),
		Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(1)),
	}
	_, err := msgServer.Deposit(ctx, msg)
	require.Error(t, err)
}

func TestDepositNotAllowedTokens(t *testing.T) {
	app, ctx, _, validators, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgDeposit{
		DepositorAddress: depositors[0].String(),
		ValidatorAddress: validators[0].String(),
		Amount:           sdk.NewCoin("mib", sdk.NewInt(1)),
	}
	_, err := msgServer.Deposit(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
}

func TestDepositZeroAmount(t *testing.T) {
	app, ctx, _, validators, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgDeposit{
		DepositorAddress: depositors[0].String(),
		ValidatorAddress: validators[0].String(),
		Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(0)),
	}
	_, err := msgServer.Deposit(ctx, msg)
	require.ErrorIs(t, err, types.ErrNonPositiveDeposit)
}

// -------------------------------------------------------------------------------------------------
// Withdraw
// -------------------------------------------------------------------------------------------------
func TestWithdrawNotValidatorBech32(t *testing.T) {
	app, ctx, _, _, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgWithdraw{
		DepositorAddress: depositors[0].String(),
		ValidatorAddress: depositors[0].String(),
		Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(1)),
	}
	_, err := msgServer.Withdraw(ctx, msg)
	require.Error(t, err)
}

func TestWithdrawNotAccountBech32(t *testing.T) {
	app, ctx, _, validators, _ := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgWithdraw{
		DepositorAddress: validators[0].String(),
		ValidatorAddress: validators[0].String(),
		Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(1)),
	}
	_, err := msgServer.Withdraw(ctx, msg)
	require.Error(t, err)
}

func TestWithdrawNotAllowedTokens(t *testing.T) {

	app, ctx, _, validators, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgWithdraw{
		DepositorAddress: depositors[0].String(),
		ValidatorAddress: validators[0].String(),
		Amount:           sdk.NewCoin("mib", sdk.NewInt(1)),
	}
	_, err := msgServer.Withdraw(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
}

func TestWithdrawZeroAmount(t *testing.T) {
	app, ctx, _, validators, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgWithdraw{
		DepositorAddress: depositors[0].String(),
		ValidatorAddress: validators[0].String(),
		Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(0)),
	}
	_, err := msgServer.Withdraw(ctx, msg)
	require.ErrorIs(t, err, types.ErrNonPositiveWithdraw)
}

// -------------------------------------------------------------------------------------------------
// Claim
// -------------------------------------------------------------------------------------------------
func TestClaimNotValidatorBech32(t *testing.T) {
	app, ctx, _, _, depositors := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgClaim{
		DelegatorAddress: depositors[0].String(),
		ValidatorAddress: depositors[0].String(),
	}
	_, err := msgServer.Claim(ctx, msg)
	require.Error(t, err)
}

func TestClaimNotAccountBech32(t *testing.T) {
	app, ctx, _, validators, _ := SetupMsgServerTest()

	msgServer := keeper.NewMsgServerImpl(app.SlashrefundKeeper)

	msg := &types.MsgClaim{
		DelegatorAddress: validators[0].String(),
		ValidatorAddress: validators[0].String(),
	}
	_, err := msgServer.Claim(ctx, msg)
	require.Error(t, err)
}

/*

func TestMsgServerClaim(t *testing.T) {

	s := bootstrapRefundTest(t, 100)
	srApp, ctx, testAddrs, valAddrs := s.srApp, s.ctx, s.testAddrs, s.valAddrs
	msgServer := keeper.NewMsgServerImpl(s.srApp.SlashrefundKeeper)

	expectedInPool := sdk.NewInt(10e6)

	module := types.ModuleName
	err := banktestutil.FundModuleAccount(srApp.BankKeeper, ctx, module, sdk.NewCoins(sdk.NewCoin(types.DefaultAllowedTokens[0], expectedInPool)))
	require.NoError(t, err)

	coin := sdk.NewCoin(types.DefaultAllowedTokens[0], expectedInPool)
	refundPool := types.NewRefundPool(valAddrs[0], coin, sdk.NewDecFromInt(expectedInPool))
	srApp.SlashrefundKeeper.SetRefundPool(ctx, refundPool)

	refAmounts := []sdk.Int{sdk.NewInt(5e6), sdk.NewInt(4e6), sdk.NewInt(1e6)}

	refund0 := types.NewRefund(testAddrs[0], valAddrs[0], sdk.NewDecFromInt(refAmounts[0]))
	srApp.SlashrefundKeeper.SetRefund(ctx, refund0)

	refund1 := types.NewRefund(testAddrs[1], valAddrs[0], sdk.NewDecFromInt(refAmounts[1]))
	srApp.SlashrefundKeeper.SetRefund(ctx, refund1)

	refund2 := types.NewRefund(testAddrs[2], valAddrs[0], sdk.NewDecFromInt(refAmounts[2]))
	srApp.SlashrefundKeeper.SetRefund(ctx, refund2)

	refunds := []types.Refund{refund0, refund1, refund2}

	for i := 3; i < 5; i++ {

		msg := &types.MsgClaim{
			DelegatorAddress: testAddrs[i].String(),
			ValidatorAddress: valAddrs[0].String(),
		}

		_, err := msgServer.Claim(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, err, types.ErrNoRefundForAddress)
	}

	s.RequireRefundPool(valAddrs[0], expectedInPool, sdk.NewDecFromInt(expectedInPool), refunds[0:3])

	for i := 0; i < 3; i++ {

		s.RequireRefundPool(valAddrs[0], expectedInPool, sdk.NewDecFromInt(expectedInPool), refunds[i:3])

		msg := &types.MsgClaim{
			DelegatorAddress: testAddrs[i].String(),
			ValidatorAddress: valAddrs[0].String(),
		}

		initialBalance := srApp.BankKeeper.GetBalance(ctx, testAddrs[i], types.DefaultAllowedTokens[0])

		s.RequireRefund(testAddrs[i], valAddrs[0], refunds[i].Shares)

		_, err = msgServer.Claim(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
		require.Equal(t,
			initialBalance.Amount.Add(refAmounts[i]),
			srApp.BankKeeper.GetBalance(ctx, testAddrs[i], types.DefaultAllowedTokens[0]).Amount)

		s.RequireNoRefund(testAddrs[i], valAddrs[0])
		expectedInPool = expectedInPool.Sub(refAmounts[i])
	}

	s.RequireNoRefundPool(valAddrs[0])

}

func TestMsgServerClaim_Errors(t *testing.T) {

	s := bootstrapRefundTest(t, 100)
	srApp, ctx, testAddrs, valAddrs := s.srApp, s.ctx, s.testAddrs, s.valAddrs
	msgServer := keeper.NewMsgServerImpl(s.srApp.SlashrefundKeeper)

	initialBalance := srApp.BankKeeper.GetBalance(ctx, testAddrs[0], types.DefaultAllowedTokens[0])

	// test case: invalid validator address
	msg := &types.MsgClaim{
		DelegatorAddress: testAddrs[0].String(),
		ValidatorAddress: "not a valid address",
	}
	_, err := msgServer.Claim(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Equal(t, initialBalance, srApp.BankKeeper.GetBalance(ctx, testAddrs[0], types.DefaultAllowedTokens[0]))

	// test case: invalid validator address
	msg = &types.MsgClaim{
		DelegatorAddress: "not a valid address",
		ValidatorAddress: valAddrs[0].String(),
	}
	_, err = msgServer.Claim(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Equal(t, initialBalance, srApp.BankKeeper.GetBalance(ctx, testAddrs[0], types.DefaultAllowedTokens[0]))

	// test case: refund not found
	msg = &types.MsgClaim{
		DelegatorAddress: testAddrs[0].String(),
		ValidatorAddress: valAddrs[0].String(),
	}
	_, err = msgServer.Claim(sdk.WrapSDKContext(ctx), msg)
	require.ErrorIs(t, err, types.ErrNoRefundForAddress)
	require.Equal(t, initialBalance, srApp.BankKeeper.GetBalance(ctx, testAddrs[0], types.DefaultAllowedTokens[0]))
}

func TestMsgServerClaim_RefundWithoutRefundPool(t *testing.T) {

	s := bootstrapRefundTest(t, 100)
	srApp, ctx, testAddrs, valAddrs := s.srApp, s.ctx, s.testAddrs, s.valAddrs
	msgServer := keeper.NewMsgServerImpl(s.srApp.SlashrefundKeeper)

	refund := types.NewRefund(testAddrs[0], valAddrs[0], sdk.NewDec(1))
	srApp.SlashrefundKeeper.SetRefund(ctx, refund)

	msg := &types.MsgClaim{
		DelegatorAddress: testAddrs[0].String(),
		ValidatorAddress: valAddrs[0].String(),
	}
	require.Panics(t, func() { msgServer.Claim(sdk.WrapSDKContext(ctx), msg) })
}

func TestMsgServerClaim_BlockedAddress(t *testing.T) {

	s := bootstrapRefundTest(t, 100)
	srApp, ctx, testAddrs, valAddrs := s.srApp, s.ctx, s.testAddrs, s.valAddrs
	msgServer := keeper.NewMsgServerImpl(s.srApp.SlashrefundKeeper)

	// get a blocked address from test keeper setup
	blockedAddress := srApp.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	require.True(t, srApp.BankKeeper.BlockedAddr(blockedAddress))
	initialBalance := srApp.BankKeeper.GetBalance(ctx, blockedAddress, types.DefaultAllowedTokens[0])

	// set refunds and refund pool
	amt1, amt2 := sdk.NewInt(123), sdk.NewInt(877)
	amt := amt1.Add(amt2)

	refund := types.NewRefund(blockedAddress, valAddrs[0], sdk.NewDecFromInt(amt1))
	srApp.SlashrefundKeeper.SetRefund(ctx, refund)

	refund = types.NewRefund(testAddrs[0], valAddrs[0], sdk.NewDecFromInt(amt2))
	srApp.SlashrefundKeeper.SetRefund(ctx, refund)

	coin := sdk.NewCoin(types.DefaultAllowedTokens[0], amt)
	refundPool := types.NewRefundPool(valAddrs[0], coin, sdk.NewDecFromInt(amt))
	srApp.SlashrefundKeeper.SetRefundPool(ctx, refundPool)

	// fund module account
	err := banktestutil.FundModuleAccount(srApp.BankKeeper, ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.DefaultAllowedTokens[0], amt)))
	require.NoError(t, err)

	// process message
	msg := &types.MsgClaim{
		DelegatorAddress: blockedAddress.String(),
		ValidatorAddress: valAddrs[0].String(),
	}
	_, err = msgServer.Claim(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)

	// check that refund has not been transferred
	refund1 := s.RequireRefund(blockedAddress, valAddrs[0], sdk.NewDecFromInt(amt1))
	refund2 := s.RequireRefund(testAddrs[0], valAddrs[0], sdk.NewDecFromInt(amt2))
	s.RequireRefundPool(valAddrs[0], amt, sdk.NewDecFromInt(amt), []types.Refund{refund1, refund2})
	require.Equal(t, initialBalance, srApp.BankKeeper.GetBalance(ctx, blockedAddress, types.DefaultAllowedTokens[0]))
}

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := testslashrefund.NewTestKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
*/
